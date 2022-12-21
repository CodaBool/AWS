output "role" {
  description = "Role name to attach policy to"
  value       = aws_iam_role.lambda_assume.name
}

output "function" {
  value       = aws_lambda_function.main
}

# output "lambda_function_invoke_arn" {
#   description = "Role name to attach policy to"
#   value       = aws_lambda_function.main.invoke_arn
# }

# output "lambda_function_name" {
#   description = "Role name to attach policy to"
#   value       = aws_lambda_function.main.function_name
# }