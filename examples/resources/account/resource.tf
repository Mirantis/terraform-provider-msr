data "msr_account" "example" {
  name_or_id = "example"
}

data "msr_accounts" "example" {
  filter = "all"
}
