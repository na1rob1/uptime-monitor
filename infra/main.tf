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

resource "twc_ssh_key" "main" {
  name = "uptime-monitor"
  body = file(var.ssh_public_key_path)
}

data "twc_os" "ubuntu" {
  name    = "ubuntu"
  version = "22.04"
}

data "twc_configurator" "ru" {
  location = "ru-1"
}

resource "twc_server" "uptime" {
  name              = "uptime-monitor"
  os_id             = data.twc_os.ubuntu.id
  availability_zone = "spb-3"

  configuration {
    configurator_id = data.twc_configurator.ru.id
    disk            = 15360
    cpu             = 1
    ram             = 1024
  }

  ssh_keys_ids = [twc_ssh_key.main.id]
}

resource "twc_floating_ip" "main" {
  availability_zone = "spb-3"
  comment           = "uptime-monitor IPv4"

  resource {
    type = "server"
    id   = twc_server.uptime.id
  }
}