build:
	@go build -v -o ./bin/abaxx-id-go .

run: build
	@./bin/abaxx-id-go

test:
	@go test ./...

clean:
	@rm -rf bin
