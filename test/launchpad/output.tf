
output "mke_lb_dns" {
  value = module.managers.lb_dns_name
}

output "msr_lb_dns" {
  value = module.msrs.lb_dns_name
}