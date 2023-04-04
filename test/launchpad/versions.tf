
terraform {
  required_version = ">= 1.0.0"
  required_providers {
    mirantis-msr-connect = {
      version = ">= 0.9.0"
      source  = "mirantis.com/providers/mirantis-msr-connect"
    }
    docker = {
      source  = "kreuzwerker/docker"
      version = "2.16.0"
    }
  }
}
