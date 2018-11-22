all:
	@make build

build:
	@go build -o build/litemysql cmd/litemysql/main.go

install:
	@make build
	@cp -rf build/litemysql /usr/local/bin/litemysql

clean:
	@rm -rf build/

.PHONY: build clean
