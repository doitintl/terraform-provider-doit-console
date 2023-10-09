package provider

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// CreateAttributionGroup - Create new attributionGroup
func (c *ClientTest) CreateAttributionGroup(attributionGroup AttributionGroup) (*AttributionGroup, error) {
	log.Println("CreateAttributionGroup************************")
	log.Println(attributionGroup)
	rb, err := json.Marshal(attributionGroup)
	if err != nil {
		return nil, err
	}

	log.Println("rb")
	log.Println(strings.NewReader(string(rb)))

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/analytics/v1/attributiongroups/?customerContext=%s", c.HostURL, c.Auth.CustomerContext), strings.NewReader(string(rb)))
	log.Println("URL----------------")
	log.Println(req.URL)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		log.Println("ERROR REQUEST----------------")
		log.Println(err)
		return nil, err
	}

	attributionGroupResponse := AttributionGroup{}
	err = json.Unmarshal(body, &attributionGroupResponse)
	if err != nil {
		log.Println("ERROR UNMARSHALL----------------")
		log.Println(err)
		return nil, err
	}
	log.Println("AttributionGroup response----------------")
	log.Println(attributionGroupResponse)
	return &attributionGroupResponse, nil
}

// UpdateAttributionGroup - Updates an attributionGroup
func (c *ClientTest) UpdateAttributionGroup(attributionGroupID string, attributionGroup AttributionGroup) (*AttributionGroup, error) {
	rb, err := json.Marshal(attributionGroup)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/analytics/v1/attributiongroups/%s/?customerContext=%s", c.HostURL, attributionGroupID, c.Auth.CustomerContext), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}
	log.Println("Update URL----------------")
	log.Println(req.URL)
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	attributionGroupResponse := AttributionGroup{}
	err = json.Unmarshal(body, &attributionGroupResponse)
	if err != nil {
		return nil, err
	}
	log.Println("AttributionGroup response----------------")
	log.Println(attributionGroupResponse)
	return &attributionGroupResponse, nil
}

func (c *ClientTest) DeleteAttributionGroup(attributionGroupID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/analytics/v1/attributiongroups/%s/?customerContext=%s", c.HostURL, attributionGroupID, c.Auth.CustomerContext), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

// GetAttributionGroup - Returns a specifc attribution
func (c *ClientTest) GetAttributionGroup(orderID string) (*AttributionGroup, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/analytics/v1/attributiongroups/%s/?customerContext=%s", c.HostURL, orderID, c.Auth.CustomerContext), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	attributionGroup := AttributionGroup{}
	err = json.Unmarshal(body, &attributionGroup)
	if err != nil {
		return nil, err
	}

	return &attributionGroup, nil
}
