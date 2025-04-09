export PATH := $(shell go env GOPATH)/bin:$(PATH)

TEST ?= $(shell go list ./... | grep -v 'vendor')
GOFMT_FILES ?= $(shell find . -name '*.go' | grep -v vendor)
PKG_NAME = vmc

default: build

build: fmtcheck
	go install

init:
	go build -o terraform-provider-vmc
	terraform init

debug:
	go build -gcflags="all=-N -l"

test: fmtcheck
	go test $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 240m

debugacc: fmtcheck
	TF_ACC=1 dlv test -o /dev/null $(TEST) -- -test.v $(TESTARGS)

vet:
	@echo "go vet ."
	@go vet $(shell go list ./... | grep -v vendor/) || { echo ""; echo "Vet found suspicious constructs. Please fix them before submitting your code."; exit 1; }

fmt:
	gofmt -w -s $(GOFMT_FILES)

fmtcheck:
	@$(CURDIR)/scripts/gofmtcheck.sh

test-compile:
	@if [ "$(TEST)" = "./..." ]; then echo "ERROR: Set TEST to a specific package. For example,"; echo "  make test-compile TEST=./$(PKG_NAME)"; exit 1; fi
	go test -c $(TEST) $(TESTARGS)

tools:
	GO111MODULE=on go install -mod=mod github.com/katbyte/terrafmt

docs-hcl-lint:
	@echo "==> Checking HCL formatting..."
	@terrafmt diff ./docs --check --pattern '*.md' --quiet || (echo; echo "Unexpected HCL differences. Run 'make docs-hcl-fix'."; exit 1)

docs-hcl-fix:
	@echo "==> Applying HCL formatting..."
	@terrafmt fmt ./docs --pattern '*.md'

.PHONY: build init test testacc debugacc fmt fmtcheck vet tools test-compile docs-hcl-lint docs-hcl-fix
