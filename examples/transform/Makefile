MODULE ?= "acme.com/k8s/transform"

.PHONY: build
build: clean
ifndef VERSION
	$(error VERSION must be specified)
endif
	@CGO_ENABLED=0 go build -ldflags="-X main.module=${MODULE}@${VERSION}"  -o ./bin/ ./cmd/*

.PHONY: clean
clean:
	@rm -rf ./bin
