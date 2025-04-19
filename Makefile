.PHONY: format lint test

format:
	go fmt ./...

lint:
	revive -config revive.toml ./...

test:
	go test -v -race -cover ./...
