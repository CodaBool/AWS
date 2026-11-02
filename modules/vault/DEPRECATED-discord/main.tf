module "lambda" {
  source               = "../lambda"
  name                 = "discord"
  path_to_dockerfile   = "${path.module}/src"
  memory               = 3072
  description          = "Posts to Discord"
  interval             = "cron(0 14 1 * ? *)" # 1st of every month 10am est
  event_input          = jsonencode({ test = false })
  environment = { for tuple in regexall("(.*?)=(.*)", file("${path.module}/src/.env")) : tuple[0] => sensitive(tuple[1]) }
}
