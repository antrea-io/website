GO                 ?= go
BINDIR             ?= $(CURDIR)/bin
GOLANGCI_LINT_VERSION := v1.54.0
GOLANGCI_LINT_BINDIR  := .golangci-bin
GOLANGCI_LINT_BIN     := $(GOLANGCI_LINT_BINDIR)/$(GOLANGCI_LINT_VERSION)/golangci-lint

.PHONY: all
all: update-main-docs freeze-version-docs

.PHONY: update-main-docs
update-main-docs:
	@mkdir -p $(BINDIR)
	$(GO) build -o $(BINDIR) antrea.io/website/scripts/update-main-docs

.PHONY: freeze-version-docs
freeze-version-docs:
	@mkdir -p $(BINDIR)
	$(GO) build -o $(BINDIR) antrea.io/website/scripts/freeze-version-docs

# code linting
$(GOLANGCI_LINT_BIN):
	@echo "===> Installing Golangci-lint <==="
	@rm -rf $(GOLANGCI_LINT_BINDIR)/* # delete old versions
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOLANGCI_LINT_BINDIR)/$(GOLANGCI_LINT_VERSION) $(GOLANGCI_LINT_VERSION)

.PHONY: golangci
golangci: $(GOLANGCI_LINT_BIN)
	@echo "===> Running golangci <==="
	@GOOS=linux $(GOLANGCI_LINT_BIN) run -c .golangci.yml

.PHONY: golangci-fix
golangci-fix: $(GOLANGCI_LINT_BIN)
	@echo "===> Running golangci-fix <==="
	@GOOS=linux $(GOLANGCI_LINT_BIN) run -c .golangci.yml --fix

.PHONY: clean
clean:
	@rm -rf $(BINDIR)
	@rm -rf $(GOLANGCI_LINT_BINDIR)

.PHONY: test-unit
test-unit:
	@echo "===> Running unit tests <==="
	$(GO) test -race antrea.io/website/scripts/pkg/... 
