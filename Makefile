APP_NAME := mdkp-backend
PKG_MAIN := ./cmd
CONFIG_PATH ?= ./config/config.yaml
GO ?= go

.PHONY: run test build build-win build-linux build-mac

run:
	$(GO) run $(PKG_MAIN) -config-path=$(CONFIG_PATH)

test:
	./run_tests.sh >> output.txt

# Универсальная сборка: make build OS=WIN|Linux|MacOS [ARCH=amd64]
build:
	@if [ -z "$(OS)" ]; then echo "Usage: make build OS=WIN|Linux|MacOS [ARCH=amd64]"; exit 2; fi
	@case "$(OS)" in \
		WIN)   GOOS=windows ;; \
		Linux) GOOS=linux ;; \
		MacOS) GOOS=darwin ;; \
		*) echo "Unknown OS=$(OS) (use WIN|Linux|MacOS)"; exit 2 ;; \
	esac; \
	ARCH=$${ARCH:-amd64}; \
	OUT=dist/$(APP_NAME)-$${GOOS}-$${ARCH}; \
	if [ "$${GOOS}" = "windows" ]; then OUT=$${OUT}.exe; fi; \
	mkdir -p dist; \
	echo "Building $$OUT"; \
	GOOS=$${GOOS} GOARCH=$${ARCH} $(GO) build -o $$OUT $(PKG_MAIN)

build-win:
	$(MAKE) build OS=WIN

build-linux:
	$(MAKE) build OS=Linux

build-mac:
	$(MAKE) build OS=MacOS

runBat:
	./main