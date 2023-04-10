TERRAFORM_PROVIDER_ROOT=mirantis.com/providers
BINARY_ROOT=terraform-provider
INSTALL_ROOT?=$(HOME)/.terraform.d/plugins
LOCAL_BIN_PATH?=./bin
TEST_TF_CHART_ROOT?=${CURDIR}/test/launchpad
TF_LOCK_FILE?=${TEST_TF_CHART_ROOT}/.terraform.lock.hcl

VERSION=0.9.0

PROVIDERS?=mirantis-msr
ARCHES?=amd64 arm64
OSES?=linux darwin

GO=$(shell which go)

default: install

.PHONY: clean
clean:
	rm -rf "$(LOCAL_BIN_PATH)"
	rm -rf "$(INSTALL_ROOT)/$(TERRAFORM_PROVIDER_ROOT)"

.PHONY: build
build:
	mkdir -p $(LOCAL_BIN_PATH)
	for PROVIDER in $(PROVIDERS); do \
		for OS in $(OSES); do \
			for ARCH in $(ARCHES); do \
				GOOS=$$OS GOARCH=$$ARCH $(GO) build -v -o "$(LOCAL_BIN_PATH)/$(BINARY_ROOT)-$$PROVIDER-$$OS_$$ARCH" "./cmd/$$PROVIDER"; \
			done; \
		done; \
	done;

.PHONY: release
release:
	goreleaser release --rm-dist --snapshot --skip-publish  --skip-sign

.PHONY: install
install: build
	for PROVIDER in $(PROVIDERS); do \
		for OS in $(OSES); do \
			for ARCH in $(ARCHES); do \
				mkdir -p "$(INSTALL_ROOT)/$(TERRAFORM_PROVIDER_ROOT)/$$PROVIDER/$(VERSION)/$${OS}_$${ARCH}"; \
				cp "$(LOCAL_BIN_PATH)/$(BINARY_ROOT)-$$PROVIDER-$$OS_$$ARCH" "$(INSTALL_ROOT)/$(TERRAFORM_PROVIDER_ROOT)/$$PROVIDER/$(VERSION)/$${OS}_$${ARCH}/$(BINARY_ROOT)-$$PROVIDER"; \
    	done; \
		done; \
	done;

.PHONY: test-unit
test-unit:
	go test -v -cover ./...

.PHONY: test-integration
test-integration:
	# Running integration tests.
	#
	# You need the following env variables:
	# -> MKE integration tests:
	#  MKE_HOST
	#  MKE_USERNAME
	#  MKE_PASSWORD
	#
	go test -v -cover --tags=integration ./...

.PHONY: test-acc
test-acc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

.PHONY: test-acceptance
test-acceptance: clean build install test-unit
	rm -f ${TF_LOCK_FILE}
	terraform -chdir=${TEST_TF_CHART_ROOT} init --upgrade
	terraform -chdir=${TEST_TF_CHART_ROOT} apply -auto-approve
	#terraform -chdir=${TEST_TF_CHART_ROOT} destroy -auto-approve

.PHONY: tf-destroy
tf-destroy:
	terraform -chdir=${TEST_TF_CHART_ROOT} destroy -auto-approve
