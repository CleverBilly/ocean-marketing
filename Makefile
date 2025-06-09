.PHONY: help build run clean test lint fmt deps docker swagger

# 默认目标
help: ## 显示帮助信息
	@echo "可用的命令："
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# 构建相关
build: ## 构建应用程序
	@echo "构建应用程序..."
	go build -o bin/server cmd/server/main.go

run: ## 运行应用程序
	@echo "启动应用程序..."
	go run cmd/server/main.go

clean: ## 清理构建文件
	@echo "清理构建文件..."
	rm -rf bin/
	rm -rf logs/
	go clean

# 依赖管理
deps: ## 安装依赖
	@echo "安装依赖..."
	go mod download
	go mod tidy

deps-update: ## 更新依赖
	@echo "更新依赖..."
	go get -u ./...
	go mod tidy

# 代码质量
test: ## 运行测试
	@echo "运行测试..."
	go test -v ./...

test-coverage: ## 运行测试并生成覆盖率报告
	@echo "运行测试覆盖率..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint: ## 运行代码检查
	@echo "运行代码检查..."
	golangci-lint run

fmt: ## 格式化代码
	@echo "格式化代码..."
	go fmt ./...
	goimports -w .

# 文档生成
swagger: ## 生成Swagger文档
	@echo "生成Swagger文档..."
	swag init -g cmd/server/main.go -o docs

swagger-install: ## 安装Swagger工具
	@echo "安装Swagger工具..."
	go install github.com/swaggo/swag/cmd/swag@latest

# Docker相关
docker-build: ## 构建Docker镜像
	@echo "构建Docker镜像..."
	docker build -t ocean-marketing:latest .

docker-run: ## 运行Docker容器
	@echo "运行Docker容器..."
	docker run -p 8080:8080 ocean-marketing:latest

docker-compose-up: ## 启动docker-compose
	@echo "启动docker-compose..."
	docker-compose up -d

docker-compose-down: ## 停止docker-compose
	@echo "停止docker-compose..."
	docker-compose down

# 部署相关
deploy-dev: ## 部署到开发环境
	@echo "部署到开发环境..."
	# 添加部署脚本

deploy-prod: ## 部署到生产环境
	@echo "部署到生产环境..."
	# 添加生产部署脚本

# 数据库相关
db-migrate: ## 运行数据库迁移
	@echo "运行数据库迁移..."
	# 添加数据库迁移命令

db-seed: ## 运行数据库种子
	@echo "运行数据库种子..."
	# 添加数据库种子命令

# 工具安装
install-tools: ## 安装开发工具
	@echo "安装开发工具..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/swaggo/swag/cmd/swag@latest

# 环境准备
dev-setup: deps install-tools swagger ## 开发环境设置
	@echo "开发环境设置完成!"

# 完整测试流程
ci: fmt lint test ## CI流程：格式化、检查、测试
	@echo "CI流程完成!"

# 版本发布
version: ## 显示版本信息
	@echo "Ocean Marketing v1.0.0"
	@echo "Go version: $(shell go version)"
	@echo "Git commit: $(shell git rev-parse --short HEAD)"

# 监控相关
metrics: ## 查看性能指标
	@echo "访问 http://localhost:8080/metrics 查看Prometheus指标"

pprof: ## 打开pprof性能分析
	@echo "访问 http://localhost:8080/debug/pprof/ 进行性能分析"

health: ## 健康检查
	@curl -s http://localhost:8080/health | json_pp || echo "服务未启动" 