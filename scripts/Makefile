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

.PHONY: clean
clean:
	@rm -rf $(BINDIR)
