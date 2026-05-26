# PRD: CLI API 全开放与白名单治理

## Overview

为 cd2-cli 实现完整的 CloudDrive2 gRPC API 暴露，并通过本地白名单机制进行治理。高风险 API 默认阻止执行，必须显式启用后才能使用。配置文件路径统一为 `~/.config/cd2-cli.yaml`。

## Goals

- CLI 暴露 CloudDrive2 gRPC API 的完整可用面
- 所有远程 API 命令纳入本地白名单管理
- 高风险 API 默认阻止，需显式加入白名单后执行
- 支持白名单 CLI 管理命令
- 支持将 API Token 写入配置文件并自动鉴权

## Quality Gates

These commands must pass for every user story:
- `go test ./...` - 单元测试
- `go vet ./...` - 静态分析
- `make build` - 构建验证

For whitelist/CLI command stories, also include:
- 运行 `go test -race ./internal/whitelist/...` 验证并发安全
- 运行 `go test -tags=integration -v ./tests/integration/...` 验证集成测试（需要 CloudDrive2 环境）

## User Stories

### US-001: 创建命令元数据注册表

As a developer, I want 一个集中的命令元数据注册表 so that 命令定义、风险等级、白名单默认值不再散落在多个文件。

**Acceptance Criteria:**
- [ ] 新增 `internal/registry/registry.go`，定义 `CommandSpec` 结构体
- [ ] `CommandSpec` 包含 ID、Category、RPC、Description、RiskLevel、DefaultOpen 字段
- [ ] ID 使用稳定小写路径格式如 `file.delete`, `system.restart`
- [ ] 提供 `Register()` 和 `Get()` 方法
- [ ] 提供按 Category/RiskLevel 筛选的查询方法
- [ ] 所有已注册命令可通过 `List()` 获取

### US-002: 修复 whitelist 包并发安全与默认配置

As a developer, I want `internal/whitelist` 包正确处理锁和默认配置 so that 不会出现自锁且默认配置有效。

**Acceptance Criteria:**
- [ ] 修复 `RegisterCommand`、`EnableCommand` 等方法的自锁问题
- [ ] `initDefaultConfig()` 根据 registry 生成默认配置
- [ ] 低/中风险命令默认 enabled=true
- [ ] 高/critical 风险命令默认 enabled=false
- [ ] 配置结构支持 `whitelist_enabled` 和 `commands` 字段
- [ ] 单元测试验证并发安全性 (`go test -race`)

### US-003: 接入白名单到 Cobra 命令执行链路

As a user, I want 远程 API 命令在执行前检查白名单 so that 高风险命令默认被阻止。

**Acceptance Criteria:**
- [ ] 在 `newRootCommand()` 的 `PersistentPreRunE` 中加载白名单配置
- [ ] 每个远程 API 命令执行前调用 `IsAllowed(commandID)`
- [ ] 白名单禁用时所有命令允许执行
- [ ] 白名单启用时未注册命令阻止执行
- [ ] 本地配置命令（whitelist、auth token set）不被白名单阻塞
- [ ] 被阻止时返回包含 commandID 和风险等级的错误信息

### US-004: 实现白名单 CLI 命令

As a user, I want 通过 CLI 管理白名单 so that 可以启用/禁用特定命令。

**Acceptance Criteria:**
- [ ] 新增 `cd2-cli whitelist list [--all] [--risk LEVEL] [--category CAT]` 命令
- [ ] 新增 `cd2-cli whitelist status COMMAND_ID` 命令
- [ ] 新增 `cd2-cli whitelist enable COMMAND_ID` 命令
- [ ] 新增 `cd2-cli whitelist disable COMMAND_ID` 命令
- [ ] 新增 `cd2-cli whitelist reset` 命令（根据 registry 重建默认值）
- [ ] 新增 `cd2-cli whitelist path` 命令（显示配置文件路径）
- [ ] 所有命令输出默认 JSON 格式
- [ ] 配置文件权限设为 `0600`
- [ ] 首次执行任何 whitelist 命令时自动启用白名单功能

### US-005: 实现配置路径迁移与 Token 写入

As a user, I want 配置文件统一存储在 `~/.config/cd2-cli.yaml` 并支持 Token 写入 so that 配置管理更规范。

**Acceptance Criteria:**
- [ ] 配置默认路径改为 `~/.config/cd2-cli.yaml`
- [ ] 新增 `cd2-cli auth token set TOKEN` 命令
- [ ] 新增 `cd2-cli auth token clear` 命令
- [ ] 新增 `cd2-cli auth token show --redact` 命令
- [ ] `auth login --save` 在登录成功后写入 token
- [ ] 写入配置时遵循 `--config` 指定路径
- [ ] 文件权限 `0600`
- [ ] JSON 输出不完整回显 token
- [ ] 继续支持 `--token` 和 `CD2_CLI_TOKEN` 环境变量覆盖

### US-006: 补齐 session 命令组

As a user, I want 通过 CLI 管理 CloudDrive2 会话 so that 可以查看和撤销会话。

**Acceptance Criteria:**
- [ ] 新增 `cd2-cli session list` 命令
- [ ] 新增 `cd2-cli session revoke SESSION_ID` 命令
- [ ] 新增 `cd2-cli session revoke-others` 命令
- [ ] 所有命令在 registry 中注册并设置风险等级
- [ ] `revoke` 和 `revoke-others` 设为高风险默认关闭
- [ ] 集成测试验证命令执行

### US-007: 补齐 token 命令组

As a user, I want 通过 CLI 管理 API Token so that 可以创建、查询、删除 Token。

**Acceptance Criteria:**
- [ ] 新增 `cd2-cli token list` 命令
- [ ] 新增 `cd2-cli token create` 命令
- [ ] 新增 `cd2-cli token modify TOKEN_ID` 命令
- [ ] 新增 `cd2-cli token remove TOKEN_ID` 命令
- [ ] 新增 `cd2-cli token info TOKEN_ID` 命令
- [ ] `create/modify/remove` 设为高风险默认关闭
- [ ] 集成测试验证命令执行

### US-008: 补齐 webdav 命令组

As a user, I want 通过 CLI 管理 WebDAV 用户和服务器配置 so that 可以配置 WebDAV 访问。

**Acceptance Criteria:**
- [ ] 新增 `cd2-cli webdav user get` 命令
- [ ] 新增 `cd2-cli webdav user add` 命令
- [ ] 新增 `cd2-cli webdav user modify` 命令
- [ ] 新增 `cd2-cli webdav user remove` 命令
- [ ] 新增 `cd2-cli webdav server get` 命令
- [ ] 新增 `cd2-cli webdav server set` 命令
- [ ] 用户管理命令设为高风险默认关闭
- [ ] 集成测试验证命令执行

### US-009: 补齐 webhook 命令组

As a user, I want 通过 CLI 管理 Webhook so that 可以配置事件通知。

**Acceptance Criteria:**
- [ ] 新增 `cd2-cli webhook template` 命令
- [ ] 新增 `cd2-cli webhook list` 命令
- [ ] 新增 `cd2-cli webhook add` 命令
- [ ] 新增 `cd2-cli webhook change` 命令
- [ ] 新增 `cd2-cli webhook remove` 命令
- [ ] `add/change/remove` 设为高风险默认关闭
- [ ] 集成测试验证命令执行

### US-010: 补齐 offline 命令组

As a user, I want 通过 CLI 管理离线下载任务 so that 可以添加和管理离线下载。

**Acceptance Criteria:**
- [ ] 新增 `cd2-cli offline add` 命令
- [ ] 新增 `cd2-cli offline remove` 命令
- [ ] 新增 `cd2-cli offline list` 命令
- [ ] 新增 `cd2-cli offline list-all` 命令
- [ ] 新增 `cd2-cli offline quota` 命令
- [ ] 新增 `cd2-cli offline clear` 命令
- [ ] 新增 `cd2-cli offline restart` 命令
- [ ] `add/remove/clear/restart` 设为高风险默认关闭
- [ ] 集成测试验证命令执行

### US-011: 补齐 copy 命令组

As a user, I want 通过 CLI 管理复制任务 so that 可以控制文件复制操作。

**Acceptance Criteria:**
- [ ] 新增 `cd2-cli copy tasks` 命令
- [ ] 新增 `cd2-cli copy merge-tasks` 命令
- [ ] 新增 `cd2-cli copy cancel` 命令
- [ ] 新增 `cd2-cli copy pause` 命令
- [ ] 新增 `cd2-cli copy restart` 命令
- [ ] 新增 `cd2-cli copy remove` 命令
- [ ] 新增 `cd2-cli copy resume` 命令
- [ ] `cancel/pause/restart/remove/resume` 设为高风险默认关闭
- [ ] 集成测试验证命令执行

### US-012: 补齐 sync 命令组

As a user, I want 通过 CLI 同步文件变更 so that 可以监听和处理文件同步。

**Acceptance Criteria:**
- [ ] 新增 `cd2-cli sync file-changes` 命令
- [ ] 新增 `cd2-cli sync start-listener` 命令
- [ ] 新增 `cd2-cli sync stop-listener` 命令
- [ ] 新增 `cd2-cli sync walk-test` 命令
- [ ] 新增 `cd2-cli sync cd1-user-data` 命令
- [ ] `start-listener/stop-listener` 设为高风险默认关闭
- [ ] 集成测试验证命令执行

### US-013: 补齐 local 命令组

As a user, I want 通过 CLI 访问本地文件系统 so that 可以管理本地目录。

**Acceptance Criteria:**
- [ ] 新增 `cd2-cli local list` 命令
- [ ] 新增 `cd2-cli local mkdir` 命令
- [ ] `mkdir` 设为高风险默认关闭
- [ ] 集成测试验证命令执行

### US-014: 补齐 remote-upload 命令组

As a user, I want 通过 CLI 管理远程上传 so that 可以控制和监控远程上传进度。

**Acceptance Criteria:**
- [ ] 新增 `cd2-cli remote-upload start` 命令
- [ ] 新增 `cd2-cli remote-upload control` 命令
- [ ] 新增 `cd2-cli remote-upload read-data` 命令
- [ ] 新增 `cd2-cli remote-upload hash-progress` 命令
- [ ] `start/control` 设为高风险默认关闭
- [ ] 集成测试验证命令执行

### US-015: 补齐 promotion 命令组

As a user, I want 通过 CLI 管理推广计划 so that 可以查询和操作推广相关功能。

**Acceptance Criteria:**
- [ ] 新增查询类 promotion 命令（默认开放）
- [ ] 新增 `cd2-cli promotion join-plan` 命令
- [ ] 新增 `cd2-cli promotion bind-cloud-account` 命令
- [ ] 新增 `cd2-cli promotion transfer-balance` 命令
- [ ] 新增 `cd2-cli promotion activate-plan` 命令
- [ ] 新增 `cd2-cli promotion send-action` 命令
- [ ] 写入操作设为高风险默认关闭
- [ ] 集成测试验证命令执行

### US-016: 补齐 system 命令组剩余接口

As a user, I want 通过 CLI 访问完整的系统管理接口 so that 可以管理缓存、更新、日志等。

**Acceptance Criteria:**
- [ ] 补齐缓存管理命令
- [ ] 补齐更新相关命令
- [ ] 补齐日志相关命令
- [ ] 补齐 WebServer 配置命令
- [ ] 补齐设备管理命令
- [ ] 补齐表查询命令
- [ ] 写入类命令设为高风险默认关闭
- [ ] 集成测试验证命令执行

### US-017: 补齐 cloudapi 命令组剩余接口

As a user, I want 通过 CLI 访问完整的云盘 API 接口 so that 可以管理云盘登录和配置。

**Acceptance Criteria:**
- [ ] 补齐所有登录方式命令
- [ ] 新增发现 SMB 命令
- [ ] 新增设置配置命令
- [ ] 新增能力查询命令
- [ ] 登录/移除/配置命令设为高风险默认关闭
- [ ] 集成测试验证命令执行

### US-018: 补齐 file 命令组剩余接口

As a user, I want 通过 CLI 访问完整的文件操作接口 so that 可以进行批量删除、加密、上传下载等操作。

**Acceptance Criteria:**
- [ ] 新增批量删除命令
- [ ] 新增永久删除命令
- [ ] 新增加密目录命令
- [ ] 新增锁/解锁命令
- [ ] 新增属性查询命令
- [ ] 新增空间查询命令
- [ ] 新增 metadata 命令
- [ ] 新增原始路径命令
- [ ] 新增上传/下载命令
- [ ] 删除/写入命令设为高风险默认关闭
- [ ] 复杂参数支持 `--request-json` 或 `--request-file`
- [ ] 集成测试验证命令执行

### US-019: 补齐 mount 命令组剩余接口

As a user, I want 通过 CLI 访问完整的挂载管理接口 so that 可以更新挂载和查询挂载能力。

**Acceptance Criteria:**
- [ ] 新增 `cd2-cli mount update` 命令
- [ ] 新增 `cd2-cli mount can-add` 命令
- [ ] 新增 `cd2-cli mount drive-letters` 命令
- [ ] 新增 `cd2-cli mount has-drive-letters` 命令
- [ ] 新增 `cd2-cli mount can-mount-both` 命令
- [ ] `update` 设为高风险默认关闭
- [ ] 集成测试验证命令执行

### US-020: 补齐 backup 命令组剩余接口

As a user, I want 通过 CLI 访问完整的备份管理接口 so that 可以配置和管理备份任务。

**Acceptance Criteria:**
- [ ] 新增 `cd2-cli backup add` 命令
- [ ] 新增 `cd2-cli backup update` 命令
- [ ] 新增 `cd2-cli backup destination` 命令
- [ ] 新增 `cd2-cli backup enabled` 命令
- [ ] 新增 `cd2-cli backup watch` 命令
- [ ] 新增 `cd2-cli backup strategies` 命令
- [ ] 新增 `cd2-cli backup can-add` 命令
- [ ] 新增 `cd2-cli backup notify` 命令
- [ ] `add/update` 设为高风险默认关闭
- [ ] 集成测试验证命令执行

### US-021: 补齐 transfer 命令组剩余接口

As a user, I want 通过 CLI 访问完整的传输管理接口 so that 可以控制下载上传任务。

**Acceptance Criteria:**
- [ ] 新增 `cd2-cli transfer download count` 命令
- [ ] 新增 `cd2-cli transfer upload count` 命令
- [ ] 新增按 key 的 cancel/pause/resume 命令
- [ ] cancel/pause/resume 设为高风险默认关闭
- [ ] 集成测试验证命令执行

### US-022: 补齐 auth/public 命令组剩余接口

As a user, I want 通过 CLI 访问完整的认证接口 so that 可以管理 2FA、邮箱确认、第三方登录等。

**Acceptance Criteria:**
- [ ] 补齐 2FA 相关命令
- [ ] 补齐邮箱确认命令
- [ ] 补齐第三方登录命令
- [ ] 补齐注册/重置命令
- [ ] 账号修改类命令设为高风险默认关闭
- [ ] 集成测试验证命令执行

### US-023: 命令注册完整性验证测试

As a developer, I want 自动验证所有 Cobra 命令都在 registry 中注册 so that 不会遗漏白名单覆盖。

**Acceptance Criteria:**
- [ ] 测试验证所有 Cobra 远程 API 命令都有对应 `CommandSpec`
- [ ] 测试验证所有 `CommandSpec` 都能在 `whitelist list` 中出现
- [ ] 测试作为 `go test ./...` 的一部分运行
- [ ] 测试失败时列出未注册的命令

## Functional Requirements

- FR-1: 所有远程 API 命令必须在 registry 中注册
- FR-2: 高风险命令默认阻止执行，需显式启用
- FR-3: 白名单启用时，未注册命令阻止执行
- FR-4: 本地配置命令（whitelist、auth token）不被白名单阻塞
- FR-5: 配置文件默认路径为 `~/.config/cd2-cli.yaml`
- FR-6: 配置文件权限必须为 `0600`
- FR-7: 首次执行 whitelist 命令时自动启用白名单
- FR-8: 复杂参数命令支持 `--request-json` 或 `--request-file`
- FR-9: 被阻止命令返回包含 commandID 和风险等级的错误信息
- FR-10: JSON 输出不完整回显敏感 token 信息

## Non-Goals

- 不实现配置文件从旧路径自动迁移（新项目无旧用户）
- 不实现系统主题自动检测
- 不实现 Web UI 管理界面
- 不实现远程白名单同步（仅本地管理）
- 不默认运行集成测试（需要 CloudDrive2 环境）

## Technical Considerations

- 使用 `internal/registry` 包管理命令元数据，避免散落
- whitelist 包需使用读写锁，避免自锁问题
- Cobra 命令的 `PersistentPreRunE` 是白名单检查的接入点
- 复用 `internal/client` 现有 wrapper，不在 CLI 中直接调用 protobuf stub
- 集成测试放在 `tests/integration/` 目录

## Success Metrics

- 所有 CloudDrive2 gRPC API 都有对应 CLI 命令
- 所有命令都有风险等级定义
- 单元测试覆盖率 > 80%
- 所有集成测试在 CloudDrive2 环境中通过
- `go test ./...`、`go vet ./...`、`make build` 全部通过

## Open Questions

- 是否需要命令分组启用/禁用功能？
- 是否需要白名单导入/导出功能？
- 是否需要审计日志记录被阻止的命令？