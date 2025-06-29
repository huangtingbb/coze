# Coze Agent Platform

基于 Gin 框架的 Coze Agent 中台微服务系统

## 功能特性

- 🚀 基于 Gin 的高性能 RESTful API
- 🔐 JWT 身份验证
- 📝 Swagger API 文档
- 🗄️ MySQL + GORM ORM
- 📦 Redis 缓存支持
- 🔍 结构化日志记录
- 🛡️ 统一错误处理
- 📱 CORS 跨域支持

## 项目结构

```
coze-agent-platform/
├── cmd/main.go                # 应用入口
├── config/                    # 配置管理
│   ├── config.go             # 配置结构体和初始化
│   └── config.yaml           # 配置文件
├── controllers/              # 控制器层
│   ├── auth.go              # 认证控制器(登录/注册)
│   ├── user.go              # 用户控制器
│   ├── agent.go             # Agent控制器
│   └── conversation.go      # 对话控制器(占位符)
├── services/                # 服务层
│   ├── user_service.go      # 用户服务
│   └── agent_service.go     # Agent服务
├── models/                  # 数据模型层
│   ├── database.go          # 数据库连接
│   ├── user.go              # 用户模型
│   ├── agent.go             # Agent模型
│   ├── conversation.go      # 对话模型
│   └── message.go           # 消息模型
├── middleware/              # 中间件
│   ├── logger.go            # 日志中间件
│   ├── recovery.go          # 恢复中间件
│   ├── cors.go              # CORS中间件
│   └── jwt.go               # JWT认证中间件
├── routers/                 # 路由配置
│   └── router.go            # 路由设置
├── utils/                   # 工具类
│   ├── logger.go            # 日志工具
│   ├── jwt.go               # JWT工具
│   └── response.go          # 统一响应工具
├── Makefile                 # 构建脚本
├── README.md                # 项目文档
├── go.mod                   # Go模块文件  
└── go.sum                   # 依赖锁定文件
```

## 快速开始

### 1. 安装依赖

```bash
make deps
```

### 2. 配置数据库

修改 `config/config.yaml` 中的数据库配置：

```yaml
database:
  host: "localhost"
  port: 3306
  username: "root"
  password: "your-password"
  dbname: "coze_agent"
```

### 3. 运行服务

```bash
make run
```

或者开发模式运行（包含依赖安装和文档生成）：

```bash
make dev
```

### 4. 访问服务

- API 服务: http://localhost:8080
- 健康检查: http://localhost:8080/health
- Swagger 文档: http://localhost:8080/swagger/index.html

## API 接口

### 认证接口

- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/register` - 用户注册

### 用户接口

- `GET /api/v1/users/profile` - 获取用户资料
- `PUT /api/v1/users/profile` - 更新用户资料

### Agent 接口

- `GET /api/v1/agents` - 获取 Agent 列表
- `POST /api/v1/agents` - 创建 Agent
- `GET /api/v1/agents/:id` - 获取 Agent 详情
- `PUT /api/v1/agents/:id` - 更新 Agent
- `DELETE /api/v1/agents/:id` - 删除 Agent

## 开发命令

```bash
# 构建项目
make build

# 运行项目
make run

# 运行测试
make test

# 生成 Swagger 文档
make swagger

# 清理构建文件
make clean
```

## 配置说明

系统配置通过环境变量或配置文件管理，支持以下配置项：

- `APP_NAME`: 应用名称
- `APP_PORT`: 服务端口
- `APP_MODE`: 运行模式 (debug/release)
- `DB_HOST`: 数据库主机
- `DB_PORT`: 数据库端口
- `DB_USERNAME`: 数据库用户名
- `DB_PASSWORD`: 数据库密码
- `DB_NAME`: 数据库名称
- `JWT_SECRET`: JWT 密钥
- `JWT_EXPIRE`: JWT 过期时间（秒）

## 依赖项

- [Gin](https://github.com/gin-gonic/gin) - HTTP Web 框架
- [GORM](https://gorm.io/) - ORM 库
- [JWT-Go](https://github.com/golang-jwt/jwt) - JWT 认证
- [Viper](https://github.com/spf13/viper) - 配置管理
- [Logrus](https://github.com/sirupsen/logrus) - 结构化日志
- [Swagger](https://github.com/swaggo/gin-swagger) - API 文档

## 许可证

MIT License 