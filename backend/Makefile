APP=server
TEST_PKG=./...

build:
	go build -o bin/${APP} cmd/${APP}/main.go

.PHONY: docs
docs:
	swag init -o ./docs/${APP} -d cmd/${APP} --parseDependency

run: build docs
	./bin/${APP}

clean:
	go clean
	rm -rf bin && mkdir bin

test:
	go test -v ${TEST_PKG}

test-app:
	go test -v ./cmd/${APP}/... ./internal/...
