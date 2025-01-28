module "lambda" {
  source              = "../lambda"
  name                = "discord-reminder"
  path_to_dockerfile  = "${path.module}/src"
  description         = "Send a reminder to Discord"
  create_function_url = true
  interval            = "cron(30 18 * * ? *)" # can -5 hours to get EST time
  event_input = jsonencode({
    QueryStringParameters : {
      action : "cron"
    }
  })
  environment = { for tuple in regexall("(.*?)=(.*)", file("${path.module}/src/.env")) : tuple[0] => sensitive(tuple[1]) }
}
