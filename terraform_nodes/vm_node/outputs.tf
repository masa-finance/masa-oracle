output "public_ip" {
  value = google_compute_address.static.address
}

output "instance_name" {
  value = google_compute_instance.default.name
}