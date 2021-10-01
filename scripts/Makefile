GO                 ?= go
BINDIR             ?= $(CURDIR)/bin

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
.golangci-bin:
	@echo "===> Installing Golangci-lint <==="
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $@ v1.41.1

.PHONY: golangci
golangci: .golangci-bin
	@echo "===> Running golangci <==="
	@GOOS=linux .golangci-bin/golangci-lint run -c .golangci.yml

.PHONY: golangci-fix
golangci-fix: .golangci-bin
	@echo "===> Running golangci-fix <==="
	@GOOS=linux .golangci-bin/golangci-lint run -c .golangci.yml --fix

.PHONY: clean
clean:
	@rm -rf $(BINDIR)
	@rm -rf .golangci-bin
