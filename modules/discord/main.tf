module "lambda" {
  source               = "../lambda"
  name                 = "discord"
  path_to_dockerfile   = "${path.module}/src"
  tag                  = "latest"
  log_retention        = 60 # default is 7
  memory               = 3072 # default 512, 3009 unlocks 3vCPU
  description          = "Scrapes things"
  interval             = "cron(0 14 1 * ? *)" # 1st of every month 9am est
  environment = local.env
}

locals {
  env = { for tuple in regexall("(.*?)=(.*)", file("${path.module}/src/.env")) : tuple[0] => sensitive(tuple[1]) }
}