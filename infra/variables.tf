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