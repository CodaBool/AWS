module "lambda" {
  source               = "../lambda"
  name                 = "scraper"
  path_to_dockerfile   = "${path.module}/src"
  memory               = 3072
  description          = "Scrapes things"
  account              = var.account
  interval             = "cron(0 12 1 * ? *)" # 1st of every month 7am est
  environment = { for tuple in regexall("(.*?)=(.*)", file("${path.module}/src/.env")) : tuple[0] => sensitive(tuple[1]) }
}

variable "account" {
  type = string
}