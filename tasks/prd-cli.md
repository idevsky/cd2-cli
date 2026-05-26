# PRD: CLI 安全性与可靠性修复

## Overview

修复 cd2-cli 中影响安全、可靠性和自动化调用正确性的关键问题。按 P0→P1→P2 优先级分批实施，确保白名单机制有效、配置文件操作安全、退出码正确、远程调用有超时保护。

## Goals

- 白名单机制可靠，无绕过路径
- 配置文件操作保留已有字段，不破坏用户设置
- 所有远程 API 调用有统一超时控制
- 命令失败时退出码非零，便于自动化判断
- 环境变量、flag、配置文件优先级正确
- 所有远程命令纳入白名单管理

## Quality Gates

These commands must pass for every user story:
- `go test ./...` - Unit tests
- `go vet ./...` - Static analysis
- `make build` - Build verification

For integration tests:
- `make test-integration` - Integration tests (when environment available)

## User Stories

### US-001: 分离白名单配置与主配置文件

As a security-conscious user, I want whitelist configuration to use a separate file so that providing a main config with server/token doesn't accidentally disable whitelist protection.

**Acceptance Criteria:**
- [ ] Add `--whitelist-config` flag for explicit whitelist file path
- [ ] Default whitelist path: `~/.cd2-cli-whitelist.yaml` (separate from main config)
- [ ] `--config` flag only controls main config, never whitelist config
- [ ] When main config lacks `whitelist_enabled`, whitelist defaults to enabled
- [ ] Test: `--config` pointing to server/token-only file still blocks high-risk commands

### US-002: 配置写入保留已有字段

As a user with existing configuration, I want `auth token set` and `auth login --save` to preserve my other config fields so that I don't lose server/tls/skip-tls-verify settings.

**Acceptance Criteria:**
- [ ] `saveTokenToConfig()` reads existing YAML to `map[string]any` before writing
- [ ] Only updates `token` field, preserves all other fields
- [ ] File written with `0600` permissions
- [ ] Use atomic write: write to temp file, chmod, rename
- [ ] `auth token clear` removes only `token` field, preserves others
- [ ] Test: config with server/tls/json fields survives `auth token set`

### US-003: 全局 --timeout flag 与统一超时控制

As an automation user, I want all remote API calls to have a configurable timeout so that network issues don't cause indefinite hangs.

**Acceptance Criteria:**
- [ ] Add global `--timeout` flag with default `30s`
- [ ] Support config file and `CD2_CLI_TIMEOUT` environment variable
- [ ] Create `runWithClient()` helper that wraps context with timeout
- [ ] All remote commands use the helper, no scattered `context.Background()`
- [ ] Long-stream APIs have `--max-events` or explicit timeout handling
- [ ] Test: blocked RPC causes command to exit non-zero after timeout

### US-004: 命令失败返回非零退出码

As an automation user, I want failed remote API calls to return non-zero exit codes so that shell scripts and agents can correctly detect failures.

**Acceptance Criteria:**
- [ ] All remote API commands use `RunE` instead of `Run`
- [ ] Errors returned from RunE propagate to root command
- [ ] Root command configured with `SilenceUsage: true`, `SilenceErrors: true`
- [ ] JSON error output handled centrally, then `os.Exit(1)`
- [ ] Local config commands also return non-zero on write failures
- [ ] Test: simulate RPC error, verify exit code is non-zero

### US-005: 移除白名单通配符支持

As a security-conscious user, I want whitelist to require explicit path lists so that `*` patterns cannot bypass whitelist protection.

**Acceptance Criteria:**
- [ ] Remove `*` wildcard support from whitelist configuration
- [ ] Whitelist only accepts explicit file paths
- [ ] Error message clearly indicates invalid path format if wildcard used
- [ ] Test: whitelist with `*` pattern is rejected or treated as literal path

### US-006: 为所有远程命令设置 command ID

As a whitelist administrator, I want all remote commands to have command IDs so that whitelist can properly control access.

**Acceptance Criteria:**
- [ ] Audit all commands with `Run`/`RunE` that access remote APIs
- [ ] Add `setCommandID()` to commands missing it: cache prefetch/cancel/close-reader, storage add/config/set-config/status/can-add, task cancel/list/status/wait variants
- [ ] `checkWhitelist()` rejects commands with empty command ID (default deny, not allow)
- [ ] Local commands explicitly marked with `annotationIsLocal`
- [ ] Test: scan Cobra tree, all non-local commands with Run/RunE have command IDs in registry

### US-007: 修复全局 flag 绑定到 Viper

As a user, I want `--tls`, `--skip-tls-verify`, `--json` flags to work correctly so that command-line overrides take effect.

**Acceptance Criteria:**
- [ ] Bind all persistent flags to Viper using `viper.BindPFlag()`
- [ ] Add `viper.SetEnvPrefix("CD2_CLI")` and `SetEnvKeyReplacer`
- [ ] `--json=false` correctly disables JSON output
- [ ] `--tls --skip-tls-verify` on CLI override config defaults
- [ ] Test: flags override config file and environment variables

### US-008: 本地命令不初始化 gRPC client

As a user running local commands, I want `whitelist list` and `auth token set` to not attempt gRPC connections so that they work without a server.

**Acceptance Criteria:**
- [ ] Local commands skip `initClient()` entirely
- [ ] `shouldSkipWhitelistCheck()` returns true without creating client
- [ ] Local commands only need config/whitelist manager
- [ ] Test: `whitelist list` succeeds without any server running

### US-009: 环境变量优先级与 token 来源判断

As a user, I want `CD2_CLI_TOKEN` to override config file token and `auth token show --redact` to report correct source so that I understand where credentials come from.

**Acceptance Criteria:**
- [ ] `CD2_CLI_*` environment variables work for all config keys
- [ ] Token source priority: flag > env > config > none
- [ ] Use `cmd.Flags().Changed("token")` to detect flag source
- [ ] `auth token show --redact` correctly reports source
- [ ] Test: `CD2_CLI_TOKEN` overrides config file token

### US-010: 修复 copy remove 错误调用

As a user, I want `copy remove` to call the correct RPC so that command names match actual behavior.

**Acceptance Criteria:**
- [ ] Verify whether `copy remove` should call `RemoveCopyTasks` or `CancelCopyTask`
- [ ] Either fix RPC call or rename command to `copy cancel`
- [ ] Add command ID if missing
- [ ] Test: command behavior matches its name

### US-011: 注册完整性双向校验

As a developer, I want tests to verify both that Cobra commands have registry entries and registry entries have Cobra commands so that no phantom or missing commands exist.

**Acceptance Criteria:**
- [ ] Add reverse test: registry commands with `Implemented=true` must exist in Cobra tree
- [ ] Add test: all Cobra remote commands have command IDs in registry
- [ ] Filter or mark registry phantom commands (cache.dir-table, system.cache-stats, etc.)
- [ ] Test: `cd2-cli task --help` shows no duplicate commands

### US-012: 使用 DefaultOpen 字段决定白名单默认状态

As a whitelist administrator, I want `CommandSpec.DefaultOpen` to be respected so that I can override risk-level-based defaults for specific commands.

**Acceptance Criteria:**
- [ ] Change `DefaultOpen` to `*bool` to distinguish unset from false
- [ ] `initDefaultConfig()` checks `DefaultOpen` before using risk level
- [ ] Write operations (copy, move, mkdir, upload) default to closed
- [ ] Read-only commands default to open
- [ ] Test: `file.move` blocked by default whitelist

### US-013: 顶层别名命令共享白名单策略

As a whitelist administrator, I want `ls` and `file list`, `rm` and `file delete` to share the same whitelist policy so that users cannot bypass whitelist through alias paths.

**Acceptance Criteria:**
- [ ] Audit all alias pairs: ls/file.list, stat/file.find, mkdir/file.mkdir, mv/file.move, cp/file.copy, rm/file.delete, upload/file.upload, download/file.download
- [ ] Aliases share same command ID or are grouped in whitelist
- [ ] Test: same RPC has consistent risk level across all CLI entry points

### US-014: 使用 protojson 处理 protobuf 输入输出

As an automation user, I want JSON input/output to follow proto3 JSON format so that fields like Timestamp, bytes, and enums are handled correctly.

**Acceptance Criteria:**
- [ ] Create helper `parseProtoJSON(data []byte, msg proto.Message)` using `protojson.UnmarshalOptions`
- [ ] Replace `encoding/json` usage for proto messages in all `--request-json` handling
- [ ] `outputResult()` uses `protojson.MarshalOptions` for proto messages
- [ ] Set `EmitUnpopulated: true` for consistent output
- [ ] Test: Timestamp fields output in standard proto JSON format

### US-015: 配置文件不存在时允许创建

As a user, I want `auth token set --config ./new.yaml` to create the config file if it doesn't exist so that I can start fresh without errors.

**Acceptance Criteria:**
- [ ] Local config write commands tolerate missing config file
- [ ] Create file with secure defaults on first write
- [ ] Read-only commands can operate with defaults if config missing
- [ ] Test: `auth token set` with non-existent `--config` succeeds

### US-016: task 命令组纳入 registry 和白名单

As a whitelist administrator, I want `task` commands to be controllable via whitelist so that task cancellation and monitoring can be restricted.

**Acceptance Criteria:**
- [ ] Add `registerTaskCommands()` to registry initialization
- [ ] Add command IDs: task.list, task.list-uploads, task.list-downloads, task.cancel-upload, task.cancel-copy, task.cancel-merge, task.wait-copy, task.wait-merge
- [ ] Remove duplicate `taskCmd.AddCommand(taskStatusCmd)`
- [ ] Test: `task cancel-upload` blocked by default whitelist

### US-017: 合并 registry 新命令到已有白名单

As a user upgrading CLI, I want new commands to appear in whitelist without losing my existing enable/disable settings so that I don't need to reset whitelist.

**Acceptance Criteria:**
- [ ] `Load()` followed by `MergeRegistryDefaults()` in `NewManager()`
- [ ] Existing commands preserve user's `Enabled` state
- [ ] New registry commands appended with default policy
- [ ] Deleted registry commands marked deprecated, not removed from user config
- [ ] Test: create old whitelist file, add new command to registry, verify merge

### US-018: 修复 task wait 未找到任务返回 completed

As an automation user, I want `task wait-*` to return `not_found` status with non-zero exit when task doesn't exist so that I can distinguish "never started" from "completed".

**Acceptance Criteria:**
- [ ] First poll returning no tasks yields `not_found` status
- [ ] Exit code non-zero for `not_found`, `failed`, `cancelled`, `timeout`
- [ ] Optional `--missing-is-complete` flag for backward compatibility
- [ ] Test: waiting for non-existent task returns `not_found`

### US-019: 敏感参数支持安全输入方式

As a security-conscious user, I want to provide passwords and tokens via environment variables or files so that they don't appear in shell history.

**Acceptance Criteria:**
- [ ] `auth login USERNAME [PASSWORD]` - if PASSWORD missing, read from `CD2_CLI_PASSWORD`
- [ ] `auth login-2fa` reads TOTP from `CD2_CLI_TOTP` if not provided
- [ ] `auth change-password` reads passwords from env vars or files
- [ ] `cloudapi login-*-oauth` reads refresh token from `CD2_CLI_REFRESH_TOKEN`
- [ ] Add `--password-file`, `--token-file` options where applicable
- [ ] Test: `auth login USER` reads password from `CD2_CLI_PASSWORD`

### US-020: storage add 与 cloudapi login 统一策略

As a whitelist administrator, I want `storage add s3` and `cloudapi login-s3` to share the same risk policy so that users can't bypass whitelist through different entry points.

**Acceptance Criteria:**
- [ ] `storage add` commands get command IDs
- [ ] `storage add s3` shares policy with `cloudapi.login-s3` (both high-risk, default closed)
- [ ] Test: same RPC has consistent default status across all entry points

### US-021: 可选 bool flag 只在显式指定时设置

As a user, I want `webdav user modify alice --password x` to not modify `enabled`/`readOnly`/`guest` fields so that I can change just one field.

**Acceptance Criteria:**
- [ ] Use `cmd.Flags().Changed("flag-name")` to detect explicit flag usage
- [ ] Only set optional bool pointers when flag explicitly provided
- [ ] Audit all commands with optional bool flags
- [ ] Test: `webdav user modify` without bool flags doesn't set those fields in request

### US-022: JSON 参数命令强制要求 flag

As a user, I want clear errors when I forget `--config` or `--options` flags so that I don't get confusing JSON parse errors.

**Acceptance Criteria:**
- [ ] Use `MarkFlagRequired()` for mandatory JSON flags
- [ ] Or check empty value at RunE start with clear error message
- [ ] Affected commands: `mount update --options`, `webdav server set --config`, `storage set-config --config`, `cloudapi set-config --config`
- [ ] Test: missing required flag returns clear error, not JSON parse error

### US-023: 文档同步与清理

As a developer, I want README and AGENTS.md to reflect current behavior so that I'm not misled by outdated documentation.

**Acceptance Criteria:**
- [ ] Update default config path references (`~/.config/cd2-cli.yaml`)
- [ ] Document `--timeout` flag after implementation
- [ ] Document `CD2_CLI_PASSWORD` usage for `auth login`
- [ ] Archive or update `Todo.md` to reflect completed items
- [ ] Document security implications of JSON output for auth commands

### US-024: 删除废弃代码与清理

As a developer, I want deprecated functions and unused code removed so that the codebase is maintainable.

**Acceptance Criteria:**
- [ ] Remove `DefaultHighRiskCommands` and `IsHighRisk()` if no longer used
- [ ] Remove duplicate command registrations
- [ ] Clean up unused global state in root.go if safe
- [ ] Test: all existing tests still pass

### US-025: 审计与默认状态调整

As a security-conscious user, I want write operations to be blocked by default so that a fresh CLI installation is safe for automation.

**Acceptance Criteria:**
- [ ] Audit all commands: read-only defaults to open, write/delete/modify defaults to closed
- [ ] `file.copy`, `file.move`, `file.mkdir`, `fs.upload` blocked by default
- [ ] `auth.logout`, `auth.register`, `auth.2fa-recovery-codes` blocked by default
- [ ] Test: default whitelist blocks write commands, allows read commands

## Functional Requirements

- FR-1: All remote API calls must have bounded context with timeout
- FR-2: Commands with empty command ID must be denied by whitelist (not allowed)
- FR-3: Config file writes must preserve existing fields
- FR-4: Exit code must be non-zero on any error (remote or local)
- FR-5: Flag values override environment variables, which override config file
- FR-6: Local commands must not initialize gRPC client
- FR-7: Whitelist must use separate file from main config
- FR-8: New commands added to registry must merge with existing whitelist
- FR-9: All `RunE` commands must return error on failure
- FR-10: JSON input/output for protobuf must use protojson, not encoding/json
- FR-11: Optional flags must only set fields when explicitly provided

## Non-Goals

- System theme auto-detection for any TUI features
- Custom color schemes beyond light/dark
- Per-command whitelist overrides
- Backward compatibility with `*` wildcard patterns in whitelist
- Supporting legacy command names without migration path

## Technical Considerations

- Use `runWithClient()` helper pattern to reduce code duplication
- Consider creating a centralized error handler for consistent exit codes
- protojson options: `Multiline: true`, `Indent: "  "`, `EmitUnpopulated: true`
- Atomic file writes: temp file → chmod 0600 → rename
- Context timeout must be configurable via flag, env, and config

## Success Metrics

- All P0 issues resolved
- All tests pass including integration tests
- No security bypass paths in whitelist
- All remote commands have command IDs and proper exit codes
- Config operations preserve existing fields
- Timeout protection on all network calls

## Open Questions

- Should `file download` be renamed to `file get-download-url` or actually implement file download?
- Should we support `--request-file` for large JSON payloads?
- How to handle deprecation of existing wildcard patterns in user configs?
- Should long-stream APIs have default `--max-events` limit?