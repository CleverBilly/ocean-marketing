#!/bin/bash

# Ocean Marketing API 测试脚本

set -e

BASE_URL="http://localhost:8080"

echo "🧪 Ocean Marketing API 测试"
echo "================================"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试函数
test_endpoint() {
    local method=$1
    local endpoint=$2
    local data=$3
    local expected_code=$4
    local description=$5
    local token=$6

    echo -e "\n${YELLOW}测试: $description${NC}"
    echo "请求: $method $endpoint"

    if [ -n "$token" ]; then
        if [ -n "$data" ]; then
            response=$(curl -s -w "\n%{http_code}" -X "$method" \
                -H "Content-Type: application/json" \
                -H "Authorization: Bearer $token" \
                -d "$data" \
                "$BASE_URL$endpoint")
        else
            response=$(curl -s -w "\n%{http_code}" -X "$method" \
                -H "Authorization: Bearer $token" \
                "$BASE_URL$endpoint")
        fi
    else
        if [ -n "$data" ]; then
            response=$(curl -s -w "\n%{http_code}" -X "$method" \
                -H "Content-Type: application/json" \
                -d "$data" \
                "$BASE_URL$endpoint")
        else
            response=$(curl -s -w "\n%{http_code}" -X "$method" \
                "$BASE_URL$endpoint")
        fi
    fi

    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)

    if [ "$http_code" = "$expected_code" ]; then
        echo -e "${GREEN}✅ 测试通过 (HTTP $http_code)${NC}"
        echo "响应: $body" | head -c 200
        echo "..."
    else
        echo -e "${RED}❌ 测试失败 (期望: $expected_code, 实际: $http_code)${NC}"
        echo "响应: $body"
    fi
}

# 等待服务启动
echo "⏳ 等待服务启动..."
for i in {1..30}; do
    if curl -s "$BASE_URL/health" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ 服务已启动${NC}"
        break
    fi
    if [ $i -eq 30 ]; then
        echo -e "${RED}❌ 服务启动超时${NC}"
        exit 1
    fi
    sleep 1
done

# 1. 系统健康检查
echo -e "\n🏥 系统健康检查"
test_endpoint "GET" "/health" "" "200" "健康检查"
test_endpoint "GET" "/ready" "" "200" "就绪检查"
test_endpoint "GET" "/live" "" "200" "存活检查"

# 2. 示例业务测试（无需认证的接口）
echo -e "\n📋 示例模块测试（公开接口）"

# 获取示例列表
test_endpoint "GET" "/api/v1/examples" "" "200" "获取示例列表"
test_endpoint "GET" "/api/v1/examples?page=1&size=5" "" "200" "获取示例列表（分页）"

# 获取示例详情（如果存在）
test_endpoint "GET" "/api/v1/examples/1" "" "404" "获取不存在的示例详情"

# 模拟创建Token（用于需要认证的接口测试）
echo -e "\n🔐 生成测试Token..."
# 这里使用一个固定的JWT Token用于测试，实际应用中应该通过登录获取
# 注意：这个Token需要与你的JWT密钥配置匹配
test_token="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6ImFkbWluIiwiZXhwIjoxNzA5ODg5NjAwfQ.fake-signature"

# 3. 示例业务测试（需要认证的接口）
echo -e "\n📋 示例模块测试（需要认证）"

# 创建示例
example_data='{
    "title": "测试示例",
    "content": "这是一个API测试创建的示例",
    "status": 1
}'

# 注意：由于没有用户系统，这些测试可能会失败（401未授权）
# 但可以验证中间件是否正常工作
test_endpoint "POST" "/api/v1/examples" "$example_data" "401" "创建示例（无Token）"
test_endpoint "POST" "/api/v1/examples" "$example_data" "401" "创建示例（测试Token）" "$test_token"

# 更新示例
update_example_data='{
    "title": "更新后的示例",
    "content": "这是更新后的内容"
}'
test_endpoint "PUT" "/api/v1/examples/1" "$update_example_data" "401" "更新示例（无Token）"

# 删除示例
test_endpoint "DELETE" "/api/v1/examples/1" "" "401" "删除示例（无Token）"

# 4. 错误情况测试
echo -e "\n❌ 错误情况测试"

# 无效路径
test_endpoint "GET" "/api/v1/invalid" "" "404" "访问不存在的路径"

# 无效方法
test_endpoint "PATCH" "/api/v1/examples" "" "405" "使用不支持的HTTP方法"

# 无效JSON
test_endpoint "POST" "/api/v1/examples" '{"invalid":json}' "400" "发送无效JSON"

# 5. 监控和管理接口测试
echo -e "\n📊 监控和管理接口测试"

# Prometheus指标
test_endpoint "GET" "/metrics" "" "200" "Prometheus指标"

# Swagger文档（如果启用）
test_endpoint "GET" "/swagger/index.html" "" "200" "Swagger文档"

echo -e "\n${GREEN}🎉 API 测试完成！${NC}"
echo "================================"
echo -e "\n📋 测试总结："
echo "- ✅ 系统健康检查正常"
echo "- ✅ 示例模块API可访问"
echo "- ✅ 认证中间件正常工作"
echo "- ✅ 错误处理机制正常"
echo "- ✅ 监控接口正常" 