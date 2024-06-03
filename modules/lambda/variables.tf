# https://docs.aws.amazon.com/eventbridge/latest/userguide/eb-cron-expressions.html
variable "interval" {
  default     = "cron(0 16 ? * MON-FRI *)" # 12am EST every weekday
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
  default     = 60
  description = "how long to retain lambda execution log data in cloudwatch logs"
}

variable "run_on_schedule" {
  default     = true
  description = "whether to create cloudwatch event scheduling resources"
}

variable "name" {
  type = string
}

variable "account" {
  type = string
  default = "919759177803"
}

variable "memory" {
  type = number
  default = 512
}

variable "event_input" {
  default = ""
  type    = string
}

variable "description" {
  type = string
}

variable "architecture" {
  default = "arm64"
  type = string
}

variable "create_function_url" {
  default = false
  type = bool
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