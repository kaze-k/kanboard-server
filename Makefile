APP_NAME=kanboard
BUILD_DIR=build
BUILD_WIN_DIR=build-win
CONFIG_DIR=config
CONFIG_FILE=config.template.toml
CONFIG_NAME=config.toml

.DEFAULT: all

all:
	@make clean
	@make build
	@make run

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME)
	cp $(CONFIG_DIR)/$(CONFIG_FILE) $(BUILD_DIR)/$(CONFIG_NAME)

build-win:
	mkdir -p $(BUILD_WIN_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_WIN_DIR)/$(APP_NAME).exe
	cp $(CONFIG_DIR)/$(CONFIG_FILE) $(BUILD_WIN_DIR)/$(CONFIG_NAME)

clean:
	rm -rf $(BUILD_DIR)
	rm -rf ${BUILD_WIN_DIR}

run:
	./$(BUILD_DIR)/$(APP_NAME)

dev:
	fresh
