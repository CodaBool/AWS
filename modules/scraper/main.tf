module "lambda" {
  source               = "../lambda"
  name                 = var.name
  path_to_dockerfile   = path.module
  tag                  = var.tag
  log_retention        = 60
  description          = "Scrapes things"
  run_on_schedule      = false
  # interval             = "cron(0 12 1 * ? *)" # 1st of every month 7am est
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
  env = { for tuple in regexall("(.*)=(.*)", file("${path.module}/.env")) : tuple[0] => tuple[1] }
}

module "api" {
  source               = "../api"
  name                 = var.name
  log_retention        = 60
  lambda_function_name = module.lambda.function.function_name
  lambda_invoke_arn    = module.lambda.function.invoke_arn
}