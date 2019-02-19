package rideshare

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type Authorization struct {
	TokenType   string `json:"token_type,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
	ExpireIn    int    `json:"expire_in,omitempty"`
	Scope       string `json:"scope,omitempty"`
}

type CostEstimate struct {
	Currency                   string      `json:"currency,omitempty"`
	RideType                   string      `json:"ride_type,omitempty"`
	DisplayName                string      `json:"display_name,omitempty"`
	PrimetimePercentage        string      `json:"primetime_percentage,omitempty"`
	PrimetimeConfirmationToken interface{} `json:"primetime_confirmation_token,omitempty"`
	CostToken                  interface{} `json:"cost_token,omitempty"`
	PriceQuoteID               string      `json:"price_quote_id,omitempty"`
	PriceGroupID               string      `json:"price_group_id,omitempty"`
	IsScheduledRide            bool        `json:"is_scheduled_ride,omitempty"`
	IsValidEstimate            bool        `json:"is_valid_estimate,omitempty"`
	EstimatedDurationSeconds   int         `json:"estimated_duration_seconds,omitempty"`
	EstimatedDistanceMiles     float64     `json:"estimated_distance_miles,omitempty"`
	EstimatedCostCentsMin      int         `json:"estimated_cost_cents_min,omitempty"`
	EstimatedCostCentsMax      int         `json:"estimated_cost_cents_max,omitempty"`
	CanRequestRide             bool        `json:"can_request_ride,omitempty"`
}
type LyftAvailabilityRideEstimate struct {
	CostEstimates []CostEstimate `json:"cost_estimates,omitempty"`
}

// LyftCostEstimate obtain the cost estimate on for all lyft rides for origin to
// destination.
func LyftCostEstimate(bearerToken string) LyftAvailabilityRideEstimate {
	var est LyftAvailabilityRideEstimate
	client := http.DefaultClient
	req, err := http.NewRequest("GET", "https://api.lyft.com/v1/cost?start_lat=37.7763&start_lng=-122.3918&end_lat=37.7972&end_lng=-122.4533", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", bearerToken)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.NewDecoder(resp.Body).Decode(&est); err != nil {
		log.Fatal(err)
	}

	return est
}

// LyftAuth sends a request to the Lyft authorization endpoint and supplies the
// client id and client secret and obtains the access token.
func LyftAuth(client *http.Client, username, passwd string) Authorization {
	var auth Authorization
	if client == nil {
		client = http.DefaultClient
	}

	body := []byte(`{"grant_type": "client_credentials", "scope": "public"}`)
	req, err := http.NewRequest("POST", "https://api.lyft.com/oauth/token", bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}

	req.SetBasicAuth(username, passwd)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	json.NewDecoder(resp.Body).Decode(&auth)
	return auth
}

type UberPrice struct {
	LocalizedDisplayName string  `json:"localized_display_name,omitempty"`
	Distance             float64 `json:"distance,omitempty"`
	DisplayName          string  `json:"display_name,omitempty"`
	ProductID            string  `json:"product_id,omitempty"`
	HighEstimate         float64 `json:"high_estimate,omitempty"`
	LowEstimate          float64 `json:"low_estimate,omitempty"`
	Duration             int     `json:"duration,omitempty"`
	Estimate             string  `json:"estimate,omitempty"`
	CurrencyCode         string  `json:"currency_code,omitempty"`
}
type UberPrices struct {
	Prices []UberPrice `json:"prices,omitempty"`
}

// UberCostEstimate obtain the cost estimate on for all Uber rides for origin to
// destination.
func UberCostEstimate(token string) UberPrices {
	var est UberPrices
	client := http.DefaultClient
	req, err := http.NewRequest("GET", "https://api.uber.com/v1.2/estimates/price?start_latitude=37.7763&start_longitude=-122.3918&end_latitude=37.7972&end_longitude=-122.4533", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", token)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.NewDecoder(resp.Body).Decode(&est); err != nil {
		log.Fatal(err)
	}

	return est
}
