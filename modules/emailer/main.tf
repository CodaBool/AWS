########## VARIABLES ##########
locals {
  name = "notify"
  email = "codabool@pm.me"
}

########## RESOURCES ##########
resource "aws_sns_topic" "notify" {
  name              = local.name
  kms_master_key_id = "alias/aws/sns"
}

resource "aws_sns_topic_subscription" "notify" {
  endpoint  = local.email
  protocol  = "email"
  topic_arn = aws_sns_topic.notify.arn
}

resource "aws_iam_role_policy_attachment" "sns_publish" {
  role       = module.lambda.role
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSIoTDeviceDefenderPublishFindingsToSNSMitigationAction"
}

######### LAMBDA ##########
module "lambda" {
  source               = "../lambda"
  name                 = local.name
  path_to_dockerfile   = "${path.module}/src"
  description          = "Sends emails"
  run_on_schedule      = false
}
