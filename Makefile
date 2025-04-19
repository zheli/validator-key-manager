.PHONY: format lint test docker-build

format:
	go fmt ./...

lint:
	revive -config revive.toml ./...

test:
	go test -v -race -cover ./...

docker-build:
	docker build -t validator-key-manager .
