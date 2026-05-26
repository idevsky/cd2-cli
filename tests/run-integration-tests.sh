#!/bin/bash
# Integration test runner script
# Runs integration tests against a CloudDrive2 instance and outputs JUnit format

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
JUNIT_OUTPUT="${JUNIT_OUTPUT:-junit.xml}"
TEST_RESULTS_DIR="${TEST_RESULTS_DIR:-$PROJECT_ROOT/test-results}"

usage() {
    echo "Usage: $0 [options] [test-pattern]"
    echo ""
    echo "Options:"
    echo "  --junit         Output results in JUnit XML format (default)"
    echo "  --verbose       Enable verbose test output"
    echo "  --env           Start test environment before running tests"
    echo "  --stop-env      Stop test environment after running tests"
    echo "  --help          Show this help message"
    echo ""
    echo "Environment Variables:"
    echo "  CD2_HOST       CloudDrive2 hostname (default: localhost)"
    echo "  CD2_PORT       CloudDrive2 port (default: 19798)"
    echo "  CD2_USER       CloudDrive2 username (default: admin)"
    echo "  CD2_PASS       CloudDrive2 password (default: admin123)"
    echo "  MINIO_ADDR     MinIO address (default: localhost:9000)"
    echo "  JUNIT_OUTPUT   JUnit output file (default: junit.xml)"
    echo ""
    echo "Examples:"
    echo "  $0                           # Run all integration tests"
    echo "  $0 --env                     # Start test environment and run tests"
    echo "  $0 --env --stop-env          # Start env, run tests, stop env"
    echo "  $0 TestIntegrationAuth       # Run specific test"
    echo "  $0 --verbose TestIntegrationFS  # Run FS tests with verbose output"
}

JUNIT_MODE=false
VERBOSE_MODE=false
START_ENV=false
STOP_ENV=false
TEST_PATTERN=""

while [[ $# -gt 0 ]]; do
    case $1 in
        --junit)
            JUNIT_MODE=true
            shift
            ;;
        --verbose)
            VERBOSE_MODE=true
            shift
            ;;
        --env)
            START_ENV=true
            shift
            ;;
        --stop-env)
            STOP_ENV=true
            shift
            ;;
        --help|-h)
            usage
            exit 0
            ;;
        -*)
            echo "Unknown option: $1"
            usage
            exit 1
            ;;
        *)
            TEST_PATTERN="$1"
            shift
            ;;
    esac
done

cd "$PROJECT_ROOT"

# Check if go-junit-report is installed for JUnit output
if [ "$JUNIT_MODE" = true ] || [ -z "$TEST_PATTERN" ]; then
    if ! command -v go-junit-report &> /dev/null; then
        echo "Installing go-junit-report for JUnit output..."
        go install github.com/jstemmer/go-junit-report/v2@latest
    fi
fi

# Start test environment if requested
if [ "$START_ENV" = true ]; then
    echo "Starting test environment..."
    "$SCRIPT_DIR/test-env.sh" start
    
    # Wait for services to be ready
    echo "Waiting for services to be ready..."
    sleep 10
    
    for i in {1..30}; do
        if curl -sf http://localhost:${CD2_PORT:-19798}/api/GetSystemInfo > /dev/null 2>&1; then
            echo "CloudDrive2 is ready!"
            break
        fi
        echo "Waiting for CloudDrive2... ($i/30)"
        sleep 2
    done
fi

# Prepare test arguments
TEST_ARGS=(-tags=integration -v)
if [ -n "$TEST_PATTERN" ]; then
    TEST_ARGS+=(-run "$TEST_PATTERN")
fi

# Create test results directory
mkdir -p "$TEST_RESULTS_DIR"

# Run tests
echo "Running integration tests..."
echo "Test arguments: ${TEST_ARGS[*]}"

if [ "$VERBOSE_MODE" = true ]; then
    # Verbose mode with JUnit output
    go test "${TEST_ARGS[@]}" ./tests/integration/... 2>&1 | tee "$TEST_RESULTS_DIR/test-output.txt"
    TEST_EXIT_CODE=${PIPESTATUS[0]}
    
    # Convert to JUnit format
    if [ "$JUNIT_MODE" = true ] || [ -z "$TEST_PATTERN" ]; then
        cat "$TEST_RESULTS_DIR/test-output.txt" | go-junit-report -set-exit-code > "$TEST_RESULTS_DIR/$JUNIT_OUTPUT"
        echo "JUnit report saved to: $TEST_RESULTS_DIR/$JUNIT_OUTPUT"
    fi
else
    # Normal mode with JUnit output
    go test "${TEST_ARGS[@]}" ./tests/integration/... 2>&1 | tee "$TEST_RESULTS_DIR/test-output.txt" | go-junit-report -set-exit-code > "$TEST_RESULTS_DIR/$JUNIT_OUTPUT"
    TEST_EXIT_CODE=${PIPESTATUS[0]}
    echo "JUnit report saved to: $TEST_RESULTS_DIR/$JUNIT_OUTPUT"
fi

# Stop test environment if requested
if [ "$STOP_ENV" = true ]; then
    echo "Stopping test environment..."
    "$SCRIPT_DIR/test-env.sh" stop
fi

# Summary
echo ""
echo "================================"
echo "Test Summary"
echo "================================"
if [ $TEST_EXIT_CODE -eq 0 ]; then
    echo "Status: PASSED"
else
    echo "Status: FAILED"
fi
echo "JUnit Report: $TEST_RESULTS_DIR/$JUNIT_OUTPUT"
echo "Test Output: $TEST_RESULTS_DIR/test-output.txt"
echo ""

exit $TEST_EXIT_CODE