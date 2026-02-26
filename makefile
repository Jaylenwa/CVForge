# ===== 可配置变量 =====
IMAGE_NAME := cvforge-backend:latest
TAR_NAME := cvforge-backend.tar
DOCKERFILE := Dockerfile
FRONTEND_IMAGE_NAME := cvforge-frontend:latest
FRONTEND_TAR_NAME := cvforge-frontend.tar
FRONTEND_DOCKERFILE := frontend/Dockerfile

REMOTE_USER := root
REMOTE_HOST := 182.254.166.74
REMOTE_PORT := 22
REMOTE_DIR := /root/docker
REMOTE_COMPOSE_DIR := /root/CVForge

# ===== 平台控制 =====
PLATFORM := linux/amd64

# ===== 默认目标 =====
.PHONY: all
all: backend frontend deploy clean

.PHONY: backend
backend: build save upload load

# ===== 构建镜像 =====
.PHONY: build
build:
	-docker image rm -f $(IMAGE_NAME)
	docker build --platform=$(PLATFORM) -f $(DOCKERFILE) -t $(IMAGE_NAME) .

.PHONY: build-frontend
build-frontend:
	-docker image rm -f $(FRONTEND_IMAGE_NAME)
	docker build --platform=$(PLATFORM) -f $(FRONTEND_DOCKERFILE) -t $(FRONTEND_IMAGE_NAME) frontend

# ===== 导出镜像 =====
.PHONY: save
save:
	docker save -o $(TAR_NAME) $(IMAGE_NAME)

.PHONY: save-frontend
save-frontend:
	docker save -o $(FRONTEND_TAR_NAME) $(FRONTEND_IMAGE_NAME)

# ===== 上传镜像（远端存在则删除）=====
.PHONY: upload
upload:
	ssh -p $(REMOTE_PORT) $(REMOTE_USER)@$(REMOTE_HOST) "mkdir -p $(REMOTE_DIR) && rm -f $(REMOTE_DIR)/$(TAR_NAME)"
	scp -P $(REMOTE_PORT) $(TAR_NAME) $(REMOTE_USER)@$(REMOTE_HOST):$(REMOTE_DIR)/

.PHONY: upload-frontend
upload-frontend:
	ssh -p $(REMOTE_PORT) $(REMOTE_USER)@$(REMOTE_HOST) "mkdir -p $(REMOTE_DIR) && rm -f $(REMOTE_DIR)/$(FRONTEND_TAR_NAME)"
	scp -P $(REMOTE_PORT) $(FRONTEND_TAR_NAME) $(REMOTE_USER)@$(REMOTE_HOST):$(REMOTE_DIR)/

# ===== 远端加载镜像 =====
.PHONY: load
load:
	ssh -p $(REMOTE_PORT) $(REMOTE_USER)@$(REMOTE_HOST) "docker load -i $(REMOTE_DIR)/$(TAR_NAME)"

.PHONY: load-frontend
load-frontend:
	ssh -p $(REMOTE_PORT) $(REMOTE_USER)@$(REMOTE_HOST) "docker load -i $(REMOTE_DIR)/$(FRONTEND_TAR_NAME)"

.PHONY: frontend
frontend: build-frontend save-frontend upload-frontend load-frontend

# ===== 清理本地 tar =====
.PHONY: clean
clean:
	rm -f $(TAR_NAME)
	rm -f $(FRONTEND_TAR_NAME)

# ===== 远端部署（加载镜像并重启 compose）=====
.PHONY: deploy-all
deploy-all: build build-frontend save save-frontend upload upload-frontend
	ssh -p $(REMOTE_PORT) $(REMOTE_USER)@$(REMOTE_HOST) "cd $(REMOTE_DIR) && docker load -i $(TAR_NAME) && docker load -i $(FRONTEND_TAR_NAME)"
	ssh -p $(REMOTE_PORT) $(REMOTE_USER)@$(REMOTE_HOST) "cd $(REMOTE_COMPOSE_DIR) && docker compose down && docker compose up -d"

# ===== 远端部署（加载镜像并重启 compose）=====
.PHONY: deploy
deploy:
	ssh -p $(REMOTE_PORT) $(REMOTE_USER)@$(REMOTE_HOST) "cd $(REMOTE_COMPOSE_DIR) && docker compose down && docker compose up -d"

.PHONY: commit
commit:
	git add .
	git commit -m "$$(date '+%Y-%m-%d %H:%M:%S')"
