# Todo: CLI API 全开放与白名单治理

## 当前状态

✅ **已完成** - 所有核心功能已实现并测试：

### 白名单系统 (已完成)
- ✅ 命令元数据注册表 (`internal/registry/registry.go`)
- ✅ 白名单管理器 (`internal/whitelist/whitelist.go`)
- ✅ 白名单 CLI 命令 (`whitelist list/status/enable/disable/reset/path`)
- ✅ 所有远程 API 命令纳入白名单管理
- ✅ 高风险/写操作命令默认禁用 (`DefaultOpen: false`)
- ✅ 别名组机制（fs.* 与 file.* 共享策略）
- ✅ 通配符拒绝验证
- ✅ MergeRegistryDefaults（升级时保留用户设置）

### 配置管理 (已完成)
- ✅ 主配置文件分离 (`~/.config/cd2-cli.yaml`)
- ✅ 白名单配置文件分离 (`~/.cd2-cli-whitelist.yaml`)
- ✅ Token 写入/清除/显示 (`auth token set/clear/show`)
- ✅ 配置保留现有字段
- ✅ AtomicWrite 文件写入

### CLI 增强 (已完成)
- ✅ 全局 `--timeout` 标志 (30s 默认)
- ✅ 环境变量支持 (`CD2_CLI_*`)
- ✅ 敏感参数安全输入 (`CD2_CLI_PASSWORD`, `--password-file`)
- ✅ protojson 输入/输出处理
- ✅ 命令返回非零退出码
- ✅ 本地命令跳过 gRPC 初始化

### 测试覆盖 (已完成)
- ✅ 白名单单元测试
- ✅ Cobra 命令测试
- ✅ 配置写入测试
- ✅ 命令注册双向验证测试

## 目标（已达成）

✅ 所有核心目标已完成：
1. CLI 暴露 CloudDrive2 gRPC API 的核心可用面（常用命令已实现）
2. 所有远程 API 命令都纳入本地白名单管理
3. 高风险/写操作 API 默认禁用
4. 白名单管理命令已实现
5. Token 写入配置文件已实现

## 原实现计划（已归档）

以下内容为原始规划，现已完成实现，保留作为设计参考：

<details>
<summary>点击展开原始规划文档</summary>

### 1. 命令元数据注册表（已完成）

实现了 `internal/registry/registry.go`，包含：
- CommandSpec 结构体（ID, Category, RPC, Description, RiskLevel, DefaultOpen, AliasGroup, Implemented）
- 命令注册函数
- ResolveAliasGroup 别名解析

### 2. 白名单系统（已完成）

实现了 `internal/whitelist/whitelist.go`，包含：
- Manager 管理器
- IsAllowed 检查
- Enable/Disable/Register 方法
- MergeRegistryDefaults 升级合并
- Config.Enabled 使用 `*bool` 区分未设置与 false

### 3. 白名单 CLI（已完成）

实现了 `cmd/cd2-cli/cmd/whitelist.go`，包含：
- `whitelist list [--all] [--risk] [--category]`
- `whitelist status COMMAND_ID`
- `whitelist enable/disable COMMAND_ID`
- `whitelist reset`
- `whitelist path`

### 4. 配置管理（已完成）

实现了 `cmd/cd2-cli/cmd/auth_token.go` 和 `cmd/cd2-cli/cmd/auth.go`：
- `auth token set/clear/show`
- `auth login` 自动保存 token
- 配置文件原子写入（atomicWriteFile）
- 配置保留现有字段（map[string]interface{}）

### 5-7. 其他功能（已完成）

- 全局 `--timeout` 标志
- 环境变量支持
- protojson 输入/输出
- 命令返回非零退出码
- 本地命令跳过 gRPC
- 通配符拒绝
- 别名组共享策略
- DefaultOpen 覆盖风险等级
- 完整测试覆盖

</details>

## 未来扩展方向

以下功能在原始规划中提及但未实现，可根据需求扩展：

- **完整 API 覆盖**: 补齐 session/token/webdav/webhook/offline/copy/sync/local/remote-upload 等更多 API
- **批量操作**: file 批量删除、永久删除、加密目录等
- **系统管理**: 缓存管理、更新管理、日志配置等
- **云盘登录**: 更多 OAuth 登录方式、SMB 发现等

## 测试要求

交付前运行：

```bash
go test ./...
go vet ./...
make build
```

不要默认跑 integration tests，除非 CloudDrive2 环境已经可用。
