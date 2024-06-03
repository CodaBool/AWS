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
    }]
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

# NEW: lambdas can be invoked by URL
resource "aws_lambda_function_url" "main" {
  count              = var.create_function_url ? 1 : 0
  function_name      = aws_lambda_function.main.function_name
  authorization_type = "NONE"
}

resource "aws_lambda_function" "main" {
  function_name    = var.name
  role             = aws_iam_role.lambda_assume.arn
  package_type     = "Image"
  description      = var.description
  memory_size      = var.memory
  architectures    = [var.architecture]
  timeout          = 900
  image_uri        = "${var.account}.dkr.ecr.us-east-1.amazonaws.com/${var.name}:latest"
  source_code_hash = split("sha256:", data.aws_ecr_image.lambda.id)[1]
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

# Build and push the Docker image on changes
resource "null_resource" "push" {
  triggers = {
    hash = md5(join("", [for f in fileset("${var.path_to_dockerfile}", "*"): filemd5("${var.path_to_dockerfile}/${f}")]))
  }
  provisioner "local-exec" {
    command     = "${path.module}/push.sh ${var.path_to_dockerfile} ${aws_ecr_repository.main.repository_url} ${var.tag}"
    interpreter = ["bash", "-c"]
  }
}