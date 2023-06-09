resource "aws_lambda_permission" "allow_cloudwatch" {
  count = var.run_on_schedule ? 1 : 0
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.main.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.event_rule[0].arn
}

resource "aws_cloudwatch_event_rule" "event_rule" {
  count = var.run_on_schedule ? 1 : 0
  name_prefix         = "scheduled-${aws_lambda_function.main.function_name}"
  schedule_expression = var.interval
  description         = "Invoke the ${aws_lambda_function.main.function_name} Lambda function"
}

resource "aws_cloudwatch_event_target" "lambda" {
  count = var.run_on_schedule ? 1 : 0
  rule  = aws_cloudwatch_event_rule.event_rule[0].id
  arn   = aws_lambda_function.main.arn
  input = var.event_input
}

data "aws_caller_identity" "current" {}

resource "aws_iam_role" "lambda_assume" {
  name               = "${var.name}-lambda-assume"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "lambda.amazonaws.com"
      }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "logs_create" {
  role       = aws_iam_role.lambda_assume.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_cloudwatch_log_group" "delete_old_logs" {
  name              = "/aws/lambda/${aws_lambda_function.main.function_name}"
  retention_in_days = var.log_retention
}

resource "aws_lambda_function" "main" {
  function_name    = var.name
  role             = aws_iam_role.lambda_assume.arn
  package_type     = "Image"
  description      = var.description
  memory_size      = 512
  timeout          = 900
  image_uri        = "${data.aws_caller_identity.current.account_id}.dkr.ecr.us-east-1.amazonaws.com/${var.name}:latest"
  source_code_hash = split("sha256:", data.aws_ecr_image.lambda.id)[1]
  # source_code_hash = filemd5("../dist/${each.value}")
  environment {
    variables = var.environment
  }
}

data "aws_ecr_image" "lambda" {
  depends_on      = [null_resource.push]
  repository_name = var.name
  image_tag       = var.tag
}

resource "aws_ecr_repository" "main" {
  name = var.name
  force_delete = true
  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_ecr_lifecycle_policy" "remove_old_images" {
  repository = aws_ecr_repository.main.name
  # example use of keeping certain tag image
  # https://github.com/mathspace/terraform-aws-ecr-docker-image/blob/master/main.tf
  policy = jsonencode({
    rules = [{
      rulePriority = 1
      description = "Delete untagged images"
      action = { 
        type = "expire" 
      }
      selection = {
        tagStatus = "untagged"
        countType = "sinceImagePushed"
        countUnit = "days"
        countNumber = 1
      }
    }]
  })
}

# Necessary since the initial push would have relied on 
data "external" "hash" {
  program = ["${path.module}/hash.sh", var.path_to_dockerfile]
}

# Build and push the Docker image whenever the hash changes
resource "null_resource" "push" {
  triggers = {
    hash = data.external.hash.result["hash"]
  }
  provisioner "local-exec" {
    command     = "${path.module}/push.sh ${var.path_to_dockerfile} ${aws_ecr_repository.main.repository_url} ${var.tag}"
    interpreter = ["bash", "-c"]
  }
}