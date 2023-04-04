variable "cluster_name" {
  default = "tf-mcc-provider-test"
}

variable "aws_region" {
  default = "us-west-2"
}

variable "vpc_cidr" {
  default = "172.31.0.0/16"
}

variable "admin_username" {
  default = "admin"
}
variable "admin_password" {
  default = "tum40PJ9lGIFySvPeNQYsUBGz8zIKlia"
}

variable "keypath" {
  description = "Path to the PEM used as an ssh key to each host."
  type = string
  default = "./ssh_keys/privatekey.pem"
}


variable "manager_count" {
  default = 1
}

variable "worker_count" {
  default = 3
}

variable "windows_worker_count" {
  default = 0
}

variable "msr_count" {
  default = 1
}

variable "manager_type" {
  default = "m5.large"
}

variable "worker_type" {
  default = "m5.large"
}

variable "msr_type" {
  default = "m5.large"
}
variable "manager_volume_size" {
  default = 100
}

variable "worker_volume_size" {
  default = 100
}

variable "msr_volume_size" {
  default = 100
}


variable "mcr_version" {
  type        = string
  default     = "20.10.7"
  description = "The mcr version to deploy across all nodes in the cluster."
}

variable "mcr_channel" {
  type        = string
  default     = "stable"
  description = "The channel to pull the mcr installer from."
}

variable "mcr_repo_url" {
  type        = string
  default     = "https://s3.amazonaws.com/repos-internal.mirantis.com"
  description = "The repository to source the mcr installer."
}

variable "mcr_install_url_linux" {
  type        = string
  default     = "https://get.mirantis.com/"
  description = "Location of Linux installer script."
}

variable "mcr_install_url_windows" {
  type        = string
  default     = "https://get.mirantis.com/install.ps1"
  description = "Location of Windows installer script."
}

variable "mke_version" {
  type        = string
  default     = "3.5.2"
  description = "The UCP version to deploy."
}

variable "mke_image_org" {
  type        = string
  default     = "docker.io/mirantis"
  description = "The repository to pull the UCP images from."
}

variable "mke_install_flags" {
  type        = list(string)
  default     = []
  description = "The UCP installer flags to use."
}

variable "mke_default_orchestrator" {
  type        = string
  default     = "swarm"
  description = "Set the MKE default orchestrator."
}

variable "nodeport_range" {
  type        = string
  default     = "32768-35535"
  description = "Kubernetes nodeport range."
}

variable "msr_version" {
  type        = string
  default     = "2.9.0"
  description = "The DTR version to deploy."
}

variable "msr_image_org" {
  type        = string
  default     = "docker.io/mirantis"
  description = "The repository to pull the DTR images from."
}

variable "msr_install_flags" {
  type        = list(string)
  default     = ["--ucp-insecure-tls"]
  description = "The DTR installer flags to use."
}

variable "msr_replica_config" {
  type        = string
  default     = "sequential"
  description = "Set to 'sequential' to generate sequential replica id's for cluster members, for example 000000000001, 000000000002, etc. ('random' otherwise)"
}


variable "windows_administrator_password" {
  default = "w!ndozePassw0rd"
}

variable "extra_tags" {
  type        = map(string)
  default     = {}
  description = "A map of arbitrary, customizable string key/value pairs to be included alongside a preset map of tags to be used across myriad AWS resources."
}
