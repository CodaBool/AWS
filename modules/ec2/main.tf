# spot instances will be evicted once price changes
# they should be close in price to an on-demand instance

# resource "aws_spot_instance_request" "main" {
#   ami                    = data.aws_ami.image.id
#   spot_price             = data.external.lowest_price.result.price
#   wait_for_fulfillment   = true          # wait up to 10min 
#   instance_type          = var.instance_type
#   subnet_id              = aws_default_subnet.a.id
#   key_name               = var.key_name
#   vpc_security_group_ids = [aws_security_group.main.id]
#   iam_instance_profile   = aws_iam_instance_profile.ec2_profile.name
# }

resource "aws_instance" "main" {
  ami                    = data.aws_ami.image.id
  instance_type          = var.instance_type
  subnet_id              = var.subnet # ipv6
  key_name               = var.key_name
  vpc_security_group_ids = [aws_security_group.main.id]
  ipv6_addresses         = [var.ip]
  iam_instance_profile   = var.name
  tags = {
    Name = var.name
  }
}

# max price to request, use aws ec2 describe-spot-price-history
# data "external" "lowest_price" {
#   program = ["bash", "${path.module}/price.sh", var.instance_type]
# }

data "aws_ami" "image" {
  most_recent = true
  owners = ["self"]
  filter {
    name = "tag:Name"
    values = ["${var.name}*"]
  }
}

data "aws_vpc" "default" {
  default = true
}

resource "aws_security_group" "main" {
  name        = var.name
  vpc_id = data.aws_vpc.default.id
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    ipv6_cidr_blocks = [var.ssh_ip] # must be ipv6 ending in /128
  }
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    ipv6_cidr_blocks = ["::/0"]
  }
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    ipv6_cidr_blocks = ["::/0"]
  }
  tags = {
    Name = var.name
  }
}


resource "aws_iam_role" "cw_assume" {
  name = var.name
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "ec2.amazonaws.com"
      }
    }]
  })
}

# Cloudwatch resources
resource "aws_iam_instance_profile" "ec2_profile" {
  name = var.name
  role = aws_iam_role.cw_assume.name
}

resource "aws_iam_role_policy_attachment" "ssm" {
  role       = aws_iam_role.cw_assume.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

resource "aws_iam_role_policy_attachment" "logs" {
  role       = aws_iam_role.cw_assume.name
  policy_arn = "arn:aws:iam::aws:policy/CloudWatchAgentServerPolicy"
}

resource "aws_iam_role_policy_attachment" "retention" {
  role       = aws_iam_role.cw_assume.name
  policy_arn = aws_iam_policy.retention.arn
}

# could rm ssm, ssm-agent requires ipv4
resource "aws_iam_policy" "retention" {
  name_prefix        = "change_retention"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "logs:PutRetentionPolicy"
      Effect   = "Allow"
      Resource = "*"
    },
    {
      Action = "ssm:*"
      Effect   = "Allow"
      Resource = "*"
    }]
  })
}
