module "lambda" {
  source               = "../lambda"
  name                 = "scraper"
  path_to_dockerfile   = "${path.module}/src"
  tag                  = "latest"
  log_retention        = 60 # default is 7
  memory               = 3072 # default 512, 3009 unlocks 3vCPU
  description          = "Scrapes things"
  # this could easily be the method of update instead of crontab but I have the crontab built
  # and having this on schedule would mean having an `update-all` with a `promises.allSettled`
  # added onto the lambda which is more work than keeping existing infrastructure
  interval             = "cron(0 12 1 * ? *)" # 1st of every month 7am est
  # event_input        = jsonencode({
  #   path = "/v1/upcoming_movies"
  #   queryStringParameters = {
  #     key = local.env["KEY"]
  #     limit = 25
  #   }
  # })
  environment = local.env
}

locals {
  env = { for tuple in regexall("(.*)=(.*)", file("${path.module}/src/.env")) : tuple[0] => tuple[1] }
}