terraform {
  required_providers {
    doit = {
      source = "hashicorp.com/edu/doit-console"
    }
  }
}

resource "doit_attribution" "attri" {
  name="attritestnewname8"
  description="attritestdiana8"
  formula="A"
  components=[{type="label", key="iris_location", values=["us"]}]
}

provider "doit" {
  host="https://api.doit.com"
}


git filter-repo --invert-paths --path examples/provider-install-verification/main copy.tf.bk

