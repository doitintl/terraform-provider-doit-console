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

resource "doit_attribution" "attribute1" {
  name="attritestnewname1"
  description="attritestdiana8"
  formula="A"
  components=[{type="label", key="iris_location", values=["us"]}]
}

resource "doit_attribution" "attribute2" {
  name="attritestnewname2"
  description="attritestdiana8"
  formula="A"
  components=[{type="label", key="iris_location", values=["us"]}]
}

resource "doit_attribution_group" "attributeGroup" {
  name="attritestnewgroup"
  description="attritestgroup"
  attributions=[doit_attribution.attribute1.id, doit_attribution.attribute2.id]
}

