############################
#╔════════════════════════╗#
#║ DEFAULT CONFIGURATION  ║#
#╚════════════════════════╝#
############################
.PHONY: help build deploy
.SILENT:

.DEFAULT_GOAL = help
ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
$(eval $(ARGS):;@:)

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

SERVICES := $(wildcard ../services/*)

clean:
	rm -rf output/*

package:
	upx --best --lzma output/**/*;

build-framework:
	echo Build Kitsune;
	CGO_ENABLED=0 go build -trimpath -buildvcs=false -ldflags="-s -w \
		-X github.com/kodmain/KitsuneFramework/internal/env.BUILD_VERSION=$$(git describe --tags --abbrev=0) \
		-X github.com/kodmain/KitsuneFramework/internal/env.BUILD_COMMIT=$$(git rev-parse --short HEAD) \
		-X github.com/kodmain/KitsuneFramework/internal/env.BUILD_APP_NAME=kitsune" \
		-o output/bin/kitsune $(CURDIR)/../cmd/main.go;

build-service:
	echo Build service $(ARGS);
	CGO_ENABLED=0 go build -trimpath -buildvcs=false -ldflags="-s -w \
		-X github.com/kodmain/KitsuneFramework/internal/env.BUILD_VERSION=$$(git describe --tags --abbrev=0) \
		-X github.com/kodmain/KitsuneFramework/internal/env.BUILD_COMMIT=$$(git rev-parse --short HEAD) \
		-X github.com/kodmain/KitsuneFramework/internal/env.BUILD_APP_NAME=$(ARGS)" \
		-o output/services/$(ARGS) $(CURDIR)/../services/$(ARGS)/main.go;

build-services:
	@for service in $(SERVICES); do \
		make build-service $$(basename $$service);\
	done

project: clean build-framework build-services package ## Build all for current env
