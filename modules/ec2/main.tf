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

# resource "aws_instance" "main" {
#   ami                    = data.aws_ami.image.id
#   instance_type = var.instance_type
#   subnet_id              = aws_default_subnet.a.id
#   key_name               = var.key_name
#   vpc_security_group_ids = [aws_security_group.main.id]
#   iam_instance_profile   = aws_iam_instance_profile.ec2_profile.name
#   tags = {
#     Name = var.name
#   }
# }


resource "aws_instance" "main" {
  ami                    = data.aws_ami.image.id
  instance_type          = var.instance_type
  subnet_id              = "subnet-02bd6f23bd2e48675" # ipv6
  key_name               = var.key_name
  vpc_security_group_ids = [aws_security_group.main.id]
  ipv6_address_count = 1
  # ipv6_addresses         = ["2600:1f18:1248:e300:813:9e07:6f2e:6f7a"]
  iam_instance_profile   = var.name
  tags = {
    Name = "sock_test"
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

# resource "aws_eip" "main" {
#   instance = aws_instance.main.id
# }

data "aws_vpc" "default" {
  default = true
}
# aws_default_vpc.default.cidr_block

# NOTE: different subnets have different spot prices
# resource "aws_default_subnet" "a" {
#   availability_zone = "us-east-1a"
# }

resource "aws_security_group" "main" {
  name        = var.name
  vpc_id = data.aws_vpc.default.id
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["${var.ssh_ip}/32"]
  }
  ingress {
    # from_port   = 0
    # to_port     = 0
    # protocol    = "-1"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    ipv6_cidr_blocks = ["::/0"]
    cidr_blocks = ["0.0.0.0/0"]
  }
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    # cidr_blocks = ["0.0.0.0/0"]
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
