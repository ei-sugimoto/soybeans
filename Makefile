.PHONY build:
build:
	@echo "build..."
	@go build .
	@sudo cp ./soybeans /usr/local/bin/soybeans