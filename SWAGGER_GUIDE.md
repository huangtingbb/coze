# Coze Agent Platform - Swagger API 文档指南

## 🚀 快速开始

### 1. 启动服务器
```bash
# 进入项目目录
cd coze-agent-platform

# 启动服务器
go run cmd/main.go
```

服务器启动后，你将看到类似以下的输出：
```
Server starting on port :8080
```

### 2. 访问 Swagger 文档

启动服务器后，可以通过以下 URL 访问 Swagger 文档：

- **Swagger UI**: http://localhost:8080/swagger/index.html
- **健康检查**: http://localhost:8080/health

## 📖 API 文档概览

### 认证相关 (Authentication)
- `POST /api/auth/login` - 用户登录
- `POST /api/auth/register` - 用户注册

### 用户管理 (Users)
- `GET /api/users/profile` - 获取用户资料
- `PUT /api/users/profile` - 更新用户资料

### Agent 管理 (Agents)
- `GET /api/agents` - 获取 Agent 列表
- `POST /api/agents` - 创建 Agent
- `GET /api/agents/{id}` - 获取 Agent 详情
- `PUT /api/agents/{id}` - 更新 Agent
- `DELETE /api/agents/{id}` - 删除 Agent

### 对话管理 (Conversations)
- `GET /api/conversations` - 获取对话列表
- `POST /api/conversations` - 创建对话
- `GET /api/conversations/{id}` - 获取对话详情
- `DELETE /api/conversations/{id}` - 删除对话

### 消息管理 (Messages)
- `GET /api/conversations/{id}/messages` - 获取消息列表
- `POST /api/conversations/{id}/messages` - 发送消息

### Coze 集成 (Coze)
- `GET /api/coze/token` - 获取 Coze 访问令牌

## 🔐 认证说明

大部分 API 需要 JWT 认证。在 Swagger UI 中：

1. 首先调用 `/auth/login` 或 `/auth/register` 获取 token
2. 点击页面右上角的 "Authorize" 按钮
3. 在弹出的对话框中输入：`Bearer <your_token>`
4. 点击 "Authorize" 完成认证

## 📋 请求示例

### 用户注册
```json
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "123456",
  "nickname": "测试用户"
}
```

### 创建 Agent
```json
{
  "name": "我的AI助手",
  "description": "一个智能对话助手",
  "avatar": "https://example.com/avatar.jpg",
  "prompt": "你是一个友好的AI助手",
  "config": "{\"temperature\": 0.7}"
}
```

## 🔄 更新文档

当修改 API 接口后，需要重新生成 Swagger 文档：

```bash
# 安装 swag 工具（如果未安装）
go install github.com/swaggo/swag/cmd/swag@latest

# 生成文档
swag init -g cmd/main.go

# 重启服务器
go run cmd/main.go
```

## 📝 Swagger 注释规范

在控制器函数上添加 Swagger 注释：

```go
// CreateAgent 创建Agent
// @Summary 创建Agent
// @Description 创建新的AI Agent
// @Tags Agent
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body CreateAgentRequest true "Agent信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /agents [post]
func CreateAgent(c *gin.Context) {
    // 实现代码...
}
```

## 🛠️ 文档结构

生成的文档文件位于：
- `docs/swagger.json` - JSON 格式的 API 文档
- `docs/swagger.yaml` - YAML 格式的 API 文档  
- `docs/docs.go` - Go 代码格式的文档（自动生成）

## 📞 支持

如有问题，请查看：
- Swagger 官方文档: https://swagger.io/docs/
- gin-swagger 文档: https://github.com/swaggo/gin-swagger 