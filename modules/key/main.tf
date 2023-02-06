resource "aws_key_pair" "main" {
  key_name   = var.key_name
  public_key = tls_private_key.rsa.public_key_openssh
}

resource "tls_private_key" "rsa" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "local_file" "key" {
  # gets placed on the root of where its ran
  content  = tls_private_key.rsa.private_key_pem
  file_permission = "600"
  filename = "${var.key_name}.pem"
}

variable "key_name" {
  type    = string
}