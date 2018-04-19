resource "digitalocean_droplet" "summary-server" {
  image = "docker-16-04"
  name = "summary-server"
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

//  provisioner "remote-exec" {
//    inline = [
//      "sudo ufw allow 80",
//      "sudo ufw allow 443",
//      "sudo apt update && appt install -y letsencrypt",
//      "sudo letsencrypt certonly --standalone -n -agree-tos --email aaronluannguyen@gmail.com -d api.aaronnluannguyen.me"
//    ]
//  }
//}
//
//resource "digitalocean_record" "api" {
//  domain = "api.aaronnluannguyen.me"
//  type = "A"
//  name = "api"
//  value = "${digitalocean_droplet.summary-server.ipv4_address}"
}