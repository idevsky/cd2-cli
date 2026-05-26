# AGENTS.md

## Project Overview

`cd2-cli` is a Go command line client and library wrapper for the CloudDrive2 gRPC API.

Main areas:
- `cmd/cd2-cli`: Cobra/Viper CLI entrypoint.
- `internal/client`: typed Go wrappers around CloudDrive2 gRPC methods.
- `pkg/proto`: generated protobuf and gRPC bindings.
- `tests/integration`: integration tests that require a real CloudDrive2 environment.
- `tests/test-env.sh` and Docker Compose files: local integration environment helpers.

The CLI emits JSON by default so other tools and agents can parse results reliably.

## Common Commands

Use these commands before handing off changes:

```bash
go test ./...
go vet ./...
make build
```

Optional checks:

```bash
go test -race ./internal/client
go test -tags=integration -v ./tests/integration/...
```

Integration tests require CloudDrive2 and supporting services. Start them with:

```bash
make integration-env-up
make test-integration
make integration-env-down
```

Do not run integration tests as part of ordinary unit-test verification unless the environment is already available.

## Configuration

CLI defaults:
- Server: `localhost:19798`
- Timeout: `30s`
- JSON output: enabled
- TLS: disabled

Supported config sources:
- Flags such as `--server`, `--token`, `--timeout`, `--tls`.
- YAML config file, defaulting to `~/.cd2-cli.yaml`.
- Environment variables with `CD2_CLI_` prefix, for example `CD2_CLI_SERVER`.
- `auth login` accepts `CD2_CLI_PASSWORD` so passwords do not need to appear in shell history.

Never commit tokens, passwords, cloud refresh tokens, cookies, or generated test data.

## Functional Notes

Implemented CLI command groups:
- `system`: public info, runtime info, capabilities, restart.
- `auth`: login, logout, account status, password change.
- `file`: list, find, mkdir, delete, rename, move, copy, search, download URL.
- `mount`: list, add, remove, start, stop.
- `cloudapi`: list, config, remove.
- `backup`: list, status, remove, restart scan.
- `transfer`: counts, download/upload lists, upload control.

The client package exposes many more CloudDrive2 APIs than the CLI surfaces. Prefer adding CLI commands by reusing existing `internal/client` wrappers instead of calling protobuf stubs directly.

## Performance Notes

- Always pass a bounded `context.Context` to network calls. The CLI wraps each command in `--timeout`.
- Avoid collecting unbounded server streams unless the API is known to terminate. Some push APIs can be long-lived.
- Keep generated protobuf files out of manual refactors unless regenerating from `pkg/proto/clouddrive.proto`.

## Security Notes

- Authentication is sent as `authorization: Bearer <token>` metadata.
- `--skip-tls-verify` is available only for explicit insecure environments; do not enable it by default.
- TLS clients enforce TLS 1.2 or newer.
- Treat JSON output from auth and download-url commands as sensitive because it may contain reusable tokens or signed URLs.
- Do not log command arguments that may include passwords.

## Editing Guidance

- Keep changes small and local to the relevant API group.
- Prefer `rg` for searching and `gofmt` after Go edits.
- Use `go test ./...`, `go vet ./...`, and `make build` for normal verification.
- Build output `./cd2-cli` is ignored and should not be committed.
