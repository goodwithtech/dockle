SOURCE := ./cmd/dockle_remote/
BINARY := dockle_remote
DOCKER_REGISTRY ?= gcr.io/eticloud/k8sec
VERSION ?= $(shell git rev-parse HEAD)
IMAGE_NAME ?= $(DOCKER_REGISTRY)/dockle:$(VERSION)
TARGET_OS ?= linux

build:
	GOOS=$(TARGET_OS) CGO_ENABLED=0 go build -o $(BINARY) $(SOURCE)

# builds the current dev docker version
build-docker:
	docker build --build-arg BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ") \
		--build-arg VCS_REF=$(shell git rev-parse --short HEAD) \
		--build-arg IMAGE_VERSION=${VERSION} \
			-t $(IMAGE_NAME) .

docker.push: build-docker
	docker push $(IMAGE_NAME)

test:
	GO111MODULE=on CGO_ENABLED=0 go test -v `go list ./...`

