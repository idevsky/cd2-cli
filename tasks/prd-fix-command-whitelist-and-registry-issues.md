# PRD: Fix Command Whitelist and Registry Issues

## Overview
修复 cd2-cli 项目中命令白名单和注册表的关键问题，确保本地命令不会误用远端 gRPC，白名单机制正确工作，以及注册表防止重复 ID。

## Goals
- P0 问题全部修复：local 命令语义正确、completion 命令可用
- P1 问题全部修复：白名单初始化无副作用、alias 支持完整、registry 防重复
- P2 问题全部修复：upload 失败清理资源
- 所有测试通过

## Quality Gates

These commands must pass for every user story:
- `go test ./...` - Unit tests
- `go vet ./...` - Static analysis
- `make build` - Build succeeds

## User Stories

### US-001: Fix local command semantics
**Description:** As a user, I want `local list/mkdir` commands to work on local filesystem only, not call remote gRPC.

**Acceptance Criteria:**
- [ ] `local list` uses local filesystem, not gRPC client
- [ ] `local mkdir` uses local filesystem, not gRPC client
- [ ] Commands marked as "local" are correctly categorized in registry

### US-002: Fix completion command blocked by whitelist
**Description:** As a user, I want `completion` command to work without whitelist configuration.

**Acceptance Criteria:**
- [ ] `completion bash` outputs shell completion script
- [ ] `completion zsh` outputs shell completion script
- [ ] No whitelist config file required for completion

### US-003: Fix whitelist initialization side effects
**Description:** As a user, I want `local token set` to not create whitelist config files.

**Acceptance Criteria:**
- [ ] `local token set` does not create whitelist config
- [ ] Only whitelist-related commands initialize whitelist config
- [ ] Non-whitelist commands skip whitelist initialization entirely

### US-004: Support alias in whitelist enable/disable
**Description:** As a user, I want to enable/disable commands using their alias names.

**Acceptance Criteria:**
- [ ] `whitelist enable storage.add.s3` enables `cloudapi.login-s3` (alias group)
- [ ] `whitelist disable cloudapi.login-s3` disables `storage.add.s3` (canonical)
- [ ] `whitelist status` shows both requested ID and canonical ID

### US-005: Prevent duplicate command IDs in registry
**Description:** As a developer, I want registry to panic on duplicate command IDs to catch bugs early.

**Acceptance Criteria:**
- [ ] `MustRegister` panics on duplicate ID
- [ ] `Register` allows override (for test compatibility)
- [ ] `registry_init.go` uses `MustRegister` for all commands

### US-006: Fix upload failure resource leak
**Description:** As a user, I want failed uploads to properly close remote file handles.

**Acceptance Criteria:**
- [ ] Upload errors close file handle before returning
- [ ] No orphaned file handles after upload failure

## Functional Requirements
- FR-1: Commands marked as `Local: true` must not call gRPC client
- FR-2: `completion` and `help` commands bypass whitelist check
- FR-3: Whitelist config only created by whitelist-related commands
- FR-4: Alias group members share whitelist policy
- FR-5: Registry `MustRegister` prevents duplicate registrations

## Non-Goals
- Refactoring entire command structure
- Adding new features or commands
- Performance optimization

## Technical Considerations
- Use `sync.Once` for registry initialization
- Add `skipWhitelistCheck` annotation for completion/help
- Use `MustRegister` in production, `Register` for tests

## Success Metrics
- All P0/P1/P2 issues resolved
- All tests pass
- No regression in existing functionality

## Open Questions
- None - requirements are clear from Problem.md