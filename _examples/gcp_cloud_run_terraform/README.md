# gcp_cloud_run_terraform
Terraform example for deploy to Cloud Run

## Usage
1. Copy `*.tf` files
2. Put and edit [terraform.tfvars](terraform.tfvars), or use as module

e.g.

```tf
module "feed_squeezer" {
  source = "./modules/path/to/feed_squeezer"

  tag = "latest"
}
```
