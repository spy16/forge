NAME="forge"
VERSION=$(shell git describe --tags --always --first-parent 2>/dev/null)
COMMIT=$(shell git rev-parse --short HEAD)
BUILD_TIME=$(shell date)
BUILD_DIR=bin

all: tidy test build

tidy:
	@echo "Tidy up..."
	@go mod tidy -v

test:
	@echo "Running tests..."
	@go test -cover ./...

build:
	@mkdir -p ${BUILD_DIR}
	@echo "Running build for '${VERSION}' in '${BUILD_DIR}/'..."
	@CGO_ENABLED=0 go build -ldflags '-X "main.Version=${VERSION}" -X "main.Commit=${COMMIT}" -X "main.BuildTime=${BUILD_TIME}"' -o ${BUILD_DIR}/${EXE} ./cmd/forge

install:
	@echo "Installing..."
	@go install ./cmd/forge