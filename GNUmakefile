default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

# Local install of the plugin
.PHONY: local
local:
	goreleaser build --clean --single-target --skip-validate
