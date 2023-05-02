.PHONY: cook
cook:
	make web
	go build ./cmd/cook

.PHONY: web
web:
	@pushd internal/web/ && npm install && npm run build

.PHONY: fmt
fmt:
	@find . -name "*.go" -type f -exec gofmt -w $$(dirname {}) \;

.PHONY: clean
clean:
	rm cook
	go clean
	go clean ./cmd/*
	rm -rf ./internal/web/dist/*
