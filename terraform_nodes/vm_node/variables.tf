variable "name" {
  type = string
}

variable "machine_type" {
  type    = string
  default = "e2-medium"
}

variable "source_image" {
  type    = string
  default = "ubuntu-2204-jammy-v20231030"
}

variable "disk_type" {
  type    = string
  default = "pd-balanced"
}

variable "pat_token" {
  type      = string
  sensitive = true
}

variable "github_name" {
  type = string
}

variable "repo_name" {
  type = string
}