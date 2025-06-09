 # Ocean Marketing 阿里云部署指南

本文档详细介绍如何将 Ocean Marketing 项目部署到阿里云环境中。

## 🏗️ 阿里云基础设施准备

### 1. 云资源准备清单

#### 必需资源
- **ECS 实例**: 用于部署应用服务器
- **RDS MySQL**: 云数据库 MySQL 版
- **Redis**: 云数据库 Redis 版
- **SLB**: 负载均衡器（可选，用于高可用）
- **VPC**: 专有网络
- **安全组**: 网络访问控制

#### 可选资源
- **OSS**: 对象存储（用于文件存储）
- **CDN**: 内容分发网络
- **SLS**: 日志服务
- **云监控**: 应用性能监控
- **消息队列 RabbitMQ 版**: 企业级消息队列

### 2. 网络架构设计

```
Internet
    |
   SLB (负载均衡)
    |
   ECS (应用服务器)
    |
   VPC (专有网络)
   ├── RDS MySQL (数据库)
   ├── Redis (缓存)
   └── RabbitMQ (消息队列)
```

## 🛠️ 部署步骤

### 步骤 1: 创建云资源

#### 1.1 创建 VPC 和交换机
```bash
# 在阿里云控制台创建 VPC
# - 地域: 根据业务需求选择（推荐华东1-杭州）
# - 网段: 192.168.0.0/16
# - 交换机: 192.168.1.0/24
```

#### 1.2 创建 RDS MySQL 实例
```bash
# 在 RDS 控制台创建 MySQL 实例
# - 版本: MySQL 8.0
# - 规格: 根据业务需求选择
# - 存储: SSD 云盘
# - 网络类型: VPC
# - 可用区: 与 ECS 同可用区
```

#### 1.3 创建 Redis 实例
```bash
# 在 Redis 控制台创建实例
# - 版本: Redis 6.0
# - 架构: 标准版或集群版
# - 网络类型: VPC
```

#### 1.4 创建 ECS 实例
```bash
# 在 ECS 控制台创建实例
# - 操作系统: CentOS 7.9 或 Ubuntu 20.04
# - 实例规格: 根据业务需求选择
# - 网络: VPC
# - 安全组: 开放 22(SSH), 8080(应用) 端口
```

### 步骤 2: 配置安全组

#### 2.1 入方向规则
```bash
# SSH 访问
22/22    TCP    授权对象: 0.0.0.0/0    描述: SSH

# 应用端口
8080/8080    TCP    授权对象: 0.0.0.0/0    描述: 应用服务

# 数据库端口（仅 VPC 内访问）
3306/3306    TCP    授权对象: 192.168.0.0/16    描述: MySQL
6379/6379    TCP    授权对象: 192.168.0.0/16    描述: Redis
```

### 步骤 3: 部署应用

#### 3.1 连接到 ECS 实例
```bash
ssh root@<ECS公网IP>
```

#### 3.2 安装 Docker 和 Docker Compose
```bash
# 安装 Docker
curl -fsSL https://get.docker.com | bash -s docker --mirror Aliyun
systemctl start docker
systemctl enable docker

# 安装 Docker Compose
curl -L "https://github.com/docker/compose/releases/download/v2.20.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose
```

#### 3.3 克隆项目代码
```bash
# 安装 Git
yum install -y git  # CentOS
# 或
apt-get update && apt-get install -y git  # Ubuntu

# 克隆代码
git clone <your-repository-url> /opt/ocean-marketing
cd /opt/ocean-marketing
```

#### 3.4 配置生产环境
```bash
# 复制生产配置
cp configs/app.production.yaml configs/app.yaml

# 编辑配置文件，填入实际的阿里云资源信息
vim configs/app.yaml
```

#### 3.5 配置重要信息
```yaml
# 修改 configs/app.yaml 中的以下配置:

database:
  host: rm-xxxxxxxxxxxxxxx.mysql.rds.aliyuncs.com  # RDS 内网地址
  username: your_db_user
  password: "your_strong_password"

redis:
  host: r-xxxxxxxxxxxxxxx.redis.rds.aliyuncs.com  # Redis 内网地址
  password: "your_redis_password"

jwt:
  secret: "your-super-secret-jwt-key-32-characters-long"

feishu:
  webhook_url: "https://open.feishu.cn/open-apis/bot/v2/hook/your-webhook-url"
```

#### 3.6 构建和启动服务
```bash
# 构建 Docker 镜像
docker build -t ocean-marketing:latest .

# 启动服务（仅应用，不启动本地数据库）
docker run -d \
  --name ocean-marketing \
  --restart unless-stopped \
  -p 8080:8080 \
  -v /opt/ocean-marketing/configs:/app/configs \
  -v /opt/ocean-marketing/logs:/app/logs \
  -v /opt/ocean-marketing/data:/app/data \
  -e TZ=Asia/Shanghai \
  ocean-marketing:latest
```

### 步骤 4: 数据库初始化

#### 4.1 创建数据库
```bash
# 连接到 RDS MySQL
mysql -h rm-xxxxxxxxxxxxxxx.mysql.rds.aliyuncs.com -u root -p

# 创建数据库和用户
CREATE DATABASE ocean_marketing CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'ocean_marketing'@'%' IDENTIFIED BY 'your_strong_password';
GRANT ALL PRIVILEGES ON ocean_marketing.* TO 'ocean_marketing'@'%';
FLUSH PRIVILEGES;
```

#### 4.2 应用会自动执行数据库迁移
```bash
# 查看应用启动日志
docker logs -f ocean-marketing
```

### 步骤 5: 配置反向代理（可选）

#### 5.1 安装 Nginx
```bash
yum install -y nginx  # CentOS
# 或
apt-get install -y nginx  # Ubuntu
```

#### 5.2 配置 Nginx
```nginx
# /etc/nginx/conf.d/ocean-marketing.conf
server {
    listen 80;
    server_name your-domain.com;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
    
    # 健康检查
    location /health {
        access_log off;
        proxy_pass http://localhost:8080/health;
    }
}
```

## 🔧 运维配置

### 1. 自动启动脚本
```bash
# 创建 systemd 服务文件
cat > /etc/systemd/system/ocean-marketing.service << EOF
[Unit]
Description=Ocean Marketing Application
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
ExecStart=/usr/bin/docker start ocean-marketing
ExecStop=/usr/bin/docker stop ocean-marketing

[Install]
WantedBy=multi-user.target
EOF

systemctl enable ocean-marketing
```

### 2. 日志管理
```bash
# 配置日志轮转
cat > /etc/logrotate.d/ocean-marketing << EOF
/opt/ocean-marketing/logs/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    copytruncate
}
EOF
```

### 3. 监控脚本
```bash
# 创建健康检查脚本
cat > /opt/ocean-marketing/health-check.sh << 'EOF'
#!/bin/bash
response=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health)
if [ $response != "200" ]; then
    echo "Application is unhealthy, restarting..."
    docker restart ocean-marketing
    # 发送告警通知
    curl -X POST "$FEISHU_WEBHOOK" \
        -H "Content-Type: application/json" \
        -d '{"msg_type":"text","content":{"text":"Ocean Marketing 应用异常重启"}}'
fi
EOF

chmod +x /opt/ocean-marketing/health-check.sh

# 添加到 crontab
echo "*/5 * * * * /opt/ocean-marketing/health-check.sh" | crontab -
```

## 🔐 安全配置

### 1. 防火墙配置
```bash
# CentOS 7
firewall-cmd --permanent --add-port=8080/tcp
firewall-cmd --permanent --add-port=80/tcp
firewall-cmd --permanent --add-port=443/tcp
firewall-cmd --reload

# Ubuntu
ufw allow 8080/tcp
ufw allow 80/tcp
ufw allow 443/tcp
ufw enable
```

### 2. SSL 证书配置
```bash
# 使用 Let's Encrypt 免费证书
snap install --classic certbot
certbot --nginx -d your-domain.com
```

## 📊 监控和告警

### 1. 云监控配置
- 在阿里云云监控控制台配置 ECS、RDS、Redis 监控
- 设置告警规则：CPU > 80%、内存 > 80%、磁盘 > 80%

### 2. 应用监控
- 访问 `http://your-domain.com/metrics` 查看 Prometheus 指标
- 配置 Grafana 仪表板

### 3. 日志监控
- 配置 SLS 日志服务收集应用日志
- 设置关键错误日志告警

## 🚀 性能优化

### 1. 数据库优化
```sql
-- 数据库参数优化
SET GLOBAL innodb_buffer_pool_size = 1073741824;  -- 1GB
SET GLOBAL max_connections = 1000;
SET GLOBAL query_cache_size = 268435456;  -- 256MB
```

### 2. Redis 优化
```bash
# Redis 配置优化
maxmemory 512mb
maxmemory-policy allkeys-lru
save 900 1
save 300 10
```

### 3. 应用优化
```yaml
# 调整应用配置
database:
  max_idle_conns: 50
  max_open_conns: 200
  conn_max_lifetime: 3600

redis:
  pool_size: 100
  min_idle_conns: 20
```

## 🔄 CI/CD 集成

### 1. 使用阿里云容器镜像服务
```bash
# 在阿里云容器镜像服务创建镜像仓库
# 配置自动构建触发器
```

### 2. 部署脚本
```bash
#!/bin/bash
# deploy.sh

# 拉取最新镜像
docker pull registry.cn-hangzhou.aliyuncs.com/your-namespace/ocean-marketing:latest

# 停止当前容器
docker stop ocean-marketing
docker rm ocean-marketing

# 启动新容器
docker run -d \
  --name ocean-marketing \
  --restart unless-stopped \
  -p 8080:8080 \
  -v /opt/ocean-marketing/configs:/app/configs \
  -v /opt/ocean-marketing/logs:/app/logs \
  -v /opt/ocean-marketing/data:/app/data \
  -e TZ=Asia/Shanghai \
  registry.cn-hangzhou.aliyuncs.com/your-namespace/ocean-marketing:latest

echo "Deployment completed!"
```

## 📞 故障排除

### 常见问题

1. **数据库连接失败**
   - 检查 RDS 白名单配置
   - 验证网络连通性
   - 确认用户名密码正确

2. **Redis 连接失败**
   - 检查 Redis 白名单配置
   - 验证密码配置

3. **应用启动失败**
   - 查看 Docker 日志：`docker logs ocean-marketing`
   - 检查配置文件语法
   - 验证端口占用情况

4. **性能问题**
   - 检查云监控指标
   - 分析应用日志
   - 优化数据库查询

## 📝 总结

通过以上步骤，您可以成功将 Ocean Marketing 应用部署到阿里云环境中。部署完成后，请务必：

1. 定期备份数据库
2. 监控应用性能
3. 及时更新安全补丁
4. 优化资源配置

如有问题，请参考阿里云官方文档或联系技术支持。