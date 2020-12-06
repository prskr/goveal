VERSION = $(shell git describe --dirty --tags --always)
REPO = github.com/baez90/goveal
BUILD_PATH = $(REPO)/cmd/goveal
PKGS = $(shell go list ./... | grep -v /vendor/)
TEST_PKGS = $(shell find . -type f -name "*_test.go" -printf '%h\n' | sort -u)
GOARGS = GOOS=linux GOARCH=amd64
GO_BUILD_ARGS = -ldflags="-w -s"
BINARY_NAME = goveal
DIR = $(dir $(realpath $(firstword $(MAKEFILE_LIST))))
DEBUG_PORT = 2345

REVEALJS_VERSION = 4.0.2
GORELEASER_VERSION = 0.132.1

.PHONY: all clean clean-all clean-vendor rebuild format revive test deps compile run debug watch watch-test cloc docs serve-docs serve-godoc ensure-revive ensure-reflex ensure-delve ensure-godoc ensure-pkger ensure-goreleaser

export CGO_ENABLED:=0

all: format compile

clean-all: clean clean-vendor

rebuild: clean format compile

format:
	@go fmt $(PKGS)

revive: ensure-revive
	@revive --config $(DIR)assets/lint/config.toml -exclude $(DIR)vendor/... -formatter friendly $(DIR)...

clean: ensure-pkger
	@rm -f debug $(BINARY_NAME)
	@rm -rf dist
	@pkger clean

clean-vendor:
	rm -rf vendor/

test:
	@go test -coverprofile=./cov-raw.out -v $(TEST_PKGS)
	@cat ./cov-raw.out | grep -v "generated" > ./cov.out

cli-cover-report:
	@go tool cover -func=cov.out

html-cover-report:
	@go tool cover -html=cov.out -o .coverage.html

deps:
	@go build -v ./...

compile: deps ensure-pkger
	@pkger
	@$(GOARGS) go build $(GO_BUILD_ARGS) -o $(DIR)/$(BINARY_NAME) $(BUILD_PATH)

run:
	@go run $(BUILD_PATH)

debug: ensure-delve
	@dlv debug \
		--headless \
		--listen=127.10.10.2:$(DEBUG_PORT) \
		--api-version=2 $(BUILD_PATH) \
		--build-flags="-tags debug" \
		-- serve --config ./examples/goveal.yaml $(DIR)/examples/slides.md

download-reveal:
	@mkdir -p $(DIR)/assets/reveal
	@curl -sL https://github.com/hakimel/reveal.js/archive/$(REVEALJS_VERSION).tar.gz | tar -xvz --strip-components=1 -C $(DIR)/assets/reveal --wildcards "*.js" --wildcards "*.css" --wildcards "*.html" --wildcards "*.woff" --wildcards "*.ttf" --exclude "test" --exclude "gruntfile.js" --exclude "examples/*.html"
	@mkdir -p $(DIR)/assets/reveal/plugin/menu $(DIR)/assets/reveal/plugin/mouse-pointer
	@git clone https://github.com/denehyg/reveal.js-menu.git $(DIR)/assets/reveal/plugin/menu
	@curl -L -o $(DIR)/assets/reveal/plugin/mouse-pointer/mouse-pointer.js https://raw.githubusercontent.com/caiofcm/plugin-revealjs-mouse-pointer/master/mouse-pointer.js

watch: ensure-reflex
	@reflex -r '\.go$$' -s -- sh -c 'make debug'

watch-test: ensure-reflex
	@reflex -r '_test\.go$$' -s -- sh -c 'make test'

cloc:
	@cloc --vcs=git --exclude-dir=.idea,.vscode,.theia,public,docs, .

serve-godoc: ensure-godoc
	@godoc -http=:6060

serve-docs: ensure-reflex docs
	@reflex -r '\.md$$' -s -- sh -c 'mdbook serve -d $(DIR)/public -n 127.0.0.1 $(DIR)/docs'

docs:
	@mdbook build -d $(DIR)/public $(DIR)/docs`

test-release: ensure-goreleaser ensure-packr2
	@goreleaser --snapshot --skip-publish --rm-dist

ensure-revive:
ifeq (, $(shell which revive))
	$(shell go get -u github.com/mgechev/revive)
endif

ensure-delve:
ifeq (, $(shell which dlv))
	$(shell go get -u github.com/go-delve/delve/cmd/dlv)
endif

ensure-reflex:
ifeq (, $(shell which reflex))
	$(shell go get -u github.com/cespare/reflex)
endif

ensure-godoc:
ifeq (, $(shell which godoc))
	$(shell go get -u golang.org/x/tools/cmd/godoc)
endif

ensure-pkger:
ifeq (, $(shell which pkger))
	$(shell go get -u github.com/markbates/pkger/cmd/pkger)
endif

ensure-goreleaser:
ifeq (, $(shell which goreleaser))
	$(shell curl -sL https://github.com/goreleaser/goreleaser/releases/download/v$(GORELEASER_VERSION)/goreleaser_Linux_x86_64.tar.gz | tar -xvz --exclude "*.md" -C $$GOPATH/bin)
endif