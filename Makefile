APP_NAME=crud-api
SRC=./cmd/server/server.go
BUILD_DIR=./build
MOCKERY := mockery

.PHONY: mocks docs test

build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(SRC)

run: build
	source .env && ${BUILD_DIR}/$(APP_NAME)

mocks:
	mockery --name PortfolioRepository --dir internal/repository --output internal/repository/mocks

test:
	go test ./.../

docs:
	rm -rf docs && swag init -g cmd/server/server.go

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)