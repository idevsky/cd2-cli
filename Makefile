.PHONY: build test test-cover test-all test-integration test-integration-junit docker-up docker-down docker-logs integration-env-up integration-env-down integration-env-clean clean

GOPROXY := https://goproxy.cn,direct
GO := GOPROXY=$(GOPROXY) go

build:
	$(GO) build -o cd2-cli ./cmd/cd2-cli

test:
	$(GO) test ./internal/client -v

test-cover:
	$(GO) test ./internal/client -v -coverprofile=coverage.out
	$(GO) tool cover -func=coverage.out | tail -1

test-all:
	$(GO) test ./... -v -coverprofile=coverage.out
	$(GO) tool cover -func=coverage.out | tail -1

test-integration:
	cd tests && ./run-integration-tests.sh

test-integration-junit:
	cd tests && ./run-integration-tests.sh --junit

docker-up:
	cd tests && docker compose up -d

docker-down:
	cd tests && docker compose down

docker-logs:
	docker logs clouddrive2-test

integration-env-up:
	cd tests && ./test-env.sh start

integration-env-down:
	cd tests && ./test-env.sh stop

integration-env-clean:
	cd tests && ./test-env.sh clean

integration-env-status:
	cd tests && ./test-env.sh status

integration-env-logs:
	cd tests && ./test-env.sh logs

clean:
	rm -f cd2-cli coverage.out
	rm -rf test-results
	cd tests && docker compose down -v