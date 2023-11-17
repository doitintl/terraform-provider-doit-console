terraform {
  required_providers {
    doit-console = {
      source  = "doitntl/doit-console"
      version = "0.2.1"
    }
  }
}

resource "doit-console_attribution" "attri" {
  name        = "attritestnewname9"
  description = "attritestdiana8"
  formula     = "A"
  components  = [{ type = "label", key = "iris_location", values = ["us"] }]
}

resource "doit-console_report" "test3" {
  name        = "test10"
  description = "test10"
  config = {
    metric = {
      type  = "basic"
      value = "cost"
    }
    include_promotional_credits = false
    advanced_analysis = {
      trending_up   = false
      trending_down = false
      not_trending  = false
      forecast      = false
    }
    aggregation   = "total"
    time_interval = "month"
    dimensions = [
      {
        id   = "year"
        type = "datetime"
      },
      {
        id   = "month"
        type = "datetime"
      }
    ]
    time_range = {
      mode            = "last"
      amount          = 12
      include_current = true
      unit            = "month"
    }
    filters = [
      {
        inverse = false
        id      = "attribution"
        type    = "attribution"
        values = [
          "1CE699ZdwN5CRBw0tInY"
        ]
      }
    ]
    group = [
      {
        id   = "BSQZmvX6hvuKGPDHX7R3"
        type = "attribution_group"
      },
      {
        id   = "cloud_provider"
        type = "fixed"
      }
    ]
    layout         = "table"
    display_values = "actuals_only"
    currency       = "USD"
  }
}

provider "doit-console" {
}
