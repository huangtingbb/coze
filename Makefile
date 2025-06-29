.PHONY: build run clean test deps swagger

# 构建
build:
	go build -o bin/coze-agent-platform cmd/main.go

# 运行
run:
	go run cmd/main.go

# 清理
clean:
	rm -rf bin/

# 测试
test:
	go test -v ./...

# 安装依赖
deps:
	go mod tidy
	go mod download

# 生成Swagger文档
swagger:
	swag init -g cmd/main.go -o docs/

# 开发环境启动
dev: deps swagger run 