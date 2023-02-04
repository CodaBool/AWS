resource "aws_sns_topic" "notify" {
  name              = var.name
  kms_master_key_id = "alias/aws/sns"
}

resource "aws_sns_topic_subscription" "notify" {
  endpoint  = var.email
  protocol  = "email"
  topic_arn = aws_sns_topic.notify.arn
}

module "lambda" {
  source               = "../lambda"
  name                 = var.name
  path_to_dockerfile   = path.module
  tag                  = var.tag
  log_retention        = 60
  description          = "Capable of sending emails or texts"
  # this could easily be the method of update instead of crontab but I have the crontab built
  # and having this on schedule would mean having an `update-all` with a `promises.allSettled`
  # added onto the lambda which is more work than keeping existing infrastructure
  run_on_schedule      = false 
  # interval             = "cron(0 12 1 * ? *)" # 1st of every month 7am est
  # event_input        = jsonencode({
  #   path = "/v1/upcoming_movies"
  #   queryStringParameters = {
  #     key = local.env["KEY"]
  #     limit = 25
  #   }
  # })
  environment = local.env
}

resource "aws_iam_role_policy_attachment" "sns_publish" {
  role       = module.lambda.role
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSIoTDeviceDefenderPublishFindingsToSNSMitigationAction"
}

locals {
  env = { for tuple in regexall("(.*)=(.*)", file("${path.module}/.env")) : tuple[0] => tuple[1] }
}