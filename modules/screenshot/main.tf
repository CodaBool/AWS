module "lambda" {
  source             = "../lambda"
  name               = "screenshot"
  path_to_dockerfile = "${path.module}/src"
  # log_retention       = 60
  memory              = 4096
  description         = "Takes a screenshot"
  create_function_url = true
  architecture        = "x86_64"
}
