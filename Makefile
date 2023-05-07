APP_NAME=crud-api
SRC=./main.go
BUILD_DIR=./build
MOCKERY := mockery

.PHONY: mocks

all: build

build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(SRC)

run: build
	source .env && ${BUILD_DIR}/$(APP_NAME)

mocks:
	mockery --name PortfolioRepository --dir repository --output repository/mocks

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)