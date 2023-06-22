module "lambda" {
  source               = "../lambda"
  name                 = "discord"
  path_to_dockerfile   = "${path.module}/src"
  log_retention        = 60
  memory               = 3072
  account              = var.account
  description          = "Posts to Discord"
  interval             = "cron(0 14 1 * ? *)" # 1st of every month 9am est
  environment = { for tuple in regexall("(.*?)=(.*)", file("${path.module}/src/.env")) : tuple[0] => sensitive(tuple[1]) }
}

variable "account" {
  type = string
}