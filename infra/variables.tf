variable "twc_token" {
  type        = string
  sensitive   = true
  description = "Timeweb Cloud API token"
}

variable "ssh_public_key_path" {
  type        = string
  description = "Public SSH key file"
  default     = "~/.ssh/timeweb_uptime.pub"
}

variable "k8s_version" {
  type        = string
  description = "Kubernetes version"
  # поменял версию с v1.31.14+k0s.0 08.09.26
  default     = "v1.32.13+k0s.0"
}
#
#variable "db_password" {
#  type        = string
#  sensitive   = true
#  description = "Postgres user password"
#}