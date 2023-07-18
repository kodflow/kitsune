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
	@grep -E '^[a-zA-Z_-]+:.*?##[\s]?.*$$' Makefile | awk 'BEGIN {FS = ":.*?##"}; {gsub(/(^ +| +$$)/, "", $$2);printf "╟─[ \033[36m%-$(PADDING)s\033[0m %s\n", $$1, "] "$$2}'
	echo '╚──────────────────────────────────────────────────>'
	echo ''

aio: ## Start services in portable version
	docker compose -f .github/build/compose.yml --profile=services up --build

run: ## Run services in portable version
	docker compose -f .github/build/compose.yml --profile=standalone run --build --rm kitsune.$(ARGS)

#test: ## Run services in portable version
#	docker compose -f .github/build/compose.yml --profile=standalone build
#	docker compose -f .github/build/compose.yml --profile=standalone run --rm kitsune.$(ARGS)

tests:
	find src -name "*.go" | grep -v "_test.go$$" | while read -r file; do test -f "$${file%.*}_test.go" || echo "package $$(grep -m 1 'package ' $$file | awk '{print $$2}')\n\nimport \"testing\"\n\nfunc TestNotExistInThisFile$$(basename $$file .go)(t *testing.T) {}\n" > "$${file%.*}_test.go"; done
	go test -v `go list ./...` -coverprofile=coverage.txt -covermode=atomic
	find . -name "*_test.go" | xargs grep -l "func TestNotExistInThisFile" | xargs rm 

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
	@for service in $(SERVICES); do \
		make build-service $$(basename $$service);\
	done

build-framework:
	echo Build Kitsune;
	@for os in $(OS); do \
		for arch in $(ARCH); do \
			ldflags="-s -w \
				-X github.com/kodmain/kitsune/src/internal/env.BUILD_VERSION=$$VERSION \
				-X github.com/kodmain/kitsune/src/internal/env.BUILD_COMMIT=$$(git rev-parse --short HEAD) \
				-X github.com/kodmain/kitsune/src/internal/env.BUILD_APP_NAME=kitsune"; \
			for binary in $$(find .generated -type f -name "*$$os-$$arch.sha1" | awk -F "/" '{print $$3}' | awk -F "-" '{print $$1}' | sort | uniq); do \
				cap_binary=$$(echo $$binary | tr '[:lower:]' '[:upper:]'); \
				ldflags+=" -X github.com/kodmain/kitsune/src/internal/env.BUILD_SERVICE_$$cap_binary=$$(cat .generated/services/$$binary-$$os-$$arch.sha1)"; \
			done; \
			CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -trimpath -buildvcs=false -ldflags="$$ldflags" \
				-o .generated/bin/kitsune-$$os-$$arch $(CURDIR)/src/cmd/main.go; \
				chmod +x .generated/bin/kitsune-$$os-$$arch; \
		done \
	done
	find .generated | grep "sha1" | xargs rm 

build-service:
	echo Build service $(ARGS);
	@for os in $(OS); do \
		for arch in $(ARCH); do \
			CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -trimpath -buildvcs=false -ldflags="-s -w \
			-X github.com/kodmain/kitsune/src/internal/env.BUILD_VERSION=$$VERSION \
			-X github.com/kodmain/kitsune/src/internal/env.BUILD_COMMIT=$$(git rev-parse --short HEAD) \
			-X github.com/kodmain/kitsune/src/internal/env.BUILD_APP_NAME=$(ARGS)" \
			-o .generated/services/$(ARGS)-$$os-$$arch $(CURDIR)/src/services/$(ARGS)/main.go; \
			chmod +x .generated/services/$(ARGS)-$$os-$$arch; \
			sha1sum .generated/services/$(ARGS)-$$os-$$arch | awk '{ print $$1 }'  | tr -d '\n' > .generated/services/$(ARGS)-$$os-$$arch.sha1; \
		done \
	done

package:
	find .generated -type f | xargs upx --best --lzma;
