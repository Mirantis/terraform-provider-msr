default: test

LOCAL_TAG=$(shell git describe --tags)

# Run local unit testr
.PHONY: test
test:
	go test ./... -timeout 120ms

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -tags=integration -timeout 120m

# Run golangci-lint
.PHONY: lint
lint:
	 docker run -ti --rm -v "$(CURDIR):/data" -w "/data" golangci/golangci-lint:latest golangci-lint run --timeout 55s

# Local install of the plugin
.PHONY: local
local:
	GORELEASER_CURRENT_TAG="$(LOCAL_TAG)" goreleaser build --clean --single-target --skip-validate

	# "Local plugin generated. Use $(CURDIR)/dist/terraform-provider-msr_linux_amd64_v1 as your dev_overrides path in a terraform config file
	#
	# my_tf_config_file:
	# ```
	# provider_installation {
	#
	#	# This disables the version and checksum verifications for this provider
	#	# and forces Terraform to look for the msr provider plugin in the
	#	# given directory.
	#	dev_overrides {
	#		"mirantis/msr" = "$(CURDIR)/dist/terraform-provider-msr_linux_amd64_v1"
	#	}
	#
	#	# For all other providers, install them directly from their origin provider
	#	# registries as normal. If you omit this, Terraform will _only_ use
	#	# the dev_overrides block, and so no other providers will be available.
	#	direct {}
	# }
	# ```
	#
	# then run terraform with a config file override:
	# ```
	#  $/> TF_CLI_CONFIG_FILE=my_tf_config_file terraform plan
	# ```
	# (or use an export)
	#
	# @see: https://developer.hashicorp.com/terraform/cli/config/config-file#development-overrides-for-provider-developers"
	#