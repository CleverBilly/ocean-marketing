# Scripts 使用说明

本目录包含了Ocean Marketing项目的常用脚本工具。

## 脚本列表

### 🚀 `start.sh` - 启动脚本
快速启动项目服务的脚本。

**功能：**
- 检查Go环境和配置文件
- 创建必要目录（logs、bin）
- 下载依赖包
- 生成Swagger文档（如果swag已安装）
- 编译并启动服务

**使用方法：**
```bash
./scripts/start.sh
```

**输出示例：**
```
🚀 启动 Ocean Marketing 服务...
📦 下载依赖...
📚 生成 Swagger 文档...
🔨 编译项目...
✅ 启动服务...

🎉 服务启动成功！
📍 访问地址：
   - 应用: http://localhost:8080
   - 健康检查: http://localhost:8080/health
   - API文档: http://localhost:8080/swagger/index.html
   - 监控指标: http://localhost:8080/metrics
```

### 🔍 `check-project.sh` - 项目完整性检查
检查项目文件结构、依赖关系和编译状态的脚本。

**功能：**
- 检查目录结构完整性
- 验证核心文件是否存在
- 检查Go编译状态
- 统计项目代码信息
- 验证关键特性模块

**使用方法：**
```bash
./scripts/check-project.sh
```

**检查内容：**
- 📂 目录结构检查
- 🗂️ 核心文件检查
- ⚙️ 配置文件检查
- 📋 内部包结构检查
- 📦 公共包检查
- 📄 文档和脚本检查
- 🔧 编译检查
- 📊 项目统计
- 🎯 关键特性检查

### 🧪 `test-api.sh` - API接口测试
自动化测试项目API接口的脚本。

**功能：**
- 系统健康检查测试
- Example模块API测试
- 错误情况处理测试
- 监控和管理接口测试

**使用方法：**
```bash
# 确保服务已启动
./scripts/start.sh &

# 运行API测试
./scripts/test-api.sh
```

**测试内容：**
- 🏥 系统健康检查 (`/health`, `/ready`, `/live`)
- 📋 Example模块测试（公开和需认证接口）
- ❌ 错误情况测试（404、401、400等）
- 📊 监控接口测试（`/metrics`, `/swagger`）

## 使用建议

### 开发流程
1. **启动开发环境**
   ```bash
   ./scripts/start.sh
   ```

2. **验证项目完整性**
   ```bash
   ./scripts/check-project.sh
   ```

3. **测试API接口**
   ```bash
   ./scripts/test-api.sh
   ```

### 部署前检查
```bash
# 完整性检查
./scripts/check-project.sh

# API功能测试
./scripts/test-api.sh

# 如果都通过，则可以进行部署
```

### 自动化集成
这些脚本可以集成到CI/CD流水线中：

```yaml
# .github/workflows/test.yml 示例
- name: Check Project
  run: ./scripts/check-project.sh

- name: Test APIs
  run: |
    ./scripts/start.sh &
    sleep 10
    ./scripts/test-api.sh
```

## 注意事项

1. **权限要求**：确保脚本有执行权限
   ```bash
   chmod +x scripts/*.sh
   ```

2. **环境要求**：
   - Go 1.21+
   - 已配置`configs/app.yaml`
   - 数据库和Redis连接正常

3. **网络要求**：
   - API测试需要服务运行在localhost:8080
   - 需要网络访问来下载Go依赖

4. **错误处理**：
   - 所有脚本都使用`set -e`，遇到错误会立即退出
   - 检查输出日志以诊断问题

## 定制化

如需定制脚本行为，可以修改以下变量：

```bash
# test-api.sh 中的基础URL
BASE_URL="http://localhost:8080"

# start.sh 中的编译输出目录
go build -o bin/ocean-marketing cmd/server/main.go
``` 