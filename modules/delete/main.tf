resource "aws_iam_role" "lifecycle" {
  name               = "lifecycle"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "dlm.amazonaws.com"
      }
    }]
  })
}

resource "aws_iam_role_policy_attachment" "lifecycle" {
  role       = aws_iam_role.lifecycle.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSDataLifecycleManagerServiceRole"
}

# TODO: find if these lifecycle policy were killing the instance
# resource "aws_dlm_lifecycle_policy" "instance" {
#   description        = "slap delete instance"
#   execution_role_arn = aws_iam_role.lifecycle.arn
#   policy_details {
#     resource_types = ["INSTANCE"]
#     schedule {
#       name = "every day"
#       create_rule {
#         cron_expression = "cron(0 12 * * ? *)"
#       }
#       retain_rule {
#         count = 1
#       }
#     }
#     target_tags = {
#       Name = "slap"
#     }
#   }
# }

# resource "aws_dlm_lifecycle_policy" "volume" {
#   description        = "slap delete volume"
#   execution_role_arn = aws_iam_role.lifecycle.arn
#   policy_details {
#     resource_types = ["VOLUME"]
#     schedule {
#       name = "every day"
#       create_rule {
#         cron_expression = "cron(0 12 * * ? *)"
#       }
#       retain_rule {
#         count = 1
#       }
#     }
#     target_tags = {
#       Name = "slap"
#     }
#   }
# }

data "external" "lowest_price" {
  program = ["bash", "${path.module}/price.sh"]
}

data "external" "my_ip" {
  program = ["curl", "https://ipinfo.io"]
}

module "ec2" {
  source        = "../ec2"
  name          = "slap-arm" # this must match tag:Name value
  instance_type = "t4g.nano"
  price         = data.external.lowest_price.result.price
  ssh_ip        = data.external.my_ip.result.ip
  profile_name = module.cloudwatch.profile
}

module "cloudwatch" {
  source = "../cloudwatch"
}