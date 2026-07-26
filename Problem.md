# Problem.md

本文档基于当前最新提交 `098a68f` 重新审视，旧内容已清空。以下只记录当前仍存在的问题和建议修复方式。

已执行验证：

- `go test ./...` 通过。
- `go vet ./...` 通过。
- `make build` 通过。
- `go test -race ./internal/whitelist ./internal/registry ./cmd/cd2-cli/cmd` 通过。
- 烟测确认：
  - `auth token set` 不再创建 whitelist 文件。
  - `system restart` 默认被 whitelist 阻断。
  - `whitelist enable cloudapi.login-s3` 可以启用 canonical `storage.add.s3`。
  - `local list/mkdir` 不再 panic。

## P0: `local list/mkdir` 被改成操作 CLI 所在机器，丢失 CloudDrive2 Local API 语义

位置：

- `cmd/cd2-cli/cmd/local.go`
- `internal/client/local.go`
- `pkg/proto/clouddrive.proto`

现状：

当前 `local list` 和 `local mkdir` 不再调用 CloudDrive2 gRPC，而是直接调用本机文件系统：

```go
entries, err := os.ReadDir(parentFolder)
err := os.MkdirAll(createdPath, 0755)
```

但项目里仍然存在 CloudDrive2 Local API wrapper：

```go
cd2Client.Local().LocalGetSubFiles(...)
cd2Client.Local().LocalCreateFolder(...)
```

proto 也明确有远端 API：

```proto
rpc LocalGetSubFiles(LocalGetSubFilesRequest) returns (stream LocalGetSubFilesResult) {}
rpc LocalCreateFolder(LocalCreateFolderRequest) returns (LocalCreateFolderResult) {}
```

问题：

- `cd2-cli --server remote:19798 local mkdir /data x` 现在会在运行 CLI 的机器上创建目录，而不是 CloudDrive2 服务端所在机器。
- 这改变了已有命令语义，属于功能回归。
- 这也让 CloudDrive2 的 `LocalGetSubFiles` / `LocalCreateFolder` API 从 CLI 暴露面中消失，违背“CLI 全部开放”的目标。
- 当前 `local.list/local.mkdir` 仍设置 command ID，但标记为 local，不进入 registry/whitelist；这让命令 ID 的语义更混乱。

建议修复：

1. 恢复 `local list/mkdir` 走 CloudDrive2 gRPC。
2. 不要把这两个命令标记为 `markAsLocal`。
3. 在 registry 中恢复 `local.list` 和 `local.mkdir`，其中 `local.mkdir` 应高风险默认关闭。
4. 如果确实需要“CLI 本机文件系统”功能，新增独立命令组，例如 `host-local list/mkdir`，不要复用 CloudDrive2 Local API 的命令名。
5. 增加测试：
   - `local mkdir` 默认被 whitelist 阻断，而不是直接改本机文件系统。
   - 启用白名单后调用 mock gRPC `LocalCreateFolder`。
   - `local list` 调用 mock gRPC `LocalGetSubFiles`，而不是 `os.ReadDir`。

## P1: `completion` 虽不再被阻断，但仍会创建 whitelist 文件

位置：

- `cmd/cd2-cli/cmd/root.go`
- `cmd/cd2-cli/cmd/registry_init.go`

现象：

实测：

```bash
tmp=$(mktemp -d)
./cd2-cli --config "$tmp/config.yaml" --whitelist-config "$tmp/whitelist.yaml" completion bash >/tmp/completion.bash
test -f "$tmp/whitelist.yaml" && echo yes
```

结果为 `yes`。

原因：

`persistentPreRunE` 先执行 `needsWhitelistInit(cmd)` 并初始化 whitelist，然后才执行 `shouldSkipWhitelistCheck(cmd)`：

```go
shouldInitWhitelist := needsWhitelistInit(cmd)
if shouldInitWhitelist && whitelistMgr == nil {
    initWhitelist()
}

if shouldSkipWhitelistCheck(cmd) {
    return nil
}
```

`completion bash` 最终能跳过检查，但在跳过之前已经创建了 whitelist 文件。

影响：

- 纯本地 completion 生成命令仍有配置写入副作用。
- shell 自动补全安装流程可能无意生成 whitelist 配置。
- 这和 `auth token set` 已修复的“本地命令不应创建 whitelist”原则不一致。

建议修复：

1. `needsWhitelistInit` 应对 `completion/help/__complete/__completeScript` 返回 false。
2. 或者在 `persistentPreRunE` 中先判断纯本地 skip，再决定是否初始化 whitelist。
3. 注意 `whitelist` 命令组本身仍需要初始化 whitelist，可通过 `annotationNeedsWhitelist` 保留。
4. 增加测试：
   - `completion bash/zsh/fish/powershell` 不创建 whitelist 文件。
   - `__complete` / `__completeScript` 不创建 whitelist 文件。
   - `whitelist list` 仍会初始化 whitelist。

## P1: `file download` 仍然不会下载文件

位置：

- `cmd/cd2-cli/cmd/file_upload.go`
- `cmd/cd2-cli/cmd/fs_transfer.go`
- `internal/client/file.go`

现状：

顶层命令 `download [remote-path] [local-file]` 已改为调用 `DownloadRemoteFile`。

但子命令 `file download [remote-path] [local-path]` 仍然只获取 URL 并输出：

```go
urlInfo, err := cd2Client.File().GetDownloadUrl(...)
return outputResult(map[string]interface{}{
    "success": true,
    "message": "Direct URL available. Use external tool to download.",
})
```

问题：

- `cd2-cli file download /remote/a.bin /tmp/a.bin` 的文案仍然是下载到本地，但不会写 `/tmp/a.bin`。
- 顶层 `download` 和 `file download` 行为不一致。
- 自动化 Agent 很容易看到 `success: true` 后误认为文件已经落盘。

建议修复：

1. 让 `file download` 复用 `cd2Client.File().DownloadRemoteFile(ctx, remotePath, localPath)`。
2. 或者把 `file download` 改名/改文案为 `file download-url`，并且不要返回代表下载完成的 `success: true`。
3. 增加测试：
   - 顶层 `download` 和 `file download` 都真正创建本地文件。
   - 下载失败时退出码非 0。

## P1: 新的 `DownloadRemoteFile` 不能正确处理非 direct URL，并忽略服务端要求的 HTTP header

位置：

- `internal/client/file.go`
- `pkg/proto/clouddrive.proto`

现状：

`DownloadRemoteFile` 优先使用 `directUrl`，否则使用 `downloadUrlPath`：

```go
downloadURL := urlInfo.DownloadUrlPath
if urlInfo.DirectUrl != nil && *urlInfo.DirectUrl != "" {
    downloadURL = *urlInfo.DirectUrl
}
```

但 proto 对 `downloadUrlPath` 的注释是：

```proto
// path and query part of the download URL with placeholders
// e.g. "/static/{SCHEME}/{HOST}/{PREVIEW}/path/to/file.txt?token=abc123"
```

也就是说 `downloadUrlPath` 不是完整 URL。当前代码直接 `http.NewRequest("GET", downloadURL, nil)`，当没有 `directUrl` 时大概率会得到 `unsupported protocol scheme`。

另外 proto 还提供：

```proto
optional string userAgent = 4;
map<string, string> additionalHeaders = 5;
```

当前实现完全没有使用这些 header。

影响：

- 对没有 direct URL 的云盘，顶层 `download` 会失败。
- 对需要特定 User-Agent 或额外 header 的 direct URL，下载可能 403 或内容错误。
- 输出仍返回 `downloadUrl/directUrl`，可能泄露签名 URL。

建议修复：

1. 如果 `directUrl` 存在，按 direct URL 下载，并设置 `userAgent/additionalHeaders`。
2. 如果只有 `downloadUrlPath`：
   - 根据当前 CloudDrive2 server 地址、TLS 配置、downloadUrlPath 组装完整 URL。
   - 正确替换 `{SCHEME}`、`{HOST}`、`{PREVIEW}` 等占位符。
3. 对 signed URL 默认不要输出，除非用户显式 `--show-url`。
4. 增加 mock HTTP server 测试：
   - direct URL 下载成功。
   - direct URL 需要 User-Agent/header。
   - 只有 downloadUrlPath 时能正确拼接 server URL。

## P1: 上传失败清理远端 file handle 时没有带鉴权，并且 warning 写到 stdout

位置：

- `internal/client/file.go`

现状：

上传失败 cleanup 中直接调用底层 gRPC client：

```go
_, closeErr := a.c.client.CloseFile(closeCtx, &pb.CloseFileRequest{
    FileHandle: fileHandle,
})
```

这个 `closeCtx` 来自 `context.Background()`，没有经过 `a.c.withAuth(closeCtx)`。

如果服务端要求 Bearer token，cleanup 的 `CloseFile` 很可能失败。

同时失败 warning 使用 `fmt.Printf` 写 stdout：

```go
fmt.Printf("warning: failed to close file handle on cleanup: %v\n", closeErr)
```

影响：

- 出错路径下可能仍然无法关闭远端 handle。
- warning 写 stdout 会污染 CLI JSON 输出，破坏自动化解析。
- 这个问题在现有 mock 测试中可能测不出来，因为 mock server 未校验 auth metadata。

建议修复：

1. cleanup 使用 `a.c.withAuth(closeCtx)`，或者调用已有 wrapper `a.CloseFile(closeCtx, ...)`。
2. warning 写 stderr，或者不要直接打印，改为把 cleanup error 包装到返回错误中。
3. 增加测试：
   - mock server 要求 metadata 中存在 `authorization`。
   - cleanup 失败时 stdout 仍是合法 JSON，stderr 才允许出现 warning。

## P2: registry 中仍有 7 个无效 RPC 名

位置：

- `cmd/cd2-cli/cmd/registry_init.go`
- `pkg/proto/clouddrive.proto`

扫描结果：

```text
file.upload: UploadLocalFile
file.download: DownloadFile
fs.upload: UploadLocalFile
fs.download: DownloadFile
storage.add: APILogin
cache.prefetch-start: StartFilePrefetch
cloudapi.login-123pan-oauth: ApiLogin123PanOAuth
```

这些不是 proto 中真实的 `rpc` 名。部分是 client wrapper 方法名，部分是聚合/虚拟名。

影响：

- registry 不能可靠用于 API 覆盖统计、审计或文档生成。
- 其他 Agent 会误以为这些是 CloudDrive2 gRPC 方法名。
- “CLI 全部开放”无法用 registry 和 proto 做准确比对。

建议修复：

1. `RPC` 字段只填写 proto 中真实 rpc 名。
2. 对一个 CLI 命令调用多个 RPC 的情况，改成 `RPCs []string`。
3. 若需要记录 wrapper 方法，新增 `Wrapper` 字段，不要混用 `RPC`。
4. 增加测试：解析 `pkg/proto/clouddrive.proto`，要求 registry 中所有 `RPC/RPCs` 都存在于 proto rpc 集合。

## P2: registry 覆盖测试被削弱，local API 消失没有被测出来

位置：

- `cmd/cd2-cli/cmd/registration_test.go`

现状：

当前 `collectCommandIDs` 会排除所有 `annotationIsLocal` 命令：

```go
if _, isLocal := cmd.Annotations[annotationIsLocal]; !isLocal {
    ...
}
```

然后 `TestCobraRemoteCommandsHaveRegistryEntries` 只检查非 local 命令是否在 registry 中。

这导致：

- `local.list/local.mkdir` 被标记为 local 后，从 registry 覆盖测试中消失。
- CloudDrive2 proto 中仍存在 `LocalGetSubFiles/LocalCreateFolder`，但 CLI 不再暴露远端 API，测试没有失败。
- 之前的 `Implemented` 相关测试被删除后，没有新增等价的“proto API 覆盖”或“wrapper API 覆盖”测试。

建议修复：

1. 增加 proto/client wrapper 到 CLI command 的覆盖测试。
2. 对于不暴露的 API，必须有显式 allowlist 和理由。
3. local 命令如果是 CloudDrive2 Local API，就不能标记为 local；如果是 CLI 本机操作，则不应占用 `local.*` 这组 CloudDrive2 API command ID。

## P2: `Register` 仍保留静默覆盖接口，后续代码可能绕过 `MustRegister`

位置：

- `internal/registry/registry.go`

现状：

新增了 `MustRegister`，当前 `registry_init.go` 已使用它。但原来的 `Register` 仍然保留，并且仍然静默覆盖重复 ID：

```go
func (r *Registry) Register(spec *CommandSpec) {
    ...
    r.commands[spec.ID] = spec
}
```

影响：

- 后续新代码或测试可能继续使用 `Register`，重复 ID 又会静默覆盖。
- 当前重复 ID 防线依赖开发者记得使用 `MustRegister`。

建议修复：

1. 让 `Register` 也拒绝重复 ID，返回 error。
2. 若要保留测试便利接口，改名为 `UpsertForTest`，并限制在测试文件中使用。
3. 增加测试：普通 registry duplicate 必须失败或 panic。

## 建议优先级

1. 先修 `local list/mkdir` 语义回归，这是最容易造成误操作的行为问题。
2. 修 completion 创建 whitelist 的副作用。
3. 统一 `file download` 与顶层 `download`，并完善 `DownloadRemoteFile` 对 `downloadUrlPath`、headers、signed URL 输出的处理。
4. 修上传失败 cleanup 的鉴权和 stdout 污染。
5. 最后补 registry RPC 有效性和 API 覆盖测试。

