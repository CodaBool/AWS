variable "log_retention" {
  default     = 7
  description = "how long to retain lambda execution log data in cloudwatch logs"
}

variable "name" {
  type = string
}

variable "lambda_function_name" {
  type = string
}

variable "lambda_invoke_arn" {
  type = string
}