
// Mirantis installing terraform provider
provider "mirantis-msr" {}

// Launchpad installer
resource "mirantis-msr_launchpad" "test" {

  skip_destroy = true

  metadata {
    name = var.cluster_name
  }
  spec {
    cluster {
      prune = true
    }

    dynamic "host" {
      for_each = module.managers.machines
      content {
        role = host.value.tags["Role"]
        ssh {
          address  = host.value.public_ip
          user     = "ubuntu"
          key_path = var.keypath
        }
      }
    }

    dynamic "host" {
      for_each = module.workers.machines
      content {
        role = host.value.tags["Role"]
        ssh {
          address  = host.value.public_ip
          user     = "ubuntu"
          key_path = var.keypath
        }
      }
    }

    dynamic "host" {
      for_each = module.msrs.machines
      content {
        role = host.value.tags["Role"]
        ssh {
          address  = host.value.public_ip
          user     = "ubuntu"
          key_path = var.keypath
        }
      }
    }

    mcr {
      channel             = "stable"
      install_url_linux   = var.mcr_install_url_linux
      install_url_windows = var.mcr_install_url_windows
      repo_url            = var.mcr_repo_url
      version             = var.mcr_version
    } // mcr

    mke {
      admin_password = var.admin_password
      admin_username = var.admin_username
      image_repo     = var.mke_image_org
      version        = var.mke_version
      install_flags = [
        "--san=${module.managers.lb_dns_name}",
        "--default-node-orchestrator=${var.mke_default_orchestrator}",
        "--nodeport-range=${var.nodeport_range}",
        "--cloud-provider=aws"
      ]
      upgrade_flags = [
        "--force-recent-backup",
        "--force-minimums"
      ]
    } // mke

    msr {
      image_repo  = var.msr_image_org
      version     = var.msr_version
      replica_ids = "sequential"
      install_flags = ["--ucp-insecure-tls",
      "--dtr-external-url ${module.msrs.lb_dns_name}"]
    } // msr
  }   // spec
}

output "mke_cluster_name" {
  value = mirantis-msr_launchpad.test.metadata[0].name
}
