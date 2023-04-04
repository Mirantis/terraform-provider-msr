output "lb_dns_name" {
    value = aws_lb.mke_manager.dns_name
}

output "public_ips" {
    value = aws_instance.mke_manager.*.public_ip
}

output "private_ips" {
    value = aws_instance.mke_manager.*.private_ip
}

output "machines" {
  value = aws_instance.mke_manager
}
