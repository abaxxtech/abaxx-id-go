build:
	@go build -ldflags="-s -w -buildid=" -v -o ./bin/abaxx-id ./cmd/abaxx-id/.

run: build
	@./bin/abaxx-id

test:
	@go test ./...

clean:
	@rm -rf bin

publish:
	@gh release upload v0.0.1 ./bin/abaxx-id
