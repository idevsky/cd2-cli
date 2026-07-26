#!/bin/bash
# Run integration tests for cd2-cli
# This script checks Docker container status, sets environment variables,
# runs all integration tests with verbose output, and reports pass/fail summary.

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
TESTS_DIR="$PROJECT_ROOT/tests"
COMPOSE_FILE="$TESTS_DIR/docker-compose.test.yml"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

usage() {
    echo "Usage: $0 [options]"
    echo ""
    echo "Run integration tests for cd2-cli against CloudDrive2 instance."
    echo ""
    echo "Options:"
    echo "  --start-env     Start Docker test environment before running tests"
    echo "  --stop-env      Stop Docker test environment after running tests"
    echo "  --keep-env      Keep environment running after tests (don't stop on failure)"
    echo "  --no-check      Skip Docker container status check"
    echo "  -v, --verbose   Enable extra verbose output"
    echo "  -h, --help      Show this help message"
    echo ""
    echo "Environment Variables (with defaults):"
    echo "  CD2_HOST        CloudDrive2 hostname (default: localhost)"
    echo "  CD2_PORT        CloudDrive2 gRPC port (default: 19798)"
    echo "  CD2_USER        CloudDrive2 username (default: admin)"
    echo "  CD2_PASS        CloudDrive2 password (default: admin123)"
    echo ""
    echo "Examples:"
    echo "  $0                         # Run tests against running container"
    echo "  $0 --start-env             # Start env, run tests"
    echo "  $0 --start-env --stop-env  # Start env, run tests, stop env"
    echo "  $0 --no-check              # Skip container check"
}

START_ENV=false
STOP_ENV=false
KEEP_ENV=false
CHECK_CONTAINER=true
EXTRA_VERBOSE=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --start-env)
            START_ENV=true
            shift
            ;;
        --stop-env)
            STOP_ENV=true
            shift
            ;;
        --keep-env)
            KEEP_ENV=true
            shift
            ;;
        --no-check)
            CHECK_CONTAINER=false
            shift
            ;;
        -v|--verbose)
            EXTRA_VERBOSE=true
            shift
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            usage
            exit 1
            ;;
    esac
done

# Set default environment variables
export CD2_HOST="${CD2_HOST:-localhost}"
export CD2_PORT="${CD2_PORT:-19798}"
export CD2_USER="${CD2_USER:-admin}"
export CD2_PASS="${CD2_PASS:-admin123}"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}cd2-cli Integration Test Runner${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo "Configuration:"
echo "  CD2_HOST:  $CD2_HOST"
echo "  CD2_PORT:  $CD2_PORT"
echo "  CD2_USER:  $CD2_USER"
echo "  CD2_PASS:  ****"
echo ""

cd "$PROJECT_ROOT"

# Function to check Docker container status
check_container_status() {
    echo -e "${YELLOW}Checking Docker container status...${NC}"
    
    if ! command -v docker &> /dev/null; then
        echo -e "${RED}Error: Docker is not installed${NC}"
        exit 1
    fi
    
    # Check if container exists and is running
    CONTAINER_NAME="cd2-test"
    CONTAINER_STATUS=$(docker ps -a --filter "name=$CONTAINER_NAME" --format "{{.Status}}" 2>/dev/null || echo "")
    
    if [ -z "$CONTAINER_STATUS" ]; then
        echo -e "${YELLOW}Container '$CONTAINER_NAME' not found.${NC}"
        if [ "$START_ENV" = true ]; then
            echo -e "${YELLOW}Will start environment...${NC}"
        else
            echo -e "${RED}Error: Container is not running. Use --start-env to start it.${NC}"
            exit 1
        fi
    elif echo "$CONTAINER_STATUS" | grep -q "(healthy)"; then
        echo -e "${GREEN}Container '$CONTAINER_NAME' is healthy${NC}"
    elif echo "$CONTAINER_STATUS" | grep -q "Up"; then
        echo -e "${YELLOW}Container '$CONTAINER_NAME' is running but not yet healthy${NC}"
        echo "Status: $CONTAINER_STATUS"
    else
        echo -e "${RED}Container '$CONTAINER_NAME' is not running${NC}"
        echo "Status: $CONTAINER_STATUS"
        if [ "$START_ENV" = true ]; then
            echo -e "${YELLOW}Will restart environment...${NC}"
        else
            exit 1
        fi
    fi
    
    # Check gRPC port connectivity
    echo -e "${YELLOW}Checking gRPC port $CD2_PORT...${NC}"
    if command -v nc &> /dev/null; then
        if nc -z "$CD2_HOST" "$CD2_PORT" 2>/dev/null; then
            echo -e "${GREEN}Port $CD2_PORT is accessible${NC}"
        else
            echo -e "${YELLOW}Warning: Port $CD2_PORT may not be accessible${NC}"
        fi
    elif command -v curl &> /dev/null; then
        if curl -sf "http://$CD2_HOST:$CD2_PORT/" > /dev/null 2>&1; then
            echo -e "${GREEN}HTTP endpoint on port $CD2_PORT is accessible${NC}"
        else
            echo -e "${YELLOW}Warning: Could not reach HTTP endpoint on port $CD2_PORT${NC}"
        fi
    fi
    
    echo ""
}

# Start test environment if requested
if [ "$START_ENV" = true ]; then
    echo -e "${BLUE}Starting test environment...${NC}"
    if [ -f "$TESTS_DIR/test-env.sh" ]; then
        "$TESTS_DIR/test-env.sh" start
    else
        echo -e "${YELLOW}Using docker-compose directly...${NC}"
        docker compose -f "$COMPOSE_FILE" up -d
        
        echo "Waiting for services to be ready..."
        sleep 5
        
        for i in {1..30}; do
            if curl -sf "http://localhost:$CD2_PORT/" > /dev/null 2>&1; then
                echo -e "${GREEN}CloudDrive2 is ready!${NC}"
                break
            fi
            echo "Waiting for CloudDrive2... ($i/30)"
            sleep 2
        done
    fi
    echo ""
fi

# Check container status
if [ "$CHECK_CONTAINER" = true ]; then
    check_container_status
fi

# Verify go build passes
echo -e "${BLUE}Building project...${NC}"
if go build ./...; then
    echo -e "${GREEN}Build successful${NC}"
else
    echo -e "${RED}Build failed${NC}"
    exit 1
fi
echo ""

# Run integration tests
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Running Integration Tests${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

TEST_ARGS=(-tags=integration -v)
if [ "$EXTRA_VERBOSE" = true ]; then
    TEST_ARGS=("${TEST_ARGS[@]}" -v)
fi

# Track start time
START_TIME=$(date +%s)

# Run tests and capture output
TEST_OUTPUT_FILE=$(mktemp)
TEST_FAILED=false

set +e
go test "${TEST_ARGS[@]}" ./tests/integration/... 2>&1 | tee "$TEST_OUTPUT_FILE"
TEST_EXIT_CODE=${PIPESTATUS[0]}
set -e

END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))

# Parse results
PASSED=$(grep -c "^--- PASS:" "$TEST_OUTPUT_FILE" 2>/dev/null || true)
FAILED=$(grep -c "^--- FAIL:" "$TEST_OUTPUT_FILE" 2>/dev/null || true)
SKIPPED=$(grep -c "^--- SKIP:" "$TEST_OUTPUT_FILE" 2>/dev/null || true)

# Ensure numeric values
PASSED=${PASSED:-0}
FAILED=${FAILED:-0}
SKIPPED=${SKIPPED:-0}

# Cleanup temp file
rm -f "$TEST_OUTPUT_FILE"

# Stop test environment if requested
if [ "$STOP_ENV" = true ]; then
    echo ""
    echo -e "${BLUE}Stopping test environment...${NC}"
    if [ -f "$TESTS_DIR/test-env.sh" ]; then
        "$TESTS_DIR/test-env.sh" stop
    else
        docker compose -f "$COMPOSE_FILE" down
    fi
fi

# Print summary
echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Test Summary${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo "Duration: ${DURATION}s"
echo ""
echo "Results:"
echo -e "  ${GREEN}Passed:${NC}  $PASSED"
if [ "$FAILED" -gt 0 ]; then
    echo -e "  ${RED}Failed:${NC}  $FAILED"
else
    echo -e "  Failed:  $FAILED"
fi
if [ "$SKIPPED" -gt 0 ]; then
    echo -e "  ${YELLOW}Skipped:${NC} $SKIPPED"
fi
echo ""

if [ $TEST_EXIT_CODE -eq 0 ]; then
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}ALL TESTS PASSED${NC}"
    echo -e "${GREEN}========================================${NC}"
    exit 0
else
    echo -e "${RED}========================================${NC}"
    echo -e "${RED}SOME TESTS FAILED${NC}"
    echo -e "${RED}========================================${NC}"
    
    if [ "$KEEP_ENV" = true ] && [ "$START_ENV" = true ]; then
        echo -e "${YELLOW}Keeping test environment running for debugging${NC}"
    fi
    
    exit 1
fi