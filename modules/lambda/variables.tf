# https://docs.aws.amazon.com/eventbridge/latest/userguide/eb-create-rule-schedule.html#eb-cron-expressions
variable "interval" {
  default     = "cron(0 16 ? * MON-FRI *)" # 11am EST every weekday
  description = "How often to invoke the function in UTC"
}

variable "environment" {
  description = "The environment variables to pass to the lambda"
  sensitive = true
  default = {
    id = "default"
  }
}

variable "log_retention" {
  default     = 7
  description = "how long to retain lambda execution log data in cloudwatch logs"
}

variable "run_on_schedule" {
  default     = true
  description = "whether to create cloudwatch event scheduling resources"
}

variable "name" {
  type = string
}

variable "event_input" {
  default = ""
  type    = string
}

variable "description" {
  type = string
}

variable "tag" {
  description = "Tag to use for deployed Docker image"
  default     = "latest"
}

variable "notify_before" {
  default     = 30
  description = "how many days in the future to look for certificate expiration"
}

variable "path_to_dockerfile" {
  description = "Path to Docker image source"
  type        = string
}