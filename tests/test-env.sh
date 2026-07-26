#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
COMPOSE_FILE="$SCRIPT_DIR/docker-compose.test.yml"

usage() {
    echo "Usage: $0 {start|stop|restart|status|logs|clean}"
    echo ""
    echo "Commands:"
    echo "  start   - Start the test environment"
    echo "  stop    - Stop the test environment"
    echo "  restart - Restart the test environment"
    echo "  status  - Show status of test containers"
    echo "  logs    - Show logs from test containers"
    echo "  clean   - Remove containers and clean up test data"
    echo ""
    echo "Environment Variables:"
    echo "  CD2_PORT       - CloudDrive2 API port (default: 19798)"
    echo "  CD2_USER       - CloudDrive2 username (default: admin)"
    echo "  CD2_PASS       - CloudDrive2 password (default: admin123)"
    echo "  MINIO_PORT     - MinIO API port (default: 9000)"
    echo "  MINIO_CONSOLE  - MinIO console port (default: 9001)"
    echo "  WEBDAV_PORT    - WebDAV port (default: 8080)"
    echo "  SFTP_PORT      - SFTP port (default: 2222)"
}

start() {
    echo "Starting test environment..."
    mkdir -p "$SCRIPT_DIR/test-storage"/{minio,webdav,sftp}
    mkdir -p "$SCRIPT_DIR/test-data"
    mkdir -p "$SCRIPT_DIR/test-config"
    
    if [ ! -f "$SCRIPT_DIR/test-config/rclone.conf" ]; then
        echo "[webdav]" > "$SCRIPT_DIR/test-config/rclone.conf"
        echo "type = webdav" >> "$SCRIPT_DIR/test-config/rclone.conf"
        echo "url = http://localhost:8080" >> "$SCRIPT_DIR/test-config/rclone.conf"
    fi
    
    docker-compose -f "$COMPOSE_FILE" up -d
    
    echo ""
    echo "Waiting for services to be ready..."
    sleep 5
    
    for i in {1..30}; do
        if curl -sf http://localhost:${CD2_PORT:-19798}/api/GetSystemInfo > /dev/null 2>&1; then
            echo "CloudDrive2 is ready!"
            break
        fi
        echo "Waiting for CloudDrive2... ($i/30)"
        sleep 2
    done
    
    echo ""
    echo "Test environment started!"
    echo "CloudDrive2 API: http://localhost:${CD2_PORT:-19798}"
    echo "MinIO API:       http://localhost:${MINIO_PORT:-9000}"
    echo "MinIO Console:   http://localhost:${MINIO_CONSOLE_PORT:-9001}"
    echo "WebDAV:          http://localhost:${WEBDAV_PORT:-8080}"
    echo "SFTP:            sftp://localhost:${SFTP_PORT:-2222}"
}

stop() {
    echo "Stopping test environment..."
    docker-compose -f "$COMPOSE_FILE" down
    echo "Test environment stopped."
}

restart() {
    stop
    start
}

status() {
    docker-compose -f "$COMPOSE_FILE" ps
}

logs() {
    docker-compose -f "$COMPOSE_FILE" logs -f "$@"
}

clean() {
    echo "Cleaning up test environment..."
    docker-compose -f "$COMPOSE_FILE" down -v
    rm -rf "$SCRIPT_DIR/test-storage"
    rm -rf "$SCRIPT_DIR/test-data"
    rm -rf "$SCRIPT_DIR/test-config"
    echo "Test environment cleaned."
}

case "${1:-}" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        restart
        ;;
    status)
        status
        ;;
    logs)
        shift
        logs "$@"
        ;;
    clean)
        clean
        ;;
    *)
        usage
        exit 1
        ;;
esac