module "lambda" {
  source               = "../lambda"
  name                 = var.name
  path_to_dockerfile   = path.module
  tag                  = var.tag
  # log_retention        = 60
  description          = "Scrapes things"
  # this could easily be the method of update instead of crontab but I have the crontab built
  # and having this on schedule would mean having an `update-all` with a `promises.allSettled`
  # added onto the lambda which is more work than keeping existing infrastructure
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

resource "aws_sqs_queue" "scraper" {
  name                       = var.name
  visibility_timeout_seconds = 900
}

resource "aws_lambda_event_source_mapping" "scraper" {
  event_source_arn = aws_sqs_queue.scraper.arn
  function_name    = module.lambda.function.arn
}

# module "api" {
#   source               = "../api"
#   name                 = var.name
#   log_retention        = 60
#   lambda_function_name = module.lambda.function.function_name
#   lambda_invoke_arn    = module.lambda.function.invoke_arn
# }