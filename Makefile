postgresinit:
	docker run --name postgres16 -p 5433:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:16-alpine
postgres:
	docker exec -it postgres16 bin/bash

createdb:
	docker exec -it postgres16 createdb --username=postgres --owner=postgres authDatabase

.PHONY: postgresinit postgres createdb

protoc:
	protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative ./pkg/protobuf/userservice/gw/user.proto
userService:
	go run cmd//user/main.go
authService:
	go run cmd//auth/main.go


