#!/bin/bash

# Ocean Marketing 项目完整性检查脚本

set -e

echo "🔍 Ocean Marketing 项目完整性检查"
echo "=================================="

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 检查函数
check_file() {
    local file=$1
    local description=$2
    
    if [ -f "$file" ]; then
        echo -e "${GREEN}✅ $description${NC}: $file"
    else
        echo -e "${RED}❌ $description${NC}: $file (不存在)"
    fi
}

check_dir() {
    local dir=$1
    local description=$2
    
    if [ -d "$dir" ]; then
        echo -e "${GREEN}✅ $description${NC}: $dir"
    else
        echo -e "${RED}❌ $description${NC}: $dir (不存在)"
    fi
}

echo -e "\n${BLUE}📂 目录结构检查${NC}"
echo "----------------------------"

# 核心目录
check_dir "cmd/server" "主程序目录"
check_dir "internal" "内部包目录"
check_dir "pkg" "公共包目录"
check_dir "configs" "配置目录"
check_dir "scripts" "脚本目录"
check_dir "web/static" "静态文件目录"
check_dir "docs" "文档目录"

echo -e "\n${BLUE}🗂️ 核心文件检查${NC}"
echo "----------------------------"

# 主要文件
check_file "go.mod" "Go模块文件"
check_file "go.sum" "Go依赖校验文件"
check_file "cmd/server/main.go" "主程序文件"
check_file "Makefile" "构建文件"
check_file "Dockerfile" "Docker文件"
check_file "docker-compose.yml" "Docker Compose文件"
check_file ".gitignore" "Git忽略文件"
check_file "README.md" "项目说明文档"

echo -e "\n${BLUE}⚙️ 配置文件检查${NC}"
echo "----------------------------"

check_file "configs/app.yaml" "应用配置文件"
check_file "configs/prometheus.yml" "Prometheus配置文件"

echo -e "\n${BLUE}📋 内部包结构检查${NC}"
echo "----------------------------"

# 配置包
check_file "internal/config/config.go" "配置管理"

# 模型包
check_file "internal/model/example.go" "示例模型"

# 服务包
check_file "internal/service/example.go" "示例服务"

# 控制器包
check_file "internal/handler/example.go" "示例控制器"
check_file "internal/handler/health.go" "健康检查控制器"
check_file "internal/handler/v1/example.go" "V1示例控制器"

# 中间件包
check_file "internal/middleware/validation.go" "验证中间件"
check_file "internal/middleware/cors.go" "CORS中间件"
check_file "internal/middleware/logger.go" "日志中间件"
check_file "internal/middleware/recovery.go" "恢复中间件"
check_file "internal/middleware/ratelimit.go" "限流中间件"
check_file "internal/middleware/tracer.go" "追踪中间件"
check_file "internal/middleware/prometheus.go" "监控中间件"
check_file "internal/middleware/auth.go" "认证中间件"

# 路由包
check_file "internal/router/router.go" "路由注册"

# 基础设施包
check_file "internal/pkg/logger/logger.go" "日志组件"
check_file "internal/pkg/database/database.go" "数据库组件"
check_file "internal/pkg/redis/redis.go" "Redis组件"
check_file "internal/pkg/tracer/tracer.go" "链路追踪组件"
check_file "internal/pkg/migration/migration.go" "数据库迁移"

echo -e "\n${BLUE}📦 公共包检查${NC}"
echo "----------------------------"

check_file "pkg/errno/errno.go" "错误码定义"
check_file "pkg/response/response.go" "响应处理"
check_file "pkg/jwt/jwt.go" "JWT处理"
check_file "pkg/email/email.go" "邮件发送"
check_file "pkg/cast/cast.go" "类型转换"
check_file "pkg/mq/mq.go" "消息队列"

echo -e "\n${BLUE}📄 文档和脚本检查${NC}"
echo "----------------------------"

check_file "docs/project-overview.md" "项目概述"
check_file "docs/validation-middleware-guide.md" "验证中间件指南"
check_file "docs/aliyun-deployment.md" "阿里云部署指南"
check_file "scripts/start.sh" "启动脚本"
check_file "scripts/test-api.sh" "API测试脚本"
check_file "scripts/check-project.sh" "项目检查脚本"
check_file "web/static/index.html" "主页文件"

echo -e "\n${BLUE}🔧 编译检查${NC}"
echo "----------------------------"

echo -n "检查 Go 编译: "
if go build cmd/server/main.go >/dev/null 2>&1; then
    echo -e "${GREEN}✅ 编译成功${NC}"
    rm -f main 2>/dev/null || true
else
    echo -e "${RED}❌ 编译失败${NC}"
fi

echo -e "\n${BLUE}📊 项目统计${NC}"
echo "----------------------------"

# 统计代码行数
echo "Go 源文件数量: $(find . -name "*.go" | wc -l)"
echo "Go 代码行数: $(find . -name "*.go" -exec wc -l {} + | tail -1 | awk '{print $1}')"
echo "配置文件数量: $(find configs -name "*.yaml" -o -name "*.yml" | wc -l)"
echo "脚本文件数量: $(find scripts -name "*.sh" | wc -l)"
echo "文档文件数量: $(find . -name "*.md" | wc -l)"

echo -e "\n${BLUE}🎯 关键特性检查${NC}"
echo "----------------------------"

# 检查关键特性是否已实现
features=(
    "验证中间件:internal/middleware/validation.go"
    "JWT认证:pkg/jwt/jwt.go"
    "限流保护:internal/middleware/ratelimit.go"
    "日志系统:internal/pkg/logger/logger.go"
    "数据库ORM:internal/pkg/database/database.go"
    "Redis缓存:internal/pkg/redis/redis.go"
    "邮件发送:pkg/email/email.go"
    "消息队列:pkg/mq/mq.go"
    "链路追踪:internal/pkg/tracer/tracer.go"
    "监控指标:internal/middleware/prometheus.go"
    "错误处理:pkg/errno/errno.go"
    "类型转换:pkg/cast/cast.go"
    "健康检查:internal/handler/health.go"
    "示例模块:internal/model/example.go"
)

for feature in "${features[@]}"; do
    IFS=':' read -r name file <<< "$feature"
    if [ -f "$file" ]; then
        echo -e "${GREEN}✅ $name${NC}"
    else
        echo -e "${RED}❌ $name${NC}"
    fi
done

echo -e "\n${GREEN}🎉 项目完整性检查完成！${NC}"
echo "==================================" 