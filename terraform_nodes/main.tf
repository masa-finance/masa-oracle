
locals {
  config = yamldecode(file("./config.yaml"))

  joining_nodes = [for key, node in local.config : key if node.bootnodes != null]
  default_nodes = [for key, node in local.config : key if node.bootnodes == null]
  master_node   = try(element(local.default_nodes, 0), null)
  worker_nodes = try(concat(
    slice(local.default_nodes, 0, index(local.default_nodes, local.master_node)),
    slice(local.default_nodes, index(local.default_nodes, local.master_node) + 1, length(local.default_nodes))
  ), null)
  nodes_with_bootnodes = { for key, node in local.config : key => node.bootnodes if node.bootnodes != null }
}

module "testing" {
  for_each    = local.config
  source      = "./vm_node"
  name        = each.key
  pat_token   = var.pat_token
  github_name = "masa-finance"
  repo_name   = "masa-oracle"
}

resource "local_file" "test" {
  content = templatefile("${path.module}/pipeline.tftpl", {
    nodes                = local.config,
    nodes_with_bootnodes = local.nodes_with_bootnodes
    joining_nodes        = local.joining_nodes,

    # deddicated nodes
    master_node  = local.master_node,
    worker_nodes = local.worker_nodes,
  })
  filename = "../.github/workflows/masa-oracle.yaml"
}

data "google_compute_network" "vpc_network" {
  name = "default"
}

resource "google_compute_firewall" "allow_ssh" {
  name    = "masa-node-ssh"
  network = data.google_compute_network.vpc_network.name
  target_tags = [
    "masa", "pipe"
  ] // this targets our tagged VM
  source_ranges = ["0.0.0.0/0"]

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }
}
