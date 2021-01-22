GO = $(shell which go)
GOTEST = $(shell which gotest)
GOLANGCI_LINT = $(shell which golangci-lint)
GO_PKGS = $(shell $(GO) list ./... | grep -v /vendor/)
TARGET_DIR = $(shell pwd)/.target

ifeq ($(GOTEST),)
GOTEST = "$(GO) test"
endif

.PHONY: ensure-target
ensure-target:
	@mkdir -p $(TARGET_DIR)/coverage

.PHONY: check-pre
check-pre: ensure-target
	@if [ -z "$(GO)" ]; then echo "Missing go command"; exit 1; fi

.PHONY: clean
clean:
	@rm -Rf $(TARGET_DIR)

.PHONY: fmt
fmt: check-pre
	@$(GO) fmt $(GO_PKGS)

.PHONY: bench
bench:
	@$(GO) test \
		-run=notest \
		-bench=. \
		-v ./...

.PHONY: test
test: check-pre
	@$(GOTEST) \
		-coverprofile=$(TARGET_DIR)/coverage/report.out \
		-count=1 \
		-failfast \
		-race \
		-v ./...

.PHONY: lint
lint: check-pre
	@if [ -z "$(GOLANGCI_LINT)" ]; then echo "Missing golangci_lint command"; exit 1; fi
	@$(GOLANGCI_LINT) run --exclude-use-default=false ./...

.PHONY: cover
cover: test
	@$(GO) tool cover -func=$(TARGET_DIR)/coverage/report.out
	@$(GO) tool cover \
		-html=$(TARGET_DIR)/coverage/report.out \
		-o $(TARGET_DIR)/coverage/report.html
