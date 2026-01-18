IMAGE_NAME := ghcr.io/orkward/wormface
GIT_TAG := $(shell git describe --tags --always)

.PHONY: all docker push config

all: cli server

cli:
	go build -o wormface-cli ./cmd/wormface-cli/main.go

server:
	go build -o wormface-server ./cmd/wormface-server/main.go

docker:
	docker build \
	  --build-arg GIT_TAG=$(GIT_TAG) \
	  --platform linux/amd64,linux/arm64 \
	  -t $(IMAGE_NAME):$(GIT_TAG) \
	  -t $(IMAGE_NAME):latest .

push: docker
	docker push $(IMAGE_NAME):$(GIT_TAG)
	docker push $(IMAGE_NAME):latest

generate-swagger:
	swag init -g cmd/wormface-server/main.go -o internal/server/docs

config:
	op inject -f -i example.env -o .env
