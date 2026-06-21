APP_NAME=integritypos
MAIN_PATH=./cmd/server

.PHONY: build-windows build-linux clean

build-windows:
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o bin/$(APP_NAME).exe $(MAIN_PATH)

build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/$(APP_NAME) $(MAIN_PATH)

clean:
	rm -rf bin/
