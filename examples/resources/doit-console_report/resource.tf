resource "doit_report" "my_report" {
  name        = "Test Terraform update"
  description = ""
  config = {
    metric = {
      type  = "basic"
      value = "cost"
    }
    aggregation = "total"
    advanced_analysis = {
      forecast      = false
      not_trending  = false
      trending_down = false
      trending_up   = true
    }
    time_interval = "day"
    dimensions = [
      {
        id   = "year"
        type = "datetime"
      },
      {
        id   = "month"
        type = "datetime"
      },
      {
        id   = "day"
        type = "datetime"
      },
    ]
    time_range = {
      amount          = 7
      include_current = false
      mode            = "last"
      unit            = "day"
    }
    include_promotional_credits = false
    layout                      = "stacked_column_chart"
    display_values              = "actuals_only"
    currency                    = "USD"
  }
}