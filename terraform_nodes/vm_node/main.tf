data "google_compute_network" "vpc_network" {
  name = "default"
}

data "google_compute_image" "ubuntu_image" {
  family  = "ubuntu-2204-lts"
  project = "ubuntu-os-cloud"
}

resource "google_compute_address" "static" {
  name = "masa-${var.name}-ipv4"
}

resource "google_compute_instance" "default" {
  name         = var.name
  machine_type = var.machine_type

  boot_disk {
    initialize_params {
      image = data.google_compute_image.ubuntu_image.self_link
      type  = "pd-balanced"
    }
  }

  network_interface {
    network = data.google_compute_network.vpc_network.name

    access_config {
      nat_ip = google_compute_address.static.address
    }
  }

  labels = {
    "app" = "masa-oracle"
  }

  metadata = {
    ssh-keys = "masa:${file("${path.module}/.ssh/masa.pub")}"
  }

  metadata_startup_script = templatefile("${path.module}/startup_script.sh",
    {
      pat_token   = "${var.pat_token}",
      name        = "${var.name}",
      github_name = "${var.github_name}",
      repo_name   = "${var.repo_name}"
  })

  connection {
    type        = "ssh"
    host        = self.network_interface.0.access_config.0.nat_ip
    user        = "masa"
    private_key = file("${path.module}/.ssh/masa")
  }

  provisioner "remote-exec" {
    when = destroy
    inline = [
      "cd /home/masa/actions-runner",
      "bash removal_script.sh",
    ]
  }

  tags = [
    "masa", "pipe"
  ]
}

resource "null_resource" "setup_vm" {
  connection {
    type        = "ssh"
    host        = google_compute_address.static.address
    user        = "masa"
    private_key = file("${path.module}/.ssh/masa")
  }

  provisioner "file" {
    source      = "archive/"
    destination = "/home/masa"
  }

  # Setup service
  provisioner "remote-exec" {
    inline = [
      "mkdir /home/masa/masa-oracle",
      "sudo tar -xzf /home/masa/workspace.tar.gz -C /home/masa/masa-oracle/",
      "sudo cp /home/masa/masa-oracle/system-service/masa-oracle.service /etc/systemd/system/masa-oracle.service",
      "sudo mkdir -p /var/log/masa-oracle",
      "sudo chown -R masa:masa /var/log/masa-oracle",
      "sudo systemctl daemon-reload",
      "sudo systemctl restart masa-oracle.service",
    ]
  }
  lifecycle {
    replace_triggered_by = [google_compute_instance.default]
  }
  depends_on = [google_compute_instance.default]
}
