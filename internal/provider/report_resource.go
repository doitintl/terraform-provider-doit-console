package provider

import (
	"context"
	"fmt"
	"time"

	"log"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// orderResourceModel maps the resource schema data.
type reportResourceModel struct {
	// Config Report configuration
	Config *ExternalConfigModel `tfsdk:"config"`
	// Description Report description
	Description types.String `tfsdk:"description"`
	// Id Report id. Leave blank when creating a new report
	Id types.String `tfsdk:"id"`
	// Name Report name
	Name        types.String `tfsdk:"name"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

// ExternalConfig Report configuration
type ExternalConfigModel struct {
	// AdvancedAnalysis Advanced analysis toggles. Each of these can be set independently
	AdvancedAnalysis *AdvancedAnalysisModel `tfsdk:"advanced_analysis"`
	Aggregation      types.String           `tfsdk:"aggregation"`
	Currency         types.String           `tfsdk:"currency"`
	Dimensions       []DimensionModel       `tfsdk:"dimensions"`
	DisplayValues    types.String           `tfsdk:"display_values"`

	// Filters The filters to use in this report
	Filters []ExternalConfigFilterModel `tfsdk:"filters"`

	// Group The groups to use in the report.
	Group []GroupModel `tfsdk:"group"`

	// IncludePromotionalCredits Whether to include credits or not.
	// If set, the report must use time interval “month”/”quarter”/”year”
	IncludePromotionalCredits types.Bool          `tfsdk:"include_promotional_credits"`
	Layout                    types.String        `tfsdk:"layout"`
	Metric                    ExternalMetricModel `tfsdk:"metric"`

	// MetricFilter {
	// "metric": {
	// "type":  "basic",
	// "value": "cost"
	// },
	// "operator" : "gt",
	// "values" : [50]
	// }
	MetricFilter *ExternalConfigMetricFilterModel `tfsdk:"metric_filter"`

	// Splits The splits to use in the report.
	Splits       []ExternalSplitModel `tfsdk:"splits"`
	TimeInterval types.String         `tfsdk:"time_interval"`

	// TimeRange Time settings for the report
	// Description: Today is the 17th of April of 2023
	// We set the mode to "last", the amount to 2 and the unit to "day"
	// If includeCurrent is not set, the range will be the 15th and 16th of April
	// If it is, then the range will be 16th and 17th
	TimeRange *TimeSettingsModel `tfsdk:"time_range"`
}

// AdvancedAnalysis Advanced analysis toggles. Each of these can be set independently
type AdvancedAnalysisModel struct {
	Forecast     types.Bool `tfsdk:"forecast"`
	NotTrending  types.Bool `tfsdk:"not_trending"`
	TrendingDown types.Bool `tfsdk:"trending_down"`
	TrendingUp   types.Bool `tfsdk:"trending_up"`
}

// GroupModel represents a group in the report.
type GroupModel struct {
	Id    types.String `tfsdk:"id"`
	Type  types.String `tfsdk:"type"`
	Limit LimitModel   `tfsdk:"limit"`
}

type LimitModel struct {
	Metric ExternalMetricModel `tfsdk:"metric"`
	Sort   types.String        `tfsdk:"sort"`
	// Value The number of items to show
	Value types.Int64 `tfsdk:"value"`
}

type ExternalMetricModel struct {
	Type  types.String `tfsdk:"type"`
	Value types.String `tfsdk:"value"`
}

type ExternalConfigMetricFilterModel struct {
	Metric   *ExternalMetricModel `tfsdk:"metric"`
	Operator types.String         `tfsdk:"operator"`
	Values   []types.Float64      `tfsdk:"values"`
}

// ExternalSplitModel represents a split in the report.
type ExternalSplitModel struct {
	// Id ID of the field to split
	Id types.String `tfsdk:"id"`

	// IncludeOrigin if set, include the origin
	IncludeOrigin types.Bool          `tfsdk:"include_origin"`
	Mode          types.String        `tfsdk:"mode"`
	Origin        ExternalOriginModel `tfsdk:"origin"`

	// Targets Targets for the split
	Targets []ExternalSplitTargetModel `tfsdk:"targets"`

	// Type Type of the split.
	// The only supported value at the moment: "attribution_group"
	Type types.String `tfsdk:"type"`
}

// ExternalOrigin defines model for ExternalOrigin.
type ExternalOriginModel struct {
	// Id ID of the origin
	Id types.String `tfsdk:"id"`
	// Type Type of the origin.
	// The only supported value at the moment: "attribution"
	Type types.String `tfsdk:"type"`
}

// ExternalSplitTarget defines model for ExternalSplitTargetModel.
type ExternalSplitTargetModel struct {
	// Id ID of the target
	Id types.String `tfsdk:"id"`
	// Type Type of the target.
	// The only supported value at the moment: "target"
	Type types.String `tfsdk:"type"`
	// Value Percent of the target, represented in float format. E.g. 30% is 0.3. Must be set only if Split Mode is custom.
	Value types.Float64 `tfsdk:"value"`
}

// TimeSettings Time settings for the report
// Description: Today is the 17th of April of 2023
// We set the mode to "last", the amount to 2 and the unit to "day"
// If includeCurrent is not set, the range will be the 15th and 16th of April
// If it is, then the range will be 16th and 17th
type TimeSettingsModel struct {
	Amount         types.Int64  `tfsdk:"amount"`
	IncludeCurrent types.Bool   `tfsdk:"include_current"`
	Mode           types.String `tfsdk:"mode"`
	Unit           types.String `tfsdk:"unit"`
}

// Dimension {
// "id" : "sku_description",
// "type" : "fixed"
// }
type DimensionModel struct {
	// Id The field to apply to the dimension.
	Id   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

// ExternalConfigFilter {
// "id" : "sku_description",
// "type" : "fixed",
// "values" : ["Nearline Storage Iowa", "Nearline Storage Frankfurt"]
// }
//
// When using attributions as a filter both the type and the ID must be "attribution", and the
// values array contains the attribution IDs.
type ExternalConfigFilterModel struct {
	// Id What field we are filtering on
	Id types.String `tfsdk:"id"`
	// Inverse If set, exclude the values
	Inverse types.Bool   `tfsdk:"inverse"`
	Type    types.String `tfsdk:"type"`
	// Values What values to filter on or exclude
	Values []types.String `tfsdk:"values"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &reportResource{}
	_ resource.ResourceWithConfigure = &reportResource{}
)

// NewreportResource is a helper function to simplify the provider implementation.
func NewReportResource() resource.Resource {
	return &reportResource{}
}

// reportResource is the resource implementation.
type reportResource struct {
	client *ClientTest
}

// Metadata returns the resource type name.
func (r *reportResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	log.Print("hello report Metadata:)")
	resp.TypeName = req.ProviderTypeName + "_report"
}

// Schema defines the schema for the resource.
func (r *reportResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	log.Print("hello report Schema:)")
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"config": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"advanced_analysis": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"forecast": schema.BoolAttribute{
								Description: "Advanced analysis toggles. Each of these can be set independently",
								Optional:    true,
							},
							"not_trending": schema.BoolAttribute{
								Description: "",
								Optional:    true,
							},
							"trending_down": schema.BoolAttribute{
								Description: "",
								Optional:    true,
							},
							"trending_up": schema.BoolAttribute{
								Description: "",
								Optional:    true,
							},
						},
						Description: "",
						Optional:    true,
					},
					"aggregation": schema.StringAttribute{
						Description: "",
						Optional:    true,
					},
					"currency": schema.StringAttribute{
						Description: "",
						Optional:    true,
					},
					"dimensions": schema.ListNestedAttribute{
						Description: "",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Description: "",
									Optional:    true,
								},
								"type": schema.StringAttribute{
									Description: "",
									Optional:    true,
								},
							},
						},
					},
					"display_values": schema.StringAttribute{
						Description: "",
						Optional:    true,
					},
					"filters": schema.ListNestedAttribute{
						Description: "The filters to use in this report",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Description: "What field we are filtering on",
									Optional:    true,
								},
								"inverse": schema.BoolAttribute{
									Description: "If set, exclude the values",
									Optional:    true,
								},
								"type": schema.StringAttribute{
									Description: "",
									Optional:    true,
								},
								"values": schema.ListAttribute{
									Description: "What values to filter on or exclude",
									ElementType: types.StringType,
									Required:    true,
								},
							},
						},
					},
					"group": schema.ListNestedAttribute{
						Description: "The groups to use in the report.",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Description: "",
									Optional:    true,
								},
								"type": schema.StringAttribute{
									Description: "",
									Optional:    true,
								},
								"limit": schema.SingleNestedAttribute{
									Attributes: map[string]schema.Attribute{
										"metric": schema.SingleNestedAttribute{
											Attributes: map[string]schema.Attribute{
												"type": schema.StringAttribute{
													Description: "",
													Optional:    true,
												},
												"value": schema.StringAttribute{
													Description: "",
													Optional:    true,
												},
											},
											Description: "",
											Optional:    true,
										},
										"sort": schema.StringAttribute{
											Description: "",
											Optional:    true,
										},
										"value": schema.Int64Attribute{
											Description: "",
											Optional:    true,
										},
									},
									Description: "",
									Optional:    true,
								},
							},
						},
					},
					"include_promotional_credits": schema.BoolAttribute{
						Description: "Whether to include credits or not. "+
							"If set, the report must use time interval “month”/”quarter”/”year”",
						Optional:    true,
					},
					"layout": schema.StringAttribute{
						Description: "",
						Optional:    true,
					},
					"metric": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Description: "",
								Optional:    true,
							},
							"value": schema.StringAttribute{
								Description: "For basic metrics the value can be one of: [\"cost\", \"usage\", \"savings\" \n"+
								"If using custom metrics, the value must refer to an existing custom or calculated metric id ",
								Optional:    true,
							},
						},
						Description: "",
						Optional:    true,
					},
					"metric_filter": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"metric": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"type": schema.StringAttribute{
										Description: "",
										Optional:    true,
									},
									"value": schema.StringAttribute{
										Description: "",
										Optional:    true,
									},
								},
								Description: "",
								Optional:    true,
							},
							"operator": schema.StringAttribute{
								Description: "",
								Optional:    true,
							},
							"values": schema.ListAttribute{
								Description: "",
								ElementType: types.Float64Type,
								Required:    true,
							},
						},
						Description: "",
						Optional:    true,
					},
					"splits": schema.ListNestedAttribute{
						Description: "The splits to use in the report.",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Description: "",
									Optional:    true,
								},
								"include_origin": schema.BoolAttribute{
									Description: "",
									Optional:    true,
								},
								"type": schema.StringAttribute{
									Description: "",
									Optional:    true,
								},
								"mode": schema.StringAttribute{
									Description: "",
									Optional:    true,
								},
								"origin": schema.SingleNestedAttribute{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Description: "",
											Optional:    true,
										},
										"type": schema.StringAttribute{
											Description: "",
											Optional:    true,
										},
									},
									Description: "",
									Optional:    true,
								},
								"targets": schema.ListNestedAttribute{
									Description: "",
									Optional:    true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												Description: "",
												Optional:    true,
											},
											"type": schema.StringAttribute{
												Description: "",
												Optional:    true,
											},
										},
									},
								},
							},
						},
					},
					"time_interval": schema.StringAttribute{
						Description: "",
						Optional:    true,
					},
					"time_range": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"amount": schema.Int64Attribute{
								Description: "",
								Optional:    true,
							},
							"include_current": schema.BoolAttribute{
								Description: "",
								Optional:    true,
							},
							"mode": schema.StringAttribute{
								Description: "",
								Optional:    true,
							},
							"unit": schema.StringAttribute{
								Description: "",
								Optional:    true,
							},
						},
						Description: "",
						Optional:    true,
					},
				},
				Description: "Report configuration",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "Report description",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "Report name",
				Optional:    true,
			},
			"last_updated": schema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Report id",
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *reportResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	log.Print("hello report Configure:)")
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*ClientTest)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ClientTest, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *reportResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	log.Println("hello report Create:)")
	log.Println(r.client.Auth.DoiTAPITOken)
	log.Println("---------------------------------------------------")
	log.Println(r.client.Auth.CustomerContext)

	// Retrieve values from plan
	var plan reportResourceModel
	log.Println("before getting plan")
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Println("after getting plan")
	// Generate API request body from plan

	config := ExternalConfig{}
	advancedAnalysis := AdvancedAnalysis{
		Forecast:     plan.Config.AdvancedAnalysis.Forecast.ValueBool(),
		NotTrending:  plan.Config.AdvancedAnalysis.NotTrending.ValueBool(),
		TrendingDown: plan.Config.AdvancedAnalysis.TrendingDown.ValueBool(),
		TrendingUp:   plan.Config.AdvancedAnalysis.TrendingUp.ValueBool(),
	}
	config.AdvancedAnalysis = &advancedAnalysis
	var dimensions []Dimension
	for _, dimension := range plan.Config.Dimensions {
		dimension := Dimension{
			Id:   dimension.Id.ValueString(),
			Type: dimension.Type.ValueString(),
		}
		dimensions = append(dimensions, dimension)
	}
	config.Dimensions = dimensions
	var filters []ExternalConfigFilter

	for _, filter := range plan.Config.Filters {
		var values []string
		for _, value := range filter.Values {
			values = append(values, value.ValueString())
		}
		filter := ExternalConfigFilter{
			Id:      filter.Id.ValueString(),
			Inverse: filter.Inverse.ValueBool(),
			Type:    filter.Type.ValueString(),
			Values:  values,
		}
		filters = append(filters, filter)
	}
	config.Filters = filters
	var groups []Group
	for _, group := range plan.Config.Group {
		limit := Limit{
			Metric: ExternalMetric{
				Type:  group.Limit.Metric.Type.ValueString(),
				Value: group.Limit.Metric.Value.ValueString(),
			},
			Sort:  group.Limit.Sort.ValueString(),
			Value: group.Limit.Value.ValueInt64(),
		}
		groups = append(groups, Group{
			Id:    group.Id.ValueString(),
			Type:  group.Type.ValueString(),
			Limit: limit,
		})
	}
	config.Group = groups
	var splits []ExternalSplit
	for _, split := range plan.Config.Splits {
		origin := ExternalOrigin{
			Id:   split.Origin.Id.ValueString(),
			Type: split.Origin.Type.ValueString(),
		}
		targets := []ExternalSplitTarget{}
		for _, target := range split.Targets {
			target := ExternalSplitTarget{
				Id:    target.Id.ValueString(),
				Type:  target.Type.ValueString(),
				Value: target.Value.ValueFloat64(),
			}
			targets = append(targets, target)
		}
		split := ExternalSplit{
			Id:            split.Id.ValueString(),
			IncludeOrigin: split.IncludeOrigin.ValueBool(),
			Mode:          split.Mode.ValueString(),
			Origin:        origin,
			Targets:       targets,
			Type:          split.Type.ValueString(),
		}
		splits = append(splits, split)
	}
	config.Splits = splits
	timeRange := TimeSettings{
		Amount:         plan.Config.TimeRange.Amount.ValueInt64(),
		IncludeCurrent: plan.Config.TimeRange.IncludeCurrent.ValueBool(),
		Mode:           plan.Config.TimeRange.Mode.ValueString(),
		Unit:           plan.Config.TimeRange.Unit.ValueString(),
	}
	if plan.Config.TimeRange != nil {
		config.TimeRange = &timeRange
	}

	var metricFilter ExternalConfigMetricFilter
	log.Println("plan.Config.MetricFilter")
	log.Println(plan.Config.MetricFilter)
	if plan.Config.MetricFilter != nil {
		var values []float64
		for _, value := range plan.Config.MetricFilter.Values {
			values = append(values, value.ValueFloat64())
		}
		metricInFilter := ExternalMetric{
			Type:  plan.Config.MetricFilter.Metric.Type.ValueString(),
			Value: plan.Config.MetricFilter.Metric.Value.ValueString(),
		}
		operator := plan.Config.MetricFilter.Operator.ValueString()
		metricFilter = ExternalConfigMetricFilter{
			Metric:   &metricInFilter,
			Operator: operator,
			Values:   values,
		}
		config.MetricFilter = &metricFilter
	}

	metric := ExternalMetric{
		Type:  plan.Config.Metric.Type.ValueString(),
		Value: plan.Config.Metric.Value.ValueString(),
	}
	config.Metric = metric
	/*config := ExternalConfig{
		AdvancedAnalysis:          advancedAnalysis,
		Aggregation:               plan.Config.Aggregation.ValueString(),
		Currency:                  plan.Config.Currency.ValueString(),
		Dimensions:                dimensions,
		Filters:                   filters,
		Group:                     groups,
		DisplayValues:             plan.Config.DisplayValues.ValueString(),
		IncludePromotionalCredits: plan.Config.IncludePromotionalCredits.ValueBool(),
		Layout:                    plan.Config.Layout.ValueString(),
		Metric:                    metric,
		Splits:                    splits,
		TimeRange:                 timeRange,
		TimeInterval:              plan.Config.TimeInterval.ValueString(),
	}*/
	report := Report{
		Config:      config,
		Description: plan.Description.ValueString(),
		Id:          plan.Id.ValueString(),
		Name:        plan.Name.ValueString(),
	}

	log.Println("before creating report")
	// Create new report
	budgeResponse, err := r.client.CreateReport(report)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating report",
			"Could not create report, unexpected error: "+err.Error(),
		)
		return
	}
	log.Println("report response---------------------------------------------------")
	log.Println(budgeResponse)
	log.Println("report id---------------------------------------------------")
	log.Println(budgeResponse.Id)
	plan.Id = types.StringValue(budgeResponse.Id)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (r *reportResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	log.Print("hello report Read:)")
	// Get current state
	var state reportResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Print("state id::::::::::::::::::::::::::)")
	log.Print(state.Id.ValueString())
	// Get refreshed report value from DoiT
	report, err := r.client.GetReport(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Doit Console Attribution",
			"Could not read Doit Console Attribution ID "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}
	log.Print("response::::::::::::::::::::::::::)")
	log.Print(report)
	state.Id = types.StringValue(report.Id)
	log.Print("a")
	state.Description = types.StringValue(report.Description)
	state.Name = types.StringValue(report.Name)
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	log.Print("b")
	if report.Config.AdvancedAnalysis != nil {
		state.Config.AdvancedAnalysis.Forecast = types.BoolValue(report.Config.AdvancedAnalysis.Forecast)
		state.Config.AdvancedAnalysis.NotTrending = types.BoolValue(report.Config.AdvancedAnalysis.NotTrending)
		state.Config.AdvancedAnalysis.TrendingDown = types.BoolValue(report.Config.AdvancedAnalysis.TrendingDown)
		state.Config.AdvancedAnalysis.TrendingUp = types.BoolValue(report.Config.AdvancedAnalysis.TrendingUp)
	}
	log.Print("c")
	state.Config.Aggregation = types.StringValue(report.Config.Aggregation)
	state.Config.Currency = types.StringValue(report.Config.Currency)
	state.Config.DisplayValues = types.StringValue(report.Config.DisplayValues)
	state.Config.IncludePromotionalCredits = types.BoolValue(report.Config.IncludePromotionalCredits)
	state.Config.Layout = types.StringValue(report.Config.Layout)
	state.Config.TimeInterval = types.StringValue(report.Config.TimeInterval)
	log.Print("c")
	if report.Config.TimeRange != nil {
		state.Config.TimeRange.Amount = types.Int64Value(report.Config.TimeRange.Amount)
		state.Config.TimeRange.IncludeCurrent = types.BoolValue(report.Config.TimeRange.IncludeCurrent)
		state.Config.TimeRange.Mode = types.StringValue(report.Config.TimeRange.Mode)
		state.Config.TimeRange.Unit = types.StringValue(report.Config.TimeRange.Unit)
	}
	log.Print("d")
	state.Config.Metric.Type = types.StringValue(report.Config.Metric.Type)
	state.Config.Metric.Value = types.StringValue(report.Config.Metric.Value)
	log.Print("e")
	if report.Config.MetricFilter != nil {
		state.Config.MetricFilter.Operator = types.StringValue(report.Config.MetricFilter.Operator)
		state.Config.MetricFilter.Metric.Type = types.StringValue(report.Config.MetricFilter.Metric.Type)
		state.Config.MetricFilter.Metric.Value = types.StringValue(report.Config.MetricFilter.Metric.Value)
		state.Config.MetricFilter.Values = []types.Float64{}
		for _, value := range report.Config.MetricFilter.Values {
			state.Config.MetricFilter.Values = append(state.Config.MetricFilter.Values, types.Float64Value(value))
		}
	}
	log.Print("f")
	state.Config.Dimensions = []DimensionModel{}
	for _, dimension := range report.Config.Dimensions {
		state.Config.Dimensions = append(state.Config.Dimensions, DimensionModel{
			Id:   types.StringValue(dimension.Id),
			Type: types.StringValue(dimension.Type),
		})
	}
	log.Print("g")
	state.Config.Filters = []ExternalConfigFilterModel{}
	for _, filter := range report.Config.Filters {
		state.Config.Filters = append(state.Config.Filters, ExternalConfigFilterModel{
			Id:   types.StringValue(filter.Id),
			Type: types.StringValue(filter.Type),
		})
	}
	log.Print("h")
	state.Config.Group = []GroupModel{}
	for _, group := range report.Config.Group {
		state.Config.Group = append(state.Config.Group, GroupModel{
			Id:   types.StringValue(group.Id),
			Type: types.StringValue(group.Type),
		})
	}
	state.Config.Splits = []ExternalSplitModel{}
	for _, split := range report.Config.Splits {
		state.Config.Splits = append(state.Config.Splits, ExternalSplitModel{
			Id:   types.StringValue(split.Id),
			Type: types.StringValue(split.Type),
		})
	}
	log.Print("i")
	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	log.Print("state read::::::::::::::::::::::::::)")
	log.Print(state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *reportResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	log.Print("hello report Update:)")
	// Retrieve values from plan
	var plan reportResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state reportResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var report Report
	report.Id = state.Id.ValueString()
	report.Description = plan.Description.ValueString()
	report.Name = plan.Name.ValueString()
	if report.Config.AdvancedAnalysis != nil {
		report.Config.AdvancedAnalysis.Forecast = plan.Config.AdvancedAnalysis.Forecast.ValueBool()
		report.Config.AdvancedAnalysis.NotTrending = plan.Config.AdvancedAnalysis.NotTrending.ValueBool()
		report.Config.AdvancedAnalysis.TrendingDown = plan.Config.AdvancedAnalysis.TrendingDown.ValueBool()
		report.Config.AdvancedAnalysis.TrendingUp = plan.Config.AdvancedAnalysis.TrendingUp.ValueBool()
	}
	report.Config.Aggregation = plan.Config.Aggregation.ValueString()
	report.Config.Currency = plan.Config.Currency.ValueString()
	report.Config.Dimensions = []Dimension{}
	for _, dimension := range plan.Config.Dimensions {
		report.Config.Dimensions = append(report.Config.Dimensions, Dimension{
			Id:   dimension.Id.ValueString(),
			Type: dimension.Type.ValueString(),
		})
	}
	report.Config.DisplayValues = plan.Config.DisplayValues.ValueString()
	for _, filter := range plan.Config.Filters {
		var values []string
		for _, value := range filter.Values {
			values = append(values, value.ValueString())
		}
		report.Config.Filters = append(report.Config.Filters, ExternalConfigFilter{
			Id:      filter.Id.ValueString(),
			Type:    filter.Type.ValueString(),
			Inverse: filter.Inverse.ValueBool(),
			Values:  values,
		})
	}
	report.Config.Group = []Group{}
	for _, group := range plan.Config.Group {
		limit := Limit{
			Metric: ExternalMetric{
				Type:  group.Limit.Metric.Type.ValueString(),
				Value: group.Limit.Metric.Value.ValueString(),
			},
			Sort:  group.Limit.Sort.ValueString(),
			Value: group.Limit.Value.ValueInt64(),
		}
		report.Config.Group = append(report.Config.Group, Group{
			Id:    group.Id.ValueString(),
			Type:  group.Type.ValueString(),
			Limit: limit,
		})
	}
	report.Config.IncludePromotionalCredits = plan.Config.IncludePromotionalCredits.ValueBool()
	report.Config.Layout = plan.Config.Layout.ValueString()
	report.Config.Metric.Type = plan.Config.Metric.Type.ValueString()
	report.Config.Metric.Value = plan.Config.Metric.Value.ValueString()
	if report.Config.MetricFilter != nil {
		report.Config.MetricFilter.Operator = plan.Config.MetricFilter.Operator.ValueString()
		report.Config.MetricFilter.Metric.Type = plan.Config.MetricFilter.Metric.Type.ValueString()
		report.Config.MetricFilter.Metric.Value = plan.Config.MetricFilter.Metric.Value.ValueString()
		report.Config.MetricFilter.Values = []float64{}
		for _, value := range plan.Config.MetricFilter.Values {
			report.Config.MetricFilter.Values = append(report.Config.MetricFilter.Values, value.ValueFloat64())
		}
	}

	report.Config.Splits = []ExternalSplit{}
	for _, split := range plan.Config.Splits {
		origin := ExternalOrigin{
			Id:   split.Origin.Id.ValueString(),
			Type: split.Origin.Type.ValueString(),
		}
		var targets []ExternalSplitTarget
		for _, target := range split.Targets {
			target := ExternalSplitTarget{
				Id:    target.Id.ValueString(),
				Type:  target.Type.ValueString(),
				Value: target.Value.ValueFloat64(),
			}
			targets = append(targets, target)
		}
		report.Config.Splits = append(report.Config.Splits, ExternalSplit{
			Id:      split.Id.ValueString(),
			Type:    split.Type.ValueString(),
			Origin:  origin,
			Targets: targets,
		})
	}
	report.Config.TimeInterval = plan.Config.TimeInterval.ValueString()
	if report.Config.TimeRange != nil {
		report.Config.TimeRange.Amount = plan.Config.TimeRange.Amount.ValueInt64()
		report.Config.TimeRange.IncludeCurrent = plan.Config.TimeRange.IncludeCurrent.ValueBool()
		report.Config.TimeRange.Mode = plan.Config.TimeRange.Mode.ValueString()
		report.Config.TimeRange.Unit = plan.Config.TimeRange.Unit.ValueString()
	}
	// Update existing report
	_, err := r.client.UpdateReport(state.Id.ValueString(), report)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Report",
			"Could not update report, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated items from GetReport as UpdateReport items are not
	// populated.
	reportResponse, err := r.client.GetReport(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Report",
			"Could not read report ID "+plan.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	// Update resource state with updated items and timestamp
	plan.Id = types.StringValue(reportResponse.Id)
	plan.Id = types.StringValue(reportResponse.Id)
	plan.Description = types.StringValue(reportResponse.Description)
	plan.Name = types.StringValue(reportResponse.Name)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	if plan.Config.AdvancedAnalysis != nil {
		plan.Config.AdvancedAnalysis.Forecast = types.BoolValue(reportResponse.Config.AdvancedAnalysis.Forecast)
		plan.Config.AdvancedAnalysis.NotTrending = types.BoolValue(reportResponse.Config.AdvancedAnalysis.NotTrending)
		plan.Config.AdvancedAnalysis.TrendingDown = types.BoolValue(reportResponse.Config.AdvancedAnalysis.TrendingDown)
		plan.Config.AdvancedAnalysis.TrendingUp = types.BoolValue(reportResponse.Config.AdvancedAnalysis.TrendingUp)
	}
	plan.Config.Aggregation = types.StringValue(reportResponse.Config.Aggregation)
	plan.Config.Currency = types.StringValue(reportResponse.Config.Currency)
	plan.Config.DisplayValues = types.StringValue(reportResponse.Config.DisplayValues)
	plan.Config.IncludePromotionalCredits = types.BoolValue(reportResponse.Config.IncludePromotionalCredits)
	plan.Config.Layout = types.StringValue(reportResponse.Config.Layout)
	plan.Config.TimeInterval = types.StringValue(reportResponse.Config.TimeInterval)
	if plan.Config.TimeRange != nil {
		plan.Config.TimeRange.Amount = types.Int64Value(reportResponse.Config.TimeRange.Amount)
		plan.Config.TimeRange.IncludeCurrent = types.BoolValue(reportResponse.Config.TimeRange.IncludeCurrent)
		plan.Config.TimeRange.Mode = types.StringValue(reportResponse.Config.TimeRange.Mode)
		plan.Config.TimeRange.Unit = types.StringValue(reportResponse.Config.TimeRange.Unit)
	}
	plan.Config.Metric.Type = types.StringValue(reportResponse.Config.Metric.Type)
	plan.Config.Metric.Value = types.StringValue(reportResponse.Config.Metric.Value)
	if plan.Config.MetricFilter != nil {
		plan.Config.MetricFilter.Operator = types.StringValue(reportResponse.Config.MetricFilter.Operator)
		plan.Config.MetricFilter.Metric.Type = types.StringValue(reportResponse.Config.MetricFilter.Metric.Type)
		plan.Config.MetricFilter.Metric.Value = types.StringValue(reportResponse.Config.MetricFilter.Metric.Value)
		plan.Config.MetricFilter.Values = []types.Float64{}
		for _, value := range reportResponse.Config.MetricFilter.Values {
			plan.Config.MetricFilter.Values = append(plan.Config.MetricFilter.Values, types.Float64Value(value))
		}
	}

	plan.Config.Dimensions = []DimensionModel{}
	for _, dimension := range reportResponse.Config.Dimensions {
		plan.Config.Dimensions = append(plan.Config.Dimensions, DimensionModel{
			Id:   types.StringValue(dimension.Id),
			Type: types.StringValue(dimension.Type),
		})
	}
	if plan.Config.Filters != nil {
		plan.Config.Filters = []ExternalConfigFilterModel{}
		for _, filter := range reportResponse.Config.Filters {
			plan.Config.Filters = append(plan.Config.Filters, ExternalConfigFilterModel{
				Id:   types.StringValue(filter.Id),
				Type: types.StringValue(filter.Type),
			})
		}
	}
	if plan.Config.Group != nil {
		plan.Config.Group = []GroupModel{}
		for _, group := range reportResponse.Config.Group {
			plan.Config.Group = append(plan.Config.Group, GroupModel{
				Id:   types.StringValue(group.Id),
				Type: types.StringValue(group.Type),
			})
		}
	}
	if plan.Config.Splits != nil {
		plan.Config.Splits = []ExternalSplitModel{}
		for _, split := range reportResponse.Config.Splits {
			plan.Config.Splits = append(plan.Config.Splits, ExternalSplitModel{
				Id:   types.StringValue(split.Id),
				Type: types.StringValue(split.Type),
			})
		}
	}
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.

func (r *reportResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	log.Print("hello report Delete:)")
	// Retrieve values from state
	var state reportResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing report
	err := r.client.DeleteReport(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting DoiT Report",
			"Could not delete report, unexpected error: "+err.Error(),
		)
		return
	}
}
