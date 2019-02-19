package rideshare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Authorization struct {
	TokenType    string    `json:"token_type,omitempty"`
	AccessToken  string    `json:"access_token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	ExpireIn     int       `json:"expire_in,omitempty"`
	Scope        string    `json:"scope,omitempty"`
	ExpiresAt    time.Time `json:"-"`
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

type LyftRideEstimate struct {
	CostEstimates []CostEstimate `json:"cost_estimates,omitempty"`
}

type Lyft struct {
	auth     Authorization
	client   *http.Client
	username string
	passwd   string
}

func NewLyft(username, passwd string) *Lyft {
	return &Lyft{
		username: username,
		passwd:   passwd,
	}
}

func (api *Lyft) SetClient(c *http.Client) {
	api.client = c
}

func (api *Lyft) GetClient() *http.Client {
	if api.client == nil {
		return http.DefaultClient
	}
	return api.client
}

func (api *Lyft) getAuth() {
	var (
		auth   Authorization
		client = api.GetClient()
	)

	body := []byte(`{"grant_type": "client_credentials", "scope": "public"}`)
	req, err := http.NewRequest("POST", "https://api.lyft.com/oauth/token", bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}

	req.SetBasicAuth(api.username, api.passwd)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	json.NewDecoder(resp.Body).Decode(&auth)
	api.setAuth(auth)
}

func (api *Lyft) setAuth(auth Authorization) {
	api.auth = auth
	api.auth.ExpiresAt = time.Now().Add(time.Duration(auth.ExpireIn) * time.Second)
}

func (api *Lyft) stale() bool {
	return api.auth.AccessToken == "" ||
		(api.auth.RefreshToken != "" && api.auth.ExpiresAt.After(time.Now()))
}

const LyftEstimateURL = "https://api.lyft.com/v1/cost?start_lat=%f&start_lng=%f&end_lat=%f&end_lng=%f"

func (api *Lyft) GetEstimate(r Route) LyftRideEstimate {
	var (
		est    LyftRideEstimate
		client = api.GetClient()
	)

	if api.stale() {
		api.getAuth()
	}

	req, err := http.NewRequest("GET", fmt.Sprintf(LyftEstimateURL, r.SLat, r.SLng, r.ELat, r.ELng), nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+api.auth.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.NewDecoder(resp.Body).Decode(&est); err != nil {
		log.Fatal(err)
	}

	return est
}

type Uber struct {
	client      *http.Client
	serverToken string
}

func NewUber(serverToken string) *Uber {
	return &Uber{
		serverToken: serverToken,
	}
}

func (api *Uber) SetClient(c *http.Client) {
	api.client = c
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
type Route struct {
	SLat float64
	SLng float64
	ELat float64
	ELng float64
}

const uberPriceURL = "https://api.uber.com/v1.2/estimates/price?start_latitude=%f&start_longitude=%f&end_latitude=%f&end_longitude=%f"

// UberCostEstimate obtain the cost estimate on for all Uber rides for origin to
// destination.
func (api *Uber) UberCostEstimate(r Route) UberPrices {
	var (
		est    UberPrices
		client = api.GetClient()
	)

	req, err := http.NewRequest("GET", fmt.Sprintf(uberPriceURL, r.SLat, r.SLng, r.ELat, r.ELng), nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "Token "+api.serverToken)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.NewDecoder(resp.Body).Decode(&est); err != nil {
		log.Fatal(err)
	}

	return est
}

func (api *Uber) GetClient() *http.Client {
	if api.client == nil {
		return http.DefaultClient
	}
	return api.client
}
