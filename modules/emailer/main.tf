# module "lambda" {
#   source             = "../lambda"
#   name               = var.name
#   path_to_dockerfile = path.module
#   tag                = var.tag
#   description        = "Emails me the size of a path on my server"
#   #   event_input        = <<EOF
#   # {
#   #   "notify_before": ${var.notify_before},
#   #   "expired_topic_arn": "${aws_sns_topic.expiring.arn}",
#   #   "expiring_topic_arn": "${aws_sns_topic.expired.arn}"
#   # }
#   # EOF
# }

# resource "aws_iam_role_policy_attachment" "sns_publish" {
#   role       = module.lambda.role
#   policy_arn = "arn:aws:iam::aws:policy/service-role/AWSIoTDeviceDefenderPublishFindingsToSNSMitigationAction"
# }

resource "aws_sns_topic" "email" {
  name              = var.name
  kms_master_key_id = "alias/aws/sns"
}

resource "aws_sns_topic_subscription" "email" {
  endpoint  = var.email
  protocol  = "email"
  topic_arn = aws_sns_topic.email.arn
}