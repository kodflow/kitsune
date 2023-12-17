############################
#╔════════════════════════╗#
#║ DEFAULT CONFIGURATION  ║#
#╚════════════════════════╝#
############################
.PHONY: help build deploy
.SILENT:

.DEFAULT_GOAL = help
ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
SERVICES := $(wildcard src/services/*)
OS := linux darwin
ARCH := amd64 arm64
$(eval $(ARGS):;@:)
$(eval GWD=$(shell git rev-parse --show-toplevel 2>/dev/null || pwd))
$(shell git config core.hooksPath $(GWD)/.github/hooks 2>/dev/null || true)
$(shell chmod +x $(GWD)/.github/hooks/*)

help: #Pour générer automatiquement l'aide ## Display all commands available
	$(eval PADDING=$(shell grep -x -E '^[a-zA-Z_-]+:.*?##[\s]?.*$$' Makefile | awk '{ print length($$1)-1 }' | sort -n | tail -n 1))
	clear
	echo '╔──────────────────────────────────────────────────╗'
	echo '║ ██╗  ██╗███████╗██╗     ██████╗ ███████╗██████╗  ║'
	echo '║ ██║  ██║██╔════╝██║     ██╔══██╗██╔════╝██╔══██╗ ║'
	echo '║ ███████║█████╗  ██║     ██████╔╝█████╗  ██████╔╝ ║'
	echo '║ ██╔══██║██╔══╝  ██║     ██╔═══╝ ██╔══╝  ██╔══██╗ ║'
	echo '║ ██║  ██║███████╗███████╗██║     ███████╗██║  ██║ ║'
	echo '║ ╚═╝  ╚═╝╚══════╝╚══════╝╚═╝     ╚══════╝╚═╝  ╚═╝ ║'
	echo '╟──────────────────────────────────────────────────╝'
	grep -E '^[a-zA-Z_-]+:.*?##[\s]?.*$$' Makefile | awk 'BEGIN {FS = ":.*?##"}; {gsub(/(^ +| +$$)/, "", $$2);printf "╟─[ \033[36m%-$(PADDING)s\033[0m %s\n", $$1, "] "$$2}'
	echo '╚──────────────────────────────────────────────────>'
	echo ''

aio: ## Start services in portable version
	docker compose -f .github/build/compose.yml --profile=services up --build

run: ## Run services in portable version
	docker compose -f .github/build/compose.yml --profile=standalone run --build --rm kitsune.$(ARGS)

#test: ## Run services in portable version
#	docker compose -f .github/build/compose.yml --profile=standalone build
#	docker compose -f .github/build/compose.yml --profile=standalone run --rm kitsune.$(ARGS)

ssl:
	echo "Generating SSL certificates..."
	rm -rf ./.generated/ssl/*
	mkdir -p ./.generated/ssl/
	openssl genpkey -algorithm RSA -out ./.generated/ssl/localhost.key
	openssl req -new -key ./.generated/ssl/localhost.key -out ./.generated/ssl/localhost.csr -subj "/CN=localhost"
	openssl req -x509 -days 365 -key ./.generated/ssl/localhost.key -in ./.generated/ssl/localhost.csr -out ./.generated/ssl/localhost.crt
	openssl req -x509 -days 365 -nodes -newkey rsa:2048 -keyout ./.generated/ssl/localhost-ca.key -out ./.generated/ssl/localhost-ca.crt -subj "/CN=Certificate Authority"
	echo "SSL certificates generated successfully!"
	
tests: install-gotestsum
	gotestsum -- -v $(go list ./... | grep -vE "src/services|generated") -run 'Test[^Pb]' -coverprofile=coverage.txt -covermode=atomic || true

install-gotestsum:
	if ! command -v gotestsum > /dev/null; then \
		go install gotest.tools/gotestsum@latest; \
	fi

update: ## Install/Update vendor
	echo "Update all dependencies"
	go get -u ./...
	go mod vendor

build: clear update build-services build-framework ## Build all services

clear:
	rm -rf .generated

binary-only:
	find . | grep -E 'sha1|md5' | xargs rm;

build-services:
	for service in $(SERVICES); do \
		make build-service $$(basename $$service);\
	done

build-framework:
	echo Build Kitsune;
	for os in $(OS); do \
		for arch in $(ARCH); do \
			ldflags="-s -w \
				-X github.com/kodflow/kitsune/src/config/config.BUILD_VERSION=$$VERSION \
				-X github.com/kodflow/kitsune/src/config/config.BUILD_COMMIT=$$(git rev-parse --short HEAD) \
				-X github.com/kodflow/kitsune/src/config/config.BUILD_APP_NAME=kitsune"; \
			CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -trimpath -buildvcs=false -ldflags="$$ldflags" \
				-o .generated/bin/kitsune-$$os-$$arch $(CURDIR)/src/cmd/main.go; \
				chmod +x .generated/bin/kitsune-$$os-$$arch; \
		done \
	done
	find .generated | grep "sha1" | xargs rm 

build-service:
	echo Build service $(ARGS);
	for os in $(OS); do \
		for arch in $(ARCH); do \
			CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -trimpath -buildvcs=false -ldflags="-s -w \
			-X github.com/kodflow/kitsune/src/config.BUILD_VERSION=$$VERSION \
			-X github.com/kodflow/kitsune/src/config.BUILD_COMMIT=$$(git rev-parse --short HEAD) \
			-X github.com/kodflow/kitsune/src/config.BUILD_APP_NAME=$(ARGS)" \
			-o .generated/services/$(ARGS)-$$os-$$arch $(CURDIR)/src/services/$(ARGS)/main.go; \
			chmod +x .generated/services/$(ARGS)-$$os-$$arch; \
			sha1sum .generated/services/$(ARGS)-$$os-$$arch | awk '{ print $$1 }'  | tr -d '\n' > .generated/services/$(ARGS)-$$os-$$arch.sha1; \
		done \
	done

# Build and push multi-architecture Docker images using buildx
build-images: build
	docker buildx create --use
	docker buildx inspect --bootstrap
	for file in .generated/services/*; do \
		name=$$(basename $$file | cut -d '-' -f1); \
		os=$$(basename $$file | cut -d '-' -f2); \
		arch=$$(basename $$file | cut -d '-' -f3); \
		\
		if [ "$$os" = "linux" ]; then \
			echo "Construit l'image pour : name=$$name, os=$$os, arch=$$arch, file=$$file"; \
			docker buildx build \
			--platform $$os/$$arch \
			--file .github/local/Dockerfile.debug \
			--tag kodflow/debug:$$name \
			--build-arg FILE=$$file \
			--push .; \
		fi \
	done;
	for instance in $$(docker ps | grep 'moby/buildkit' | awk '{print $$1}'); do \
		docker kill $$instance; \
		docker rm $$instance; \
	done

build-images-clear:
	docker image rm $$(docker images | grep 'moby/buildkit' | awk '{print $$3}')

run-images-with-publish: build-images run-images

run-images:
	docker compose -f .github/local/compose.yml pull
	docker compose -f .github/local/compose.yml up

generate:
#protoc --proto_path=$(CURDIR)/src/internal/data --go_out=$(CURDIR) $(CURDIR)/src/internal/data/proto/*
#protoc --proto_path=$(CURDIR) --go_out=$(CURDIR) $(CURDIR)/src/internal/core/server/transport/proto/*
	find ${CURDIR} -name '*.proto' -exec protoc --proto_path=${CURDIR} --go_out=${CURDIR} {} \;

package:
	find .generated -type f | xargs upx --best --lzma;

copy-to-infra:
	rm -rf ~/Documents/Projects/Infrastructure/organizations/IEF2I/IT-F2I:387672757226/.server/bench/usr/local/bin
	mkdir -p ~/Documents/Projects/Infrastructure/organizations/IEF2I/IT-F2I:387672757226/.server/bench/usr/local/bin
	find .generated | grep "linux-arm64" | xargs -I {} cp {} ~/Documents/Projects/Infrastructure/organizations/IEF2I/IT-F2I:387672757226/.server/bench/usr/local/bin
