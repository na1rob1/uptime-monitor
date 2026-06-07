#output "server_ip" {
#  value       = twc_floating_ip.main.ip
#  description = "Public IPv4"
#}
#
#output "ssh_command" {
#  value       = "ssh root@${twc_floating_ip.main.ip}"
#  description = "SSH command"
#}

output "k8s_cluster_id" {
  value       = twc_k8s_cluster.uptime.id
  description = "k8s cluster ID - kubeconfig скачать через panel/CLI"
}

locals {
  postgres_public_host = [
  for n in twc_database_cluster.main.networks : n.ips[0].ip
  if n.type == "public"
  ][0]
}

output "postgres_host" {
  value       = local.postgres_public_host
  description = "Postgres public host"
}

output "postgres_port" {
  value       = twc_database_cluster.main.port
  description = "Postgres port"
}

output "database_url" {
  value       = "postgres://${twc_database_user.main.login}:${var.db_password}@${local.postgres_public_host}:${twc_database_cluster.main.port}/${twc_database_instance.main.name}?sslmode=require"
  sensitive   = true
  description = "DATABASE_URL для backend"
}

output "kubeconfig" {
  value       = twc_k8s_cluster.uptime.kubeconfig
  sensitive   = true
  description = "Kubeconfig — записать в ~/.kube/config-uptime"
}

resource "local_file" "kubeconfig" {
  content         = twc_k8s_cluster.uptime.kubeconfig
  filename        = pathexpand("~/.kube/config-uptime")
  file_permission = "0600"
}