resource "digitalocean_droplet" "summary-client" {
  image = "docker-16-04"
  name = "summary-client"
  region = "sfo2"
  size = "s-1vcpu-1gb"
  private_networking = false
  ssh_keys = [
    "${var.ssh_fingerprint}"
  ]

  connection {
    user = "root"
    type = "ssh"
    private_key = "${file(var.pvt_key)}"
    timeout = "2m"
  }
}