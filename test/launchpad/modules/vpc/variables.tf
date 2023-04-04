variable "cluster_name" {}

variable "host_cidr" {
  description = "CIDR IPv4 range to assign to EC2 nodes"
  default     = "172.31.0.0/16"
}

variable "extra_tags" {
  type        = map(string)
  default     = {}
  description = "A map of arbitrary, customizable string key/value pairs to be included alongside a preset map of tags to be used across myriad AWS resources."
}
