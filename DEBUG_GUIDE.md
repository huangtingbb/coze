# Coze Agent 微服务调试指南

## 概述

本项目已配置了完整的 VS Code 调试环境，支持多种调试场景。本指南将帮助您快速上手调试功能。

## 前置条件

### 必需的 VS Code 扩展

项目已配置推荐扩展列表，请确保安装以下关键扩展：

1. **Go** (`golang.go`) - Go 语言支持
2. **Docker** (`ms-azuretools.vscode-docker`) - Docker 容器支持
3. **REST Client** (`humao.rest-client`) - API 测试
4. **YAML** (`redhat.vscode-yaml`) - YAML 文件支持

VS Code 会自动提示安装推荐扩展，或者您可以手动安装。

### 必需的 Go 工具

确保安装以下 Go 工具：

```bash
# 安装 Delve 调试器
go install github.com/go-delve/delve/cmd/dlv@latest

# 安装 goimports
go install golang.org/x/tools/cmd/goimports@latest

# 安装 golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 安装 swag (Swagger 文档生成)
go install github.com/swaggo/swag/cmd/swag@latest
```

## 调试配置说明

### 1. Launch Server (推荐)

**用途**: 开发环境下启动服务器进行调试

**配置**:
- 自动设置 `GIN_MODE=debug`
- 启用详细日志输出
- 预构建任务支持

**使用方法**:
1. 在代码中设置断点
2. 按 `F5` 或选择 "Launch Server" 配置
3. 服务器将在调试模式下启动

### 2. Launch Server (Release Mode)

**用途**: 生产环境模式下测试

**配置**:
- 设置 `GIN_MODE=release`
- 简化日志输出
- 性能优化模式

### 3. Debug Current File

**用途**: 调试当前打开的单个 Go 文件

**使用方法**:
1. 打开要调试的 Go 文件
2. 设置断点
3. 选择 "Debug Current File" 配置

### 4. Debug Test in Current File

**用途**: 调试当前文件中的测试函数

**使用方法**:
1. 修改配置中的 `TestFunction` 为实际测试函数名
2. 或者使用 Go 扩展的 CodeLens 功能直接调试测试

### 5. Debug All Tests

**用途**: 运行并调试所有测试

**配置**:
- 启用 `-test.v` 详细输出
- 支持断点调试

### 6. Attach to Process

**用途**: 附加到正在运行的 Go 进程

**使用方法**:
1. 先启动应用程序
2. 选择 "Attach to Process"
3. 从进程列表中选择目标进程

### 7. Connect to Remote Delve

**用途**: 连接到远程 Delve 调试服务器

**配置**:
- 默认端口: 2345
- 默认主机: 127.0.0.1

## 常用调试技巧

### 设置断点

1. **行断点**: 在代码行号左侧点击
2. **条件断点**: 右键点击断点，设置条件
3. **日志断点**: 在断点处输出日志而不暂停

### 调试面板功能

- **变量**: 查看当前作用域中的变量
- **监视**: 添加表达式进行监视
- **调用堆栈**: 查看函数调用堆栈
- **断点**: 管理所有断点

### 调试控制台

支持以下命令：
- `p <variable>`: 打印变量值
- `pp <variable>`: 格式化打印变量
- `locals`: 显示本地变量
- `args`: 显示函数参数

## 项目特定的调试设置

### 环境变量

调试配置中设置了以下环境变量：

```json
{
    "GIN_MODE": "debug",     // Gin 框架调试模式
    "GO_ENV": "development"  // 应用环境标识
}
```

### 配置文件

确保 `config/config.yaml` 中的配置适合调试：

```yaml
app:
  mode: debug
  port: 8080

database:
  host: localhost
  port: 3306
  # ... 其他配置
```

## 快速开始

### 1. 启动调试会话

```bash
# 1. 克隆项目并进入目录
cd coze-agent-platform

# 2. 安装依赖
go mod tidy

# 3. 启动 VS Code
code .

# 4. 按 F5 开始调试
```

### 2. 测试 API

在 VS Code 中创建 `.http` 文件进行 API 测试：

```http
### 健康检查
GET http://localhost:8080/health

### 用户登录
POST http://localhost:8080/api/auth/login
Content-Type: application/json

{
    "username": "admin",
    "password": "123456"
}
```

### 3. 查看 Swagger 文档

启动服务器后，访问：http://localhost:8080/swagger/index.html

## 常见问题解决

### 1. 调试器无法启动

**问题**: `could not launch process: fork/exec`

**解决方案**:
```bash
# 检查 Go 版本
go version

# 重新安装 Delve
go install github.com/go-delve/delve/cmd/dlv@latest
```

### 2. 断点不生效

**问题**: 断点显示为灰色或被跳过

**解决方案**:
1. 检查是否启用了代码优化
2. 确保文件已保存
3. 重新启动调试会话

### 3. 找不到模块

**问题**: `cannot find module`

**解决方案**:
```bash
# 清理模块缓存
go clean -modcache

# 重新下载依赖
go mod download
```

## 高级调试技巧

### 1. 远程调试

在服务器上启动 Delve 服务：

```bash
# 构建调试版本
go build -gcflags="all=-N -l" -o bin/coze-agent-platform ./cmd/main.go

# 启动 Delve 服务器
dlv --listen=:2345 --headless=true --api-version=2 exec ./bin/coze-agent-platform
```

### 2. 容器调试

使用 Docker 进行调试：

```dockerfile
# 在 Dockerfile 中添加调试支持
FROM golang:1.21-alpine AS debug

RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /app
COPY . .

RUN go build -gcflags="all=-N -l" -o bin/coze-agent-platform ./cmd/main.go

EXPOSE 8080 2345

CMD ["dlv", "--listen=:2345", "--headless=true", "--api-version=2", "exec", "./bin/coze-agent-platform"]
```

### 3. 性能分析

启用 pprof 进行性能分析：

```go
import (
    _ "net/http/pprof"
    "net/http"
)

// 在 main.go 中添加
go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```

访问 `http://localhost:6060/debug/pprof/` 查看性能数据。

## 总结

本调试配置提供了完整的开发和调试支持，涵盖了从本地开发到远程调试的各种场景。通过合理使用这些配置，可以大大提高开发效率和代码质量。

如有问题，请参考 VS Code Go 扩展文档或提交 Issue。
