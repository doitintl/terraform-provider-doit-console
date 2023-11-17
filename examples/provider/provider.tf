terraform {
  required_providers {
    doit-console = {
      source  = "doitintl/doit-console"
      version = "0.2.1"
    }
  }
}

provider "doit-console" {
  # Configuration options prefer to use environment variables
  # DOIT_API_TOKEN, DOIT_HOST=https://api.doit.com and DOIT_CUSTOMER_CONTEXT
}
