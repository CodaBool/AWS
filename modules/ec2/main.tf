resource "aws_spot_instance_request" "main" {
  ami                  = data.aws_ami.image.id
  # ami                  = "ami-03a45a5ac837f33b7" # ami-084237e82d7842286
  spot_price           = "0.002400" # max price to request, use aws ec2 describe-spot-price-history
  wait_for_fulfillment = true     # wait up to 10min 
  instance_type        = "t4g.nano"
  subnet_id            = aws_default_subnet.a.id
  # user_data            = "sudo yum update" # doesnt seem like these were applied?
  # user_data_replace_on_change = true
  key_name               = var.key_name
  vpc_security_group_ids = [aws_security_group.main.id]
}

data "aws_ami" "image" {
  most_recent = true
  owners = ["self"]

  filter {
    name   = "name"
    values = [var.ami_name]
  }
}
# resource "aws_key_pair" "main" {
#   key_name   = var.name
#   public_key = tls_private_key.rsa.public_key_openssh
# }

# resource "tls_private_key" "rsa" {
#   algorithm = "RSA"
#   rsa_bits  = 4096
# }

# resource "local_file" "key" {
#   # gets placed on the root of where its ran
#   content  = tls_private_key.rsa.private_key_pem
#   file_permission = "600"
#   filename = "${var.name}.pem"
# }

data "external" "my_ip" {
  program = ["curl", "https://ipinfo.io"]
}

resource "aws_default_vpc" "default" {}
# aws_default_vpc.default.cidr_block

# pick the cheapest subnet
resource "aws_default_subnet" "a" {
  availability_zone = "us-east-1a"
}
# resource "aws_default_subnet" "b" {
#   availability_zone = "us-east-1b"
# }
# resource "aws_default_subnet" "c" {
#   availability_zone = "us-east-1c"
# }
# resource "aws_default_subnet" "d" {
#   availability_zone = "us-east-1d"
# }
# resource "aws_default_subnet" "f" {
#   availability_zone = "us-east-1f"
# }

resource "aws_security_group" "main" {
  name   = var.name
  vpc_id = aws_default_vpc.default.id
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["${data.external.my_ip.result.ip}/32"]
  }
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}