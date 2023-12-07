variable "project_id" {
  default = "masa-chain"
}

variable "region" {
  default = "us-central1"
}

variable "zone" {
  default = "us-central1-a"
}

variable "pat_token" {
  type      = string
  sensitive = true
}