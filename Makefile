build:
	@$(shell go env GOPATH)/bin/packr2
	@go build -o bin/bliki main.go
	@$(shell go env GOPATH)/bin/packr2 clean
	@echo "OK"

clean:
	@rm -rf ./bin

packr:
	@cd && GO111MODULE=on go get github.com/gobuffalo/packr/v2/packr2@v2.8.1
	@echo "OK"
