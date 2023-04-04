variable "cluster_name" {}

variable "vpc_id" {}

variable "keypath" {
  description = "Path to the PEM used as an ssh key to each host."
  default = "./ssh_keys/privatekey.pem"
}

variable "extra_tags" {
  type        = map(string)
  default     = {}
  description = "A map of arbitrary, customizable string key/value pairs to be included alongside a preset map of tags to be used across myriad AWS resources."
}