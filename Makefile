run-lint:
	golangci-lint run ./...

protoc-user:
	protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative ./pkg/protobuf/userservice/gw/user.proto
protoc-auth:
	protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative ./pkg/protobuf/authorizationservice/gw/authorization_service.proto
user-service:
	go run cmd//user/main.go
auth-service:
	go run cmd//auth/main.go
gallery-service:
	go run cmd//auth/main.go

gallery-migrations-up:
	migrate -path database/gallery_migrations/ -database "postgresql://postgres:postgres@localhost:5432/galleryService?sslmode=disable" -verbose up
user-migrations-up:
	migrate -path database/user_migrations/ -database "postgresql://postgres:postgres@localhost:5432/authDatabase?sslmode=disable" -verbose up
user-migrations-down:
	migrate -path database/user_migrations/ -database "postgresql://postgres:postgres@localhost:5432/authDatabase?sslmode=disable" -verbose down
gallery-migrations-down:
	migrate -path database/gallery_migrations/ -database "postgresql://postgres:postgres@localhost:5432/galleryService?sslmode=disable" -verbose down
