# Coze Agent Platform

基于 Gin 框架的 Coze Agent 中台微服务系统，提供完整的对话管理和消息存储功能。

## 核心特性

- **RESTful API**: 基于 Gin 框架，采用标准的 REST 架构
- **分层架构**: Controller -> Service -> Model 清晰分层
- **JWT 身份验证**: 完整的用户认证和授权系统
- **Swagger 文档**: 自动生成的 API 文档
- **消息存储**: 支持会话消息持久化存储
- **历史消息**: 每次对话自动加载最近20条历史消息
- **流式响应**: 支持 Server-Sent Events (SSE) 流式消息推送

## 技术栈

### 后端框架
- **Gin**: HTTP Web 框架
- **GORM**: ORM 库，支持 MySQL

### 数据存储
- **MySQL 8.0**: 主数据库，存储用户、对话、消息等数据
- **Redis**: 缓存和会话存储
- **Milvus**: 向量数据库（预留）

### 消息队列
- **Kafka**: 消息队列（预留）

### 搜索服务
- **ElasticSearch**: 全文搜索（预留）

### 第三方集成
- **Coze SDK**: Coze 官方 SDK 集成，支持双向通信

## 核心功能

### 消息存储功能
- 自动将用户消息和AI回复存储到数据库
- 支持消息的增删查改操作
- 记录消息的元数据（token数量、模型ID等）

### 历史消息功能
- 每次发送消息时自动加载最近20条历史消息
- 按时间顺序排列，提供完整的对话上下文
- 支持分页查询历史消息

### 流式对话功能
- 支持 Server-Sent Events (SSE) 协议
- 实时推送AI回复内容
- 自动保存流式对话的完整内容

## API 端点

### 对话管理
- `GET /api/conversations` - 获取对话列表
- `POST /api/conversations` - 创建新对话
- `GET /api/conversations/{id}` - 获取对话详情
- `DELETE /api/conversations/{id}` - 删除对话

### 消息管理
- `GET /api/conversations/{id}/messages` - 获取消息列表
- `POST /api/conversations/{id}/messages` - 发送消息
- `POST /api/conversations/{id}/messages/stream` - 流式发送消息

### 用户认证
- `POST /api/auth/login` - 用户登录
- `POST /api/auth/register` - 用户注册
- `GET /api/users/profile` - 获取用户资料

## 数据库设计

### 用户表 (user)
- 存储用户基本信息
- 支持用户角色管理
- 软删除支持

### 对话表 (conversation)
- 存储对话会话信息
- 关联 Coze 对话ID
- 支持对话标题自定义

### 消息表 (message)
- 存储所有消息内容
- 记录消息类型（用户/AI）
- 关联对话ID和用户ID
- 记录token消耗情况

## 配置说明

系统配置通过环境变量注入，支持以下配置项：

```yaml
app:
  name: "Coze Agent Platform"
  port: "8080"
  mode: "development"

database:
  host: "localhost"
  port: 3306
  username: "root"
  password: "password"
  dbname: "chatbot"

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0

jwt:
  secret: "your-jwt-secret"
  expire: 7200

coze:
  api_url: "https://api.coze.cn"
  client_id: "your-client-id"
  private_key: "your-private-key"
  public_key_id: "your-public-key-id"
  bot_id: "your-bot-id"
```

## 快速开始

1. **克隆项目**
```bash
git clone <repository-url>
cd coze-agent-platform
```

2. **安装依赖**
```bash
go mod download
```

3. **配置环境**
```bash
cp config/config.yaml.example config/config.yaml
# 修改配置文件中的相关参数
```

4. **初始化数据库**
```bash
# 导入数据库架构
mysql -u root -p < sql/schema.sql
```

5. **启动服务**
```bash
go run cmd/main.go
```

6. **访问 API 文档**
```
http://localhost:8080/swagger/index.html
```

## 健康检查

系统提供健康检查端点：
- `GET /health` - 服务健康状态检查

## 日志记录

- 结构化日志记录（JSON 格式）
- 支持不同日志级别配置
- 请求追踪和错误记录

## 错误处理

- 统一的错误处理中间件
- 标准化的错误响应格式
- 详细的错误信息记录

## 开发说明

### 添加新的API端点
1. 在 `controllers` 目录下添加控制器方法
2. 在 `services` 目录下添加业务逻辑
3. 在 `models` 目录下定义数据模型
4. 在 `routers` 目录下注册路由

### 数据库迁移
使用 GORM 的自动迁移功能：
```go
db.AutoMigrate(&models.User{}, &models.Conversation{}, &models.Message{})
```

## 注意事项

1. 确保 Coze SDK 的配置正确
2. 数据库连接配置需要正确
3. Redis 服务需要正常运行
4. JWT 密钥需要保密

## 许可证

[MIT License](LICENSE) 