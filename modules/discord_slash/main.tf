module "lambda" {
  source              = "../lambda"
  name                = "discord_slash"
  path_to_dockerfile  = "${path.module}/src"
  log_retention       = 60
  memory              = 3072
  description         = "Responds to Slash commands"
  create_function_url = true
  environment = { for tuple in regexall("(.*?)=(.*)", file("${path.module}/src/.env")) : tuple[0] => sensitive(tuple[1]) }
}