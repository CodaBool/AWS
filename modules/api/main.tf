resource "aws_apigatewayv2_api" "main" {
  name          = var.name
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_stage" "main" {
  api_id      = aws_apigatewayv2_api.main.id
  name        = "v1" # env_stage
  # determines if updates to an API automatically trigger a new deployment 
  auto_deploy = true # optional, default is false
  access_log_settings { # optional
    destination_arn = aws_cloudwatch_log_group.main.arn
    format = jsonencode({
      requestId               = "$context.requestId"
      sourceIp                = "$context.identity.sourceIp"
      requestTime             = "$context.requestTime"
      protocol                = "$context.protocol"
      httpMethod              = "$context.httpMethod"
      resourcePath            = "$context.resourcePath"
      routeKey                = "$context.routeKey"
      status                  = "$context.status"
      responseLength          = "$context.responseLength"
      integrationErrorMessage = "$context.integrationErrorMessage"
      }
    )
  }
}

resource "aws_apigatewayv2_integration" "main" {
  api_id = aws_apigatewayv2_api.main.id
  integration_uri    = var.lambda_invoke_arn
  integration_type   = "AWS_PROXY"
  integration_method = "POST"
}

resource "aws_apigatewayv2_route" "main" {
  api_id = aws_apigatewayv2_api.main.id
  route_key = "$default" # ANY /api
  target    = "integrations/${aws_apigatewayv2_integration.main.id}"
}

resource "aws_cloudwatch_log_group" "main" {
  name = "/aws/api/${aws_apigatewayv2_api.main.name}"
  retention_in_days = var.log_retention
}

resource "aws_lambda_permission" "api" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_function_name
  principal     = "apigateway.amazonaws.com"
  source_arn = "${aws_apigatewayv2_api.main.execution_arn}/*/*"
}