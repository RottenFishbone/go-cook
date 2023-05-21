.PHONY: cook
cook:
	make web
	make fmt
	go build ./cmd/cook

.PHONY: web
web:
	@pushd internal/web/ && npm install && npm run build

.PHONY: fmt
fmt:
	@find . -name "*.go" -type f -exec gofmt -w $$(dirname {}) \;

.PHONY: canonical
canonical:
	go build ./internal/cmd/test_gen
	./test_gen
	gofmt -w -s canonical_test.go

.PHONY: clean
clean:
	rm -f cook
	go clean
	go clean ./cmd/*
	rm -rf ./internal/web/dist/*
