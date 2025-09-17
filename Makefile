SWAG_DIRS=./cmd/api,./internal/transport/http,./internal/domain

bin/swagger:
	GOBIN=$(shell pwd)/bin go install github.com/swaggo/swag/cmd/swag@v1.8.12

swagger: bin/swagger
	./bin/swag init -g main.go -o internal/transport/http/docs -d $(SWAG_DIRS)
