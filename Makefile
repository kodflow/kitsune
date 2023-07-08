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
	go test -v `go list ./...` -coverprofile=coverage.txt -covermode=atomic

update: ## Install/Update vendor
	echo "Update all dependencies"
	go get -u ./...
	go mod tidy

build: update build-framework build-services ## Build all services

build-services:
	@for service in $(SERVICES); do \
		make build-service $$(basename $$service);\
	done

binary-only:
	find . | grep -E 'sha1|md5' | xargs rm;

build-framework:
	echo Build Kitsune;
	CGO_ENABLED=0 go build -trimpath -buildvcs=false -ldflags="-s -w \
		-X github.com/kodmain/kitsune/src/internal/env.BUILD_VERSION=$$VERSION \
		-X github.com/kodmain/kitsune/src/internal/env.BUILD_COMMIT=$$(git rev-parse --short HEAD) \
		-X github.com/kodmain/kitsune/src/internal/env.BUILD_APP_NAME=kitsune" \
		-o .generated/bin/kitsune $(CURDIR)/src/cmd/main.go;
	sha1sum .generated/bin/kitsune | awk '{ print $$1 }' | tr -d '\n' >> .generated/bin/kitsune.sha1;
	md5sum  .generated/bin/kitsune | awk '{ print $$1 }' | tr -d '\n' >> .generated/bin/kitsune.md5;

build-service:
	echo Build service $(ARGS);
	CGO_ENABLED=0 go build -trimpath -buildvcs=false -ldflags="-s -w \
		-X github.com/kodmain/kitsune/src/internal/env.BUILD_VERSION=$$VERSION \
		-X github.com/kodmain/kitsune/src/internal/env.BUILD_COMMIT=$$(git rev-parse --short HEAD) \
		-X github.com/kodmain/kitsune/src/internal/env.BUILD_APP_NAME=$(ARGS)" \
		-o .generated/services/$(ARGS) $(CURDIR)/src/services/$(ARGS)/main.go;
	sha1sum .generated/services/$(ARGS) | awk '{ print $$1 }'  | tr -d '\n' > .generated/services/$(ARGS).sha1;
	md5sum  .generated/services/$(ARGS) | awk '{ print $$1 }'  | tr -d '\n' > .generated/services/$(ARGS).md5;

package:
	find .generated -type f ! -name "*.md5" ! -name "*.sha1" | xargs upx --best --lzma;
