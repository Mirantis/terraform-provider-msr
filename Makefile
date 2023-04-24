PROVIDER=msr
TERRAFORM_PROVIDER_NAMESPACE=registry.terraform.io/mirantis
BINARY_ROOT=terraform-provider
INSTALL_ROOT?=$(HOME)/.terraform.d/plugins
LOCAL_BIN_PATH?=./bin
TEST_TF_CHART_ROOT?=${CURDIR}/test/launchpad
TF_LOCK_FILE?=${TEST_TF_CHART_ROOT}/.terraform.lock.hcl

TAG=$(shell git describe --tags)
VERSION?=$(TAG:v%=%)

ARCHES?=amd64 arm64
OSES?=linux darwin

GO=$(shell which go)

default: install

.PHONY: version
version:
	@echo $(VERSION)

.PHONY: clean
clean:
	rm -rf "$(LOCAL_BIN_PATH)"
	rm -rf "$(INSTALL_ROOT)/$(TERRAFORM_PROVIDER_NAMESPACE)/$(PROVIDER)"

.PHONY: build
build:
	mkdir -p $(LOCAL_BIN_PATH)
	for OS in $(OSES); do \
		for ARCH in $(ARCHES); do \
			GOOS=$${OS} GOARCH=$${ARCH} $(GO) build -v -o "$(LOCAL_BIN_PATH)/$(BINARY_ROOT)-$(PROVIDER)-$${OS}_$${ARCH}" "./cmd/$(PROVIDER)"; \
		done; \
	done;

.PHONY: release
release:
	goreleaser release --rm-dist --snapshot --skip-publish  --skip-sign

.PHONY: install
install: build
	for OS in $(OSES); do \
		for ARCH in $(ARCHES); do \
			mkdir -p "$(INSTALL_ROOT)/$(TERRAFORM_PROVIDER_NAMESPACE)/$(PROVIDER)/$(VERSION)/$${OS}_$${ARCH}"; \
			cp "$(LOCAL_BIN_PATH)/$(BINARY_ROOT)-$(PROVIDER)-$${OS}_$${ARCH}" "$(INSTALL_ROOT)/$(TERRAFORM_PROVIDER_NAMESPACE)/$(PROVIDER)/$(VERSION)/$${OS}_$${ARCH}/$(BINARY_ROOT)-$(PROVIDER)_v$(VERSION)"; \
		done; \
	done;

.PHONY: test-unit
test-unit:
	go test -v -cover ./...

.PHONY: test-acceptance
test-acceptance: clean build install test-unit
	rm -f ${TF_LOCK_FILE}
	terraform -chdir=${TEST_TF_CHART_ROOT} init --upgrade
	terraform -chdir=${TEST_TF_CHART_ROOT} apply -auto-approve
	terraform -chdir=${TEST_TF_CHART_ROOT} destroy -auto-approve
