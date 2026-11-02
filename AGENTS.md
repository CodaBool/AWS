# AGENTS.md

Personal AWS infrastructure in Terraform, all in one account (`us-east-1`, state in the `codabool-tf` S3 bucket).

## Layout

- `main.tf` is the source of truth for what is deployed. A module that `main.tf` does not reference is dormant. Today it wires only `key`, `actions`, and the `/env` SSM backup; no Lambdas are live.
- `DEPRECATED-*` modules (under `modules/vault/`, plus `modules/DEPRECATED-discord_slash`) are retired Lambdas kept as reference code only. Treat them as read-only; when one needs to come back, move it out of `vault/` and wire it into `main.tf`.
- `modules/lambda` is the shared Lambda factory: ECR repo, image build + push, IAM role, log group, optional EventBridge schedule and function URL. A service module wraps it with `source = "../lambda"` and keeps its code in `src/`, which is self-contained (own `package.json`/`go.mod`, `Dockerfile`, `justfile`).
- The `sock` EC2 instance (t4g.nano + 8 GB gp3) is managed outside Terraform; `modules/ec2` is dormant.

## Deploying is live

`terraform apply` changes production directly. When any file in a wired module's `src/` changes, the `null_resource.push` in `modules/lambda` runs `push.sh`, which builds the image with `docker buildx --platform linux/arm64` and pushes it to ECR before the Lambda updates. Run `terraform plan` freely; ask before `apply`.

The `discord_slash` Lambda and its ECR repo were deleted by hand (2026-11-01), so state may still hold `module.lambda_discord_slash` resources. Expect the next plan to drop them; if it errors on them, `terraform state rm` the leftovers.

## Secrets and `.env`

- A wired module reads `src/.env` at plan time (`file("${path.module}/src/.env")`) to set Lambda environment variables. The file is gitignored, so a fresh checkout has none and `plan` fails until it exists.
- `readenv.sh` collects every `.env` file in the tree outside `node_modules`, `DEPRECATED-*` included, and stores them in the SSM SecureString `/env` as a backup. A new `.env` anywhere in the repo gets uploaded on the next apply.
- `.env` and `*.pem` hold live credentials: keep their values out of output, commits, and docs.

## Cost

The owner keeps the monthly bill to `sock` (about $4.73) and treats anything else over $0.50/mo as worth removing. Query Cost Explorer through the AWS MCP (`GetCostAndUsage` grouped by `SERVICE` and `USAGE_TYPE`); history goes back 14 months. Two known orphaned charges have no resource left to delete and are with AWS Support:

- `TimedStorage-GlacierStaging` (~$0.64/mo, ~30 GB): leftover uploads from the deleted bucket `codabool-backups`, visible only in CloudWatch `AWS/S3` metrics.
- `EBS:SnapshotUsage` (~$0.09/mo, ~1.9 GB): billed with no snapshots in any region, Recycle Bin, or AWS Backup.

## Lambda conventions

- Architecture is `arm64` unless a module overrides it; base images and build flags in each `Dockerfile` match.
- Local testing goes through each `src/justfile` (`just run`, then `just test`); it expects AWS credentials already exported in the shell (the owner uses `awsume`).
