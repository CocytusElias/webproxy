BUILD_PATH:=build
DOCKER_VERSION?=latest
GO_BUILD_IMAGE?=golang:1.19.0

.PHONY: build docker
all: build

$(MODULES): build:
deps:
	go mod tidy && go mod download
gen:
	CGO_ENABLED=0 go run ./client-gen.go

build:
	$(MAKE) build-service && $(MAKE) build-client
build-service:
	$(MAKE) deps && CGO_ENABLED=0 go build -o $(BUILD_PATH)/service ./service.go
build-client:
	$(MAKE) deps && $(MAKE) gen && CGO_ENABLED=0 go build -o $(BUILD_PATH)/client ./client.go

arm64:
	$(MAKE) arm64-service && $(MAKE) arm64-client
arm64-service:
	$(MAKE) deps && CGO_ENABLED=0 GOARCH=arm64 go build -ldflags "-w -s" -o $(BUILD_PATH)/service-arm64 ./service.go
arm64-client:
	$(MAKE) deps && $(MAKE) gen && CGO_ENABLED=0 GOARCH=arm64 go build -ldflags "-w -s" -o $(BUILD_PATH)/client-arm64 ./client.go

amd64:
	$(MAKE) amd64-service && $(MAKE) amd64-client
amd64-service:
	$(MAKE) deps && CGO_ENABLED=0 GOARCH=amd64 go build -ldflags "-w -s" -o $(BUILD_PATH)/service-amd64 ./service.go
amd64-client:
	$(MAKE) deps && $(MAKE) gen  && CGO_ENABLED=0 GOARCH=amd64 go build -ldflags "-w -s" -o $(BUILD_PATH)/client-amd64 ./client.go

run-service:
	$(MAKE) build-service && $(BUILD_PATH)/service
run-client:
	$(MAKE) build-client && $(BUILD_PATH)/client

clean:
	go clean -modcache && rm -rf ./build

$(MODULES): docker:
docker:
	$(MAKE) docker-service && $(MAKE) docker-client
docker-service:
	docker build --build-arg GO_BUILD_IMAGE=$(GO_BUILD_IMAGE) -t cocytuselias2023/webproxy:service-$(DOCKER_VERSION) --no-cache --target service .
docker-client:
	docker build --build-arg GO_BUILD_IMAGE=$(GO_BUILD_IMAGE) -t cocytuselias2023/webproxy:client-$(DOCKER_VERSION) --no-cache --target client .

