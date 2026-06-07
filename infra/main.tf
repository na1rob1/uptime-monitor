terraform {
  required_providers {
    twc = {
      source = "tf.timeweb.cloud/timeweb-cloud/timeweb-cloud"
    }
  }
  required_version = ">= 1.0"
}

provider "twc" {
  token = var.twc_token
}

# LAB 2
#resource "twc_ssh_key" "main" {
#  name = "uptime-monitor"
#  body = file(var.ssh_public_key_path)
#}
#
#data "twc_os" "ubuntu" {
#  name    = "ubuntu"
#  version = "22.04"
#}
#
#data "twc_configurator" "ru" {
#  location = "ru-1"
#}
#
#resource "twc_server" "uptime" {
#  name              = "uptime-monitor"
#  os_id             = data.twc_os.ubuntu.id
#  availability_zone = "spb-3"
#
#  configuration {
#    configurator_id = data.twc_configurator.ru.id
#    disk            = 15360
#    cpu             = 1
#    ram             = 1024
#  }
#
#  ssh_keys_ids = [twc_ssh_key.main.id]
#}
#
#resource "twc_floating_ip" "main" {
#  availability_zone = "spb-3"
#  comment           = "uptime-monitor IPv4"
#
#  resource {
#    type = "server"
#    id   = twc_server.uptime.id
#  }
#}

# LAB 3

data "twc_k8s_preset" "master" {
  cpu  = 2
  type = "master"
}

data "twc_k8s_preset" "worker" {
  cpu  = 2
  type = "worker"
}

resource "twc_k8s_cluster" "uptime" {
  name           = "uptime"
  description    = "Lab 3 cluster"
  version        = var.k8s_version
  network_driver = "flannel"
  ingress        = true
  preset_id      = data.twc_k8s_preset.master.id
}

resource "twc_k8s_node_group" "workers" {
  cluster_id     = twc_k8s_cluster.uptime.id
  name           = "workers"
  preset_id      = data.twc_k8s_preset.worker.id
  node_count     = 1
  is_autoscaling = true
  min_size       = 1
  max_size       = 3
}

data "twc_database_preset" "postgres" {
  location = "ru-1"
  type     = "postgres"
  disk     = 8192

  price_filter {
    from = 100
    to   = 500
  }
}

resource "twc_database_cluster" "main" {
  name           = "uptime-cluster"
  type           = "postgres17"
  preset_id      = data.twc_database_preset.postgres.id
  is_external_ip = true
}

resource "twc_database_instance" "main" {
  cluster_id = twc_database_cluster.main.id
  name       = "uptime"
}

resource "twc_database_user" "main" {
  cluster_id = twc_database_cluster.main.id
  login      = "uptime"
  password   = var.db_password

  instance {
    instance_id = twc_database_instance.main.id
    privileges  = ["SELECT", "INSERT", "UPDATE", "DELETE", "CREATE"]
  }
}

