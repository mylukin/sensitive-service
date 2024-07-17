# Makefile

IMAGE_NAME := lukin/sensitive-service
IMAGE_TAG := $(shell date +%Y%m%d)

.PHONY: run stop restart docker-build docker-push docker-pull docker-run docker-stop docker-restart

# 运行 Go 应用程序
run:
	go run main.go

# 停止 Go 应用程序
stop:
	@pkill -f "main"

# 重启 Go 应用程序
restart: stop
	$(MAKE) run

# 构建 Docker 镜像
docker-build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .
	docker tag $(IMAGE_NAME):$(IMAGE_TAG) $(IMAGE_NAME):latest

# 推送 Docker 镜像到仓库
docker-push: docker-build
	docker push $(IMAGE_NAME):$(IMAGE_TAG)
	docker push $(IMAGE_NAME):latest

# 拉取 Docker 镜像
docker-pull:
	docker pull $(IMAGE_NAME):latest

# 运行 Docker 容器
docker-run:
	docker run -d --name sensitive-service -p 8080:8080 -v ./tables_data.json:/app/tables_data.json $(IMAGE_NAME):latest

# 停止并移除 Docker 容器
docker-stop:
	docker stop sensitive-service
	docker rm sensitive-service

# 重启 Docker 容器
docker-restart: docker-stop docker-run
