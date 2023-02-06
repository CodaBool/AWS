resource "aws_spot_instance_request" "main" {
  ami                    = data.aws_ami.image.id
  spot_price             = var.price    # max price to request, use aws ec2 describe-spot-price-history
  wait_for_fulfillment   = true          # wait up to 10min 
  instance_type          = var.instance_type
  subnet_id              = aws_default_subnet.a.id
  key_name               = var.key_name
  vpc_security_group_ids = [aws_security_group.main.id]
  iam_instance_profile   = var.profile_name
}

data "aws_ami" "image" {
  most_recent = true
  owners = ["self"]
  filter {
    name = "tag:Name"
    values = [var.name]
  }
}

# can use this to get ip with data.external.my_ip.result.ip
# however, if running terraform in a pipeline this is moot
# data "external" "my_ip" {
#   program = ["curl", "https://ipinfo.io"]
# }

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
    cidr_blocks = ["${var.ssh_ip}/32"]
  }
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}