module "lambda" {
  source               = "../lambda"
  name                 = "hibernate"
  path_to_dockerfile   = "${path.module}/src"
  description          = "Automated starting and stopping EC2 to save on costs"
  interval             = "cron(0 16 * * ? *)" # Every day at 12pm
  account              = var.account
  event_input          = jsonencode({ start = true })
}

resource "aws_iam_role_policy_attachment" "ec2_full_access" {
  role       = module.lambda.role
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2FullAccess"
}

resource "aws_iam_role_policy_attachment" "sns_publish" {
  role       = module.lambda.role
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSIoTDeviceDefenderPublishFindingsToSNSMitigationAction"
}

resource "aws_lambda_permission" "allow_cloudwatch" {
  statement_id  = "AllowExecutionFromCloudWatchSecondary"
  action        = "lambda:InvokeFunction"
  function_name = module.lambda.function.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.event_rule.arn
}

resource "aws_cloudwatch_event_rule" "event_rule" {
  name_prefix         = "scheduled-${module.lambda.function.function_name}"
  schedule_expression = "cron(0 4 * * ? *)" # Every day at 12am
  description         = "Invoke the ${module.lambda.function.function_name} Lambda function"
}

resource "aws_cloudwatch_event_target" "lambda" {
  rule  = aws_cloudwatch_event_rule.event_rule.id
  arn   = module.lambda.function.arn
  input = jsonencode({ start = false })
}

variable "account" {
  type = string
}