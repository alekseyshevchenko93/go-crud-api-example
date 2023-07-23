APP_NAME=crud-api
SRC=./cmd/server/server.go
BUILD_DIR=./bin
MOCKERY := mockery

.PHONY: build build-prod mocks docs test

build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(SRC)

run: build
	source .env && ${BUILD_DIR}/$(APP_NAME)

mocks:
	mockery --name PortfolioRepository --dir internal/repository --output internal/repository/mocks

test:
	go test -race -v  ./.../

docs:
	rm -rf docs && swag init -g cmd/server/server.go

coverage:
	go test -coverprofile=coverage.out ./... ;    go tool cover -html=coverage.out
	
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)