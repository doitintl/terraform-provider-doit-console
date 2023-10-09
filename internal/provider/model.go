package provider

// Attribution -
type Attribution struct {
	Id          string      `json:"id,omitempty"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Formula     string      `json:"formula"`
	LastUpdated string      `json:"last_updated"`
	Components  []Component `json:"components"`
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
	NullFallBack string   `json:"nullfallback"`
	LastUpdated  string   `json:"last_updated"`
	Attributions []string `json:"attributions"`
}
