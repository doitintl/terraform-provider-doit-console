package provider

// Attribution -
type Attribution struct {
	Id          string      `json:"id,omitempty"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Formula     string      `json:"formula,omitempty"`
	LastUpdated string      `json:"last_updated,omitempty"`
	Components  []Component `json:"components,omitempty"`
}

// Component -
type Component struct {
	TypeComponent string   `json:"type"`
	Key           string   `json:"key"`
	Values        []string `json:"values"`
}

// Attribution -
type AttributionGroup struct {
	Id           string   `json:"id,omitempty"`
	Name         string   `json:"name"`
	Description  string   `json:"description,omitempty"`
	LastUpdated  string   `json:"last_updated,omitempty"`
	Attributions []string `json:"attributions"`
}

// Attribution -
type AttributionGroupGet struct {
	Id           string        `json:"id,omitempty"`
	Name         string        `json:"name"`
	Description  string        `json:"description,omitempty"`
	LastUpdated  string        `json:"last_updated"`
	Attributions []Attribution `json:"attributions"`
}
