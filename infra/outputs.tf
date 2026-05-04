output "server_ip" {
  value       = twc_floating_ip.main.ip
  description = "Public IPv4"
}

output "ssh_command" {
  value       = "ssh root@${twc_floating_ip.main.ip}"
  description = "SSH command"
}