data "aws_iam_policy" "admin_access" {
  name = "AdministratorAccess"
}

resource "aws_iam_openid_connect_provider" "default" {
  url = "https://token.actions.githubusercontent.com"
  client_id_list = ["sts.amazonaws.com"]
  thumbprint_list = ["6938fd4d98bab03faadb97b34396831e3780aea1"]
}

# for some reason terraform thinks this is invalid policy but it's an exact match with what is in the account 🤷
resource "aws_iam_role" "actions" {
  name        = "gh-action-assume"
  description = "Role to assume to create the infrastructure."
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow",
        Principal = {
          "Federated" : "arn:aws:iam::${var.account}:oidc-provider/token.actions.githubusercontent.com"
        },
        Action = "sts:AssumeRoleWithWebIdentity",
        Condition = {
          StringLike = {
            "token.actions.githubusercontent.com:aud": "sts.amazonaws.com",
            "token.actions.githubusercontent.com:sub": [
              "repo:CodaBool/*",
              "repo:TampaDevs/bayhacks.dev:*"
            ]
          }
        }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "githubiac" {
  role       = aws_iam_role.actions.name
  policy_arn = data.aws_iam_policy.admin_access.arn
}

output "role_arn" {
  value       = aws_iam_role.actions.arn
  description = "github actions assume role"
}

variable "github_repository" {
  default = "CodaBool/p12-slap"
}

variable "account" {
  type = string
}