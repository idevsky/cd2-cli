# PRD: cd2-cli 集成测试环境配置与验证

## Overview
配置本地 Docker CloudDrive2 实例环境，确保集成测试可以正常运行，并完成所有 gRPC API 模块的集成验证。单元测试已完成（覆盖率 90.9%），本项目仅关注集成测试部分。

## Goals
- 解决 Docker 容器 unhealthy 状态，确保 gRPC 服务可用
- 配置集成测试环境变量和连接参数
- 完成所有模块的集成测试验证（fs、storage、task）
- 确保集成测试覆盖所有 gRPC API 调用场景

## Quality Gates

These commands must pass for every user story:
- `go build ./...` - 编译通过
- `go test ./tests/integration/... -v` - 集成测试通过

For Docker stories, also include:
- `docker ps` shows healthy container status
- `grpcurl` can connect to CloudDrive2 gRPC endpoint

## User Stories

### US-001: 诊断并修复 Docker 容器状态
**Description:** As a developer, I want the CloudDrive2 Docker container to be healthy so that integration tests can connect to it.

**Acceptance Criteria:**
- [ ] Identify root cause of unhealthy status
- [ ] Fix container configuration or restart service
- [ ] `docker ps` shows healthy status
- [ ] gRPC port 19798 is accessible from host

### US-002: 配置集成测试环境变量
**Description:** As a developer, I want integration tests to use correct connection parameters so they can connect to the CloudDrive2 instance.

**Acceptance Criteria:**
- [ ] Define environment variables for gRPC endpoint (CD2_HOST, CD2_PORT)
- [ ] Add authentication credentials configuration (CD2_USER, CD2_PASS)
- [ ] Create test helper to load environment configuration
- [ ] Document required environment variables in tests/README.md

### US-003: 验证文件系统模块集成测试
**Description:** As a developer, I want fs_* API integration tests to pass so that file operations are verified against real CloudDrive2 instance.

**Acceptance Criteria:**
- [ ] Test fs_list_files with valid path
- [ ] Test fs_get_file_info for existing files
- [ ] Test fs_mkdir and fs_remove operations
- [ ] Test fs_copy and fs_move operations
- [ ] Test fs_rename operation
- [ ] All tests in auth_test.go and related fs tests pass

### US-004: 验证存储管理模块集成测试
**Description:** As a developer, I want storage_* API integration tests to pass so that storage operations are verified against real CloudDrive2 instance.

**Acceptance Criteria:**
- [ ] Test storage_list to enumerate mounted storages
- [ ] Test storage_add for adding storage mount
- [ ] Test storage_remove for removing storage mount
- [ ] Test storage_update for updating storage configuration
- [ ] All storage-related integration tests pass

### US-005: 验证任务管理模块集成测试
**Description:** As a developer, I want task_* API integration tests to pass so that task operations are verified against real CloudDrive2 instance.

**Acceptance Criteria:**
- [ ] Test task_list to enumerate running tasks
- [ ] Test task_create for creating new tasks
- [ ] Test task_cancel for canceling running tasks
- [ ] Test task_status for checking task progress
- [ ] All task-related integration tests pass

### US-006: 创建集成测试运行脚本
**Description:** As a developer, I want a single command to run all integration tests so that verification is automated.

**Acceptance Criteria:**
- [ ] Create scripts/run-integration.sh script
- [ ] Script checks Docker container status first
- [ ] Script sets required environment variables
- [ ] Script runs all integration tests with verbose output
- [ ] Script reports pass/fail summary

## Functional Requirements
- FR-1: Docker 容器必须以 healthy 状态运行
- FR-2: gRPC 端口 19798 必须可以从主机访问
- FR-3: 集成测试必须能够连接到 CloudDrive2 实例
- FR-4: 所有集成测试必须在 60 秒内完成或超时
- FR-5: 测试失败时必须提供清晰的错误信息

## Non-Goals
- 单元测试补充（已完成）
- 性能测试
- 压力测试
- 多实例负载测试

## Technical Considerations
- 使用 docker inspect 检查容器健康状态
- 使用 grpcurl 或 grpc_health_probe 验证 gRPC 服务可用性
- 集成测试应该使用 testify/suite 管理测试生命周期
- 考虑使用 testcontainers 如果需要动态容器管理
- 使用 .env 文件管理本地开发环境变量

## Success Metrics
- Docker 容器状态显示 healthy
- 所有集成测试通过（go test ./tests/integration/... -v）
- 集成测试覆盖 fs、storage、task 所有主要 API
- 可以通过单一脚本命令运行所有集成测试

## Open Questions
- CloudDrive2 的默认认证凭据是什么？
- 是否需要预先配置网盘账号才能运行某些测试？
- 测试数据应该创建在哪个目录路径？