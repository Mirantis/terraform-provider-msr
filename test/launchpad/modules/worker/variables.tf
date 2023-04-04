variable "cluster_name" {}

variable "vpc_id" {}

variable "instance_profile_name" {}

variable "security_group_id" {}

variable "subnet_ids" {
  type = list(string)
}

variable "image_id" {}

variable "kube_cluster_tag" {}

variable "ssh_key" {
  description = "SSH key name"
}

variable "worker_count" {
  default = 3
}

variable "worker_type" {
  default = "m5.large"
}

variable "worker_volume_size" {
  default = 100
}

variable "extra_tags" {
  type        = map(string)
  default     = {}
  description = "A map of arbitrary, customizable string key/value pairs to be included alongside a preset map of tags to be used across myriad AWS resources."
}
