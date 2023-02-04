variable "name" {
  default = "notify"
  type    = string
}

variable "email" {
  description = "Email to send the SNS topic and results of lambda to"
  default     = "codabool@pm.me"
}

variable "tag" {
  description = "Tag to use for deployed Docker image"
  default     = "latest"
}