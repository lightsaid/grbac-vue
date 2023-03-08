## up: 在后台启动 docker compose 构建基础环境
up:
	docker-compose up -d

## down: 停止 docker compose
down:
	docker-compose down

## serve: 启动服务
serve:
	go run ./cmd/server/*.go

## mailer: 启动发送邮件服务
mailer: 
	go run ./cmd/mailer/*.go

## doc: swag init
doc:
	swag init -d ./cmd/server

.PHONY: up down serve mailer doc 