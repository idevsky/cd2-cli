# Integration Tests

This directory contains integration tests for cd2-cli that run against a real CloudDrive2 instance.

## Test Structure

- `integration/` - Integration test files organized by functionality
  - `auth_test.go` - Authentication and login tests
  - `filesystem_test.go` - File system operations tests
  - `mount_test.go` - Mount point management tests
  - `storage_test.go` - Storage provider tests (S3, WebDAV, etc.)
  - `task_test.go` - Task management tests (upload, download, copy)
  - `transfer_test.go` - File transfer tests
  - `system_test.go` - System information and management tests
  - `misc_test.go` - Miscellaneous tests (2FA, backup, sessions, etc.)

## Test Environment

The test environment is managed via Docker Compose and includes:

- **CloudDrive2** - Main application instance (port 19798)
- **MinIO** - S3-compatible storage for testing S3 providers (port 9000)
- **WebDAV** - WebDAV server for testing WebDAV providers (port 8080)
- **SFTP** - SFTP server for testing SFTP providers (port 2222)

### Environment Management

Use `test-env.sh` to manage the test environment:

```bash
# Start test environment
./tests/test-env.sh start

# Stop test environment
./tests/test-env.sh stop

# Restart test environment
./tests/test-env.sh restart

# Show environment status
./tests/test-env.sh status

# View container logs
./tests/test-env.sh logs

# Clean up environment (removes all data)
./tests/test-env.sh clean
```

### Environment Variables

The integration tests use the following environment variables for connection configuration:

| Variable | Description | Default |
|----------|-------------|---------|
| `CD2_HOST` | CloudDrive2 hostname or IP | `localhost` |
| `CD2_PORT` | CloudDrive2 gRPC port | `19798` |
| `CD2_USER` | CloudDrive2 username | `admin` |
| `CD2_PASS` | CloudDrive2 password | `admin123` |
| `MINIO_PORT` | MinIO API port | `9000` |
| `MINIO_CONSOLE_PORT` | MinIO console port | `9001` |
| `WEBDAV_PORT` | WebDAV port | `8080` |
| `SFTP_PORT` | SFTP port | `2222` |

#### Setting Environment Variables

```bash
# Example: Set custom connection parameters
export CD2_HOST=192.168.1.100
export CD2_PORT=19798
export CD2_USER=myuser
export CD2_PASS=mypassword

# Then run tests
go test -tags=integration -v ./tests/integration/...
```

Or inline:
```bash
CD2_HOST=192.168.1.100 CD2_USER=myuser CD2_PASS=mypassword go test -tags=integration -v ./tests/integration/...
```

## Running Tests

### Quick Start

```bash
# Start test environment
make integration-env-up

# Run all integration tests
make test-integration

# Stop test environment
make integration-env-down
```

### Using Test Runner Script

The `run-integration-tests.sh` script provides more control:

```bash
# Run all tests with JUnit output
./tests/run-integration-tests.sh

# Start environment, run tests, stop environment
./tests/run-integration-tests.sh --env --stop-env

# Run specific tests
./tests/run-integration-tests.sh TestIntegrationAuth

# Run with verbose output
./tests/run-integration-tests.sh --verbose TestIntegrationFS

# Show help
./tests/run-integration-tests.sh --help
```

### Direct Go Test Commands

```bash
# Run all integration tests
go test -tags=integration -v ./tests/integration/...

# Run specific test
go test -tags=integration -v -run TestIntegrationAuth ./tests/integration/...

# Run tests matching a pattern
go test -tags=integration -v -run TestIntegrationFS ./tests/integration/...
```

### Test Variables

- `CD2_HOST` - CloudDrive2 hostname (default: localhost)
- `CD2_PORT` - CloudDrive2 gRPC port (default: 19798)
- `CD2_USER` - CloudDrive2 username (default: admin)
- `CD2_PASS` - CloudDrive2 password (default: admin123)
- `MINIO_ADDR` - MinIO address for S3 tests (default: localhost:9000)

## JUnit Output

Integration tests can output results in JUnit XML format for CI/CD integration:

```bash
# Generate JUnit report
make test-integration-junit

# JUnit output location
test-results/junit.xml
```

The JUnit report is generated using `go-junit-report` and is suitable for:
- Jenkins
- GitLab CI
- GitHub Actions
- CircleCI
- TeamCity

## Makefile Targets

```
make test-integration          # Run integration tests
make test-integration-junit    # Run tests with JUnit output
make integration-env-up        # Start test environment
make integration-env-down      # Stop test environment
make integration-env-clean     # Clean up test environment
make integration-env-status    # Show environment status
make integration-env-logs      # View container logs
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Integration Tests

on: [push, pull_request]

jobs:
  integration-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.26'
      
      - name: Start test environment
        run: make integration-env-up
      
      - name: Wait for services
        run: sleep 30
      
      - name: Run integration tests
        run: make test-integration-junit
      
      - name: Upload test results
        uses: actions/upload-artifact@v3
        with:
          name: junit-results
          path: test-results/junit.xml
      
      - name: Stop test environment
        run: make integration-env-down
```

### GitLab CI Example

```yaml
integration-tests:
  stage: test
  image: golang:1.26
  services:
    - docker:dind
  script:
    - make integration-env-up
    - sleep 30
    - make test-integration-junit
  artifacts:
    reports:
      junit: test-results/junit.xml
    paths:
      - test-results/
  after_script:
    - make integration-env-down
```

## Test Coverage

Integration tests cover:

1. **Authentication** - Login, logout, status, 2FA
2. **File System** - List, stat, mkdir, upload, download, remove
3. **Mount Points** - List, add, remove, status
4. **Storage Providers** - List, add, remove (S3, WebDAV, SFTP)
5. **Task Management** - Upload/download/copy task listing and status
6. **Transfers** - Upload URL, download URL, task management
7. **System** - Info, runtime, settings, capabilities, logs
8. **Miscellaneous** - Backup, sessions, tokens, cache, webhooks

## Test Scenarios

### Typical Scenarios Covered

1. **User Authentication Flow**
   - Login with username/password
   - Get account status
   - Logout

2. **File Operations**
   - List directory contents
   - Get file information
   - Create folder
   - Upload file
   - Delete file/folder

3. **Mount Point Management**
   - List mount points
   - Add mount point
   - Remove mount point

4. **Storage Provider Management**
   - List all providers
   - Add S3 provider (MinIO)
   - Remove provider

5. **Task Monitoring**
   - List upload tasks
   - List download tasks
   - List copy tasks
   - Get task counts

6. **System Information**
   - Get system info
   - Get runtime info
   - Get settings
   - Get capabilities
   - Check for updates

## Troubleshooting

### Tests Fail to Connect

1. Ensure test environment is running: `make integration-env-status`
2. Check logs: `make integration-env-logs`
3. Verify CloudDrive2 is ready: `curl http://localhost:19798/api/GetSystemInfo`

### Docker Issues

1. Clean up environment: `make integration-env-clean`
2. Restart environment: `make integration-env-up` (after clean)
3. Check Docker logs: `docker logs cd2-test`

### Test Timeout

Increase timeout in test files or environment variables:
```bash
export CD2_TEST_ADDR=localhost:19798
go test -tags=integration -v -timeout 5m ./tests/integration/...
```

## Best Practices

1. **Run tests against clean environment**
   ```bash
   make integration-env-clean && make integration-env-up
   ```

2. **Use JUnit output for CI/CD**
   ```bash
   make test-integration-junit
   ```

3. **Run specific test suites**
   ```bash
   ./tests/run-integration-tests.sh TestIntegrationAuth
   ```

4. **Check environment status before running**
   ```bash
   make integration-env-status
   ```

5. **Clean up after tests**
   ```bash
   make integration-env-down
   ```