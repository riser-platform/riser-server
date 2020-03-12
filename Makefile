IMG ?= riserplatform/riser-server
TAG ?= latest

# Run tests.
test: fmt lint tidy test-cmd
	$(TEST_COMMAND)
	# Nested go modules are not tested for some reason, so test them separately
	cd api/v1/model && $(TEST_COMMAND)
	cd pkg/sdk && $(TEST_COMMAND)

test-cmd:
ifeq (, $(shell which gotestsum))
TEST_COMMAND=go test ./...
else
TEST_COMMAND=gotestsum
endif

tidy:
	go mod tidy
	cd api/v1/model && go mod tidy
	cd pkg/sdk && go mod tidy

# Runs the server
run:
	go run ./main.go || true

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
lint:
	golangci-lint run
	cd api/v1/model && golangci-lint run
	cd pkg/sdk && golangci-lint run

# compile and run unit tests on change. Always "make test" before comitting.
# requires filewatcher and gotestsum
watch:
	filewatcher -d api,pkg gotestsum

docker-build:
	docker build . -t riser-server
	docker tag riser-server ${IMG}:${TAG}

docker-push:
	docker push ${IMG}:${TAG}

docker-run: docker-build
	docker run -it --rm -p 8000:8000 -e "TEST_DIR=/riser-server" riser-server

# Updates snapshot tests.
update-snapshot:
	@UPDATESNAPSHOT=true go test ./...
	@echo "Snapshot updated. Check the diff when you commit to ensure that the updates are what you expect!"

