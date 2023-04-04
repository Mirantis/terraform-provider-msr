locals {
  global_tags = merge(
    {
      "Name"                                        = var.cluster_name
      "kubernetes.io/cluster/${var.cluster_name}"   = "shared"
    },
    var.extra_tags
  )
}

provider "aws" {
  region = var.aws_region
  # shared_credentials_file = var.aws_shared_credentials_file
  # profile = var.aws_profile
}

module "vpc" {
  source       = "./modules/vpc"
  cluster_name = var.cluster_name
  host_cidr    = var.vpc_cidr
  extra_tags   = local.global_tags
}

module "common" {
  source       = "./modules/common"
  cluster_name = var.cluster_name
  vpc_id       = module.vpc.id
  keypath      = var.keypath
  extra_tags   = local.global_tags
}

module "managers" {
  source                = "./modules/manager"
  manager_count          = var.manager_count
  vpc_id                = module.vpc.id
  cluster_name          = var.cluster_name
  subnet_ids            = module.vpc.public_subnet_ids
  security_group_id     = module.common.security_group_id
  image_id              = module.common.image_id
  kube_cluster_tag      = module.common.kube_cluster_tag
  ssh_key               = var.cluster_name
  instance_profile_name = module.common.instance_profile_name
  extra_tags            = local.global_tags
}

module "msrs" {
  source                = "./modules/msr"
  msr_count             = var.msr_count
  vpc_id                = module.vpc.id
  cluster_name          = var.cluster_name
  subnet_ids            = module.vpc.public_subnet_ids
  security_group_id     = module.common.security_group_id
  image_id              = module.common.image_id
  kube_cluster_tag      = module.common.kube_cluster_tag
  ssh_key               = var.cluster_name
  instance_profile_name = module.common.instance_profile_name
  extra_tags            = local.global_tags
}

module "workers" {
  source                = "./modules/worker"
  worker_count          = var.worker_count
  vpc_id                = module.vpc.id
  cluster_name          = var.cluster_name
  subnet_ids            = module.vpc.public_subnet_ids
  security_group_id     = module.common.security_group_id
  image_id              = module.common.image_id
  kube_cluster_tag      = module.common.kube_cluster_tag
  ssh_key               = var.cluster_name
  instance_profile_name = module.common.instance_profile_name
  worker_type           = var.worker_type
  extra_tags            = local.global_tags
}

module "windows_workers" {
  source                         = "./modules/windows_worker"
  worker_count                   = var.windows_worker_count
  vpc_id                         = module.vpc.id
  cluster_name                   = var.cluster_name
  subnet_ids                     = module.vpc.public_subnet_ids
  security_group_id              = module.common.security_group_id
  image_id                       = module.common.windows_2019_image_id
  kube_cluster_tag               = module.common.kube_cluster_tag
  instance_profile_name          = module.common.instance_profile_name
  worker_type                    = var.worker_type
  windows_administrator_password = var.windows_administrator_password
}


