---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "doit-console Provider"
subcategory: ""
description: |-
  
---

# doit-console Provider



## Example Usage

```terraform
terraform {
  required_providers {
    doit-console = {
      source  = "doitintl/doit-console"
      version = "0.3.1"
    }
  }
}

provider "doit-console" {
  # Configuration options prefer to use environment variables
  # DOIT_API_TOKEN, DOIT_HOST=https://api.doit.com and DOIT_CUSTOMER_CONTEXT
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `api_token` (String, Sensitive) API Token to access DoiT API. May also be provided by DOIT_API_TOKEN environment variable. Refer to https://developer.doit.com/docs/start
- `customer_context` (String) Customer context. May also be provided by DOIT_CUSTOMER_CONTEXT environment variable. This field is requiered just for DoiT employees
- `host` (String) URI for DoiT API. May also be provided via DOIT_HOST environment variable.
