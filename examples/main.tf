terraform {
  required_providers {
    doit = {
      source = "hashicorp.com/edu/doit-console"
    }
  }
  backend "gcs" {
    bucket = "iac-doit-console-prod"
    prefix = "terraform/state"
  }

}

resource "doit_attribution" "attri" {
  name        = "attritestnewname9"
  description = "attritestdiana8"
  formula     = "A"
  components  = [{ type = "label", key = "iris_location", values = ["us"] }]
}

provider "doit" {
  host = "https://api.doit.com"
}

resource "doit_attribution" "attribute1" {
  name        = "attritestnewname3"
  description = "attritestdiana8"
  formula     = "A"
  components  = [{ type = "label", key = "iris_location", values = ["us"] }]
}

resource "doit_attribution" "attribute2" {
  name        = "attritestnewname4"
  description = "attritestdiana8"
  formula     = "A"
  components  = [{ type = "label", key = "iris_location", values = ["us"] }]
}

resource "doit_attribution_group" "attributeGroup" {
  name         = "attritestnewgroup1"
  description  = "attritestgroup"
  attributions = [doit_attribution.attribute1.id, doit_attribution.attribute2.id]
}

