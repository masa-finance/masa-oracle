terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "5.7.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "3.5.1"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
  zone    = var.zone
}