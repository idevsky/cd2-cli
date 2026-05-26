# cd2-cli

CloudDrive2 CLI - 基于 gRPC API 的命令行工具，用于管理远程 CloudDrive2 实例。

## 功能特性

- Go 客户端封装覆盖 CloudDrive2 大部分 gRPC API，CLI 暴露常用管理命令
- 结构化 JSON 输出，适合 AI Agent 直接调用和解析
- 支持 TLS 连接和跳过证书验证
- 配置文件支持（YAML 格式）
- 完整的单元测试覆盖

## 安装

```bash
go build -o cd2-cli ./cmd/cd2-cli
```

## 配置

可以通过命令行参数或配置文件（`~/.config/cd2-cli.yaml`）配置：

```yaml
server: localhost:19798
token: your-jwt-token
timeout: 30s
tls: false
skip-tls-verify: false
json: true
```

## 使用方法

### 系统信息

```bash
# 获取系统信息（无需认证）
cd2-cli system info

# 获取运行时信息（需要认证）
cd2-cli --token YOUR_TOKEN system runtime

# 获取服务能力
cd2-cli --token YOUR_TOKEN system capabilities

# 重启服务
cd2-cli --token YOUR_TOKEN system restart
```

### 文件操作

```bash
# 列出文件
cd2-cli --token YOUR_TOKEN file list /path

# 查找文件
cd2-cli --token YOUR_TOKEN file find /path/to/file

# 创建文件夹
cd2-cli --token YOUR_TOKEN file mkdir /parent/path folder-name

# 删除文件
cd2-cli --token YOUR_TOKEN file delete /path/to/file

# 重命名文件
cd2-cli --token YOUR_TOKEN file rename /old/path new-name

# 移动文件
cd2-cli --token YOUR_TOKEN file move /source1 /source2 /dest

# 复制文件
cd2-cli --token YOUR_TOKEN file copy /source1 /source2 /dest

# 搜索文件
cd2-cli --token YOUR_TOKEN file search query --path / --fuzzy

# 获取下载链接
cd2-cli --token YOUR_TOKEN file download-url /path/to/file
```

### 挂载点管理

```bash
# 列出挂载点
cd2-cli --token YOUR_TOKEN mount list

# 添加挂载点
cd2-cli --token YOUR_TOKEN mount add /source/path /mount/point

# 移除挂载点
cd2-cli --token YOUR_TOKEN mount remove /mount/point

# 启动挂载
cd2-cli --token YOUR_TOKEN mount start /mount/point

# 停止挂载
cd2-cli --token YOUR_TOKEN mount stop /mount/point
```

### 云 API 管理

```bash
# 列出所有云 API
cd2-cli --token YOUR_TOKEN cloudapi list

# 获取云 API 配置
cd2-cli --token YOUR_TOKEN cloudapi config cloud-name user-name

# 移除云 API
cd2-cli --token YOUR_TOKEN cloudapi remove cloud-name user-name
```

### 备份管理

```bash
# 列出备份
cd2-cli --token YOUR_TOKEN backup list

# 获取备份状态
cd2-cli --token YOUR_TOKEN backup status backup-id

# 移除备份
cd2-cli --token YOUR_TOKEN backup remove backup-id

# 重启备份扫描
cd2-cli --token YOUR_TOKEN backup restart backup-id
```

### 传输任务管理

```bash
# 获取任务计数
cd2-cli --token YOUR_TOKEN transfer count

# 列出下载任务
cd2-cli --token YOUR_TOKEN transfer downloads

# 列出上传任务
cd2-cli --token YOUR_TOKEN transfer uploads

# 取消所有上传
cd2-cli --token YOUR_TOKEN transfer cancel-uploads

# 暂停所有上传
cd2-cli --token YOUR_TOKEN transfer pause-uploads

# 恢复所有上传
cd2-cli --token YOUR_TOKEN transfer resume-uploads
```

### 认证管理

```bash
# 登录获取令牌
cd2-cli auth login username password

# 或避免在命令行历史中暴露密码
CD2_CLI_PASSWORD='password' cd2-cli auth login username

# 登出
cd2-cli --token YOUR_TOKEN auth logout

# 获取账户状态
cd2-cli --token YOUR_TOKEN auth status

# 修改密码
cd2-cli --token YOUR_TOKEN auth change-password old-pass new-pass
```

## 命令行参数

| 参数 | 简写 | 说明 |
|------|------|------|
| `--server` | `-s` | CloudDrive2 服务器地址（默认: localhost:19798） |
| `--token` | `-t` | JWT 认证令牌 |
| `--timeout` | | 远程 API 调用超时时间（默认: 30s） |
| `--json` | `-j` | JSON 输出格式（默认: true） |
| `--tls` | | 使用 TLS 连接 |
| `--skip-tls-verify` | | 跳过 TLS 证书验证 |
| `--config` | | 配置文件路径（默认: ~/.config/cd2-cli.yaml） |
| `--whitelist-config` | | 白名单配置文件路径（默认: ~/.cd2-cli-whitelist.yaml） |

## 环境变量

支持 `CD2_CLI_` 前缀的环境变量，用于避免在命令行历史中暴露敏感信息：

| 环境变量 | 说明 |
|---------|------|
| `CD2_CLI_SERVER` | 服务器地址 |
| `CD2_CLI_TOKEN` | 认证令牌 |
| `CD2_CLI_TIMEOUT` | 超时时间 |
| `CD2_CLI_TLS` | 是否使用 TLS |
| `CD2_CLI_SKIP_TLS_VERIFY` | 跳过 TLS 验证 |
| `CD2_CLI_PASSWORD` | auth login 密码 |
| `CD2_CLI_TOTP` | 2FA 验证码 |
| `CD2_CLI_OLD_PASSWORD` | auth change-password 旧密码 |
| `CD2_CLI_NEW_PASSWORD` | auth change-password 新密码 |
| `CD2_CLI_REFRESH_TOKEN` | OAuth 登录刷新令牌 |

## 白名单管理

CLI 支持命令白名单，用于限制可执行的远程 API 命令：

```bash
# 列出所有命令及其状态
cd2-cli whitelist list

# 列出特定类别的命令
cd2-cli whitelist list --category file

# 列出高风险命令
cd2-cli whitelist list --risk high

# 启用命令
cd2-cli whitelist enable file.delete

# 禁用命令
cd2-cli whitelist disable system.restart

# 查看白名单配置路径
cd2-cli whitelist path
```

白名单配置存储在 `~/.cd2-cli-whitelist.yaml`，与主配置文件分离。新安装时，写操作命令默认被禁用，需要显式启用。

## 安全注意事项

- **认证输出敏感**: `auth login` 和 `auth token show` 的 JSON 输出包含认证令牌，请勿记录或分享
- **下载链接敏感**: `file download-url` 的 JSON 输出包含签名 URL，具有时效性但应避免泄露
- **密码文件**: 使用 `--password-file` 或 `CD2_CLI_PASSWORD` 环境变量可避免密码出现在命令行历史
- **文件权限**: 配置文件使用 0600 权限，仅当前用户可读写
- **白名单**: 新安装时写操作默认禁用，显式启用后才可执行

## Docker 部署 CloudDrive2

```bash
docker run -d \
  --name clouddrive2 \
  -p 19798:19798 \
  -v /path/to/data:/CloudNAS \
  -v /path/to/config:/Config \
  --privileged \
  cloudnas/clouddrive2-unstable:latest
```

## 开发

### 项目结构

```
.
├── cmd/cd2-cli/          # CLI 命令实现
│   ├── main.go           # 入口文件
│   └── cmd/              # Cobra 命令定义
├── internal/client/      # gRPC 客户端封装
├── pkg/proto/            # Proto 文件和生成的代码
└── tests/                # 测试文件
```

### 运行测试

```bash
go test ./... -v -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## API 覆盖

`internal/client` 包封装了 CloudDrive2 的主要 gRPC API，包括：

- 公共方法（无需认证）：GetSystemInfo, GetToken, Login, Register 等
- 认证管理：2FA, 会话管理, 密码修改等
- 文件操作：列表、搜索、创建、删除、重命名、移动、复制等
- 挂载点管理：添加、删除、挂载、卸载等
- 云 API 管理：各种网盘登录、配置管理等
- 备份管理：备份创建、配置、状态查询等
- 传输任务：上传下载管理、暂停恢复等
- 系统管理：运行信息、缓存管理、服务控制等

## 许可证

MIT License
