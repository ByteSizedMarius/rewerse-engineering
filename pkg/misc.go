package rewerse

import (
	"fmt"
)

type recallsResponse struct {
	Data struct {
		ProductRecalls struct {
			Products Recalls `json:"products"`
		} `json:"productRecalls"`
	} `json:"data"`
}

// Recalls is the struct for Rewe Product-Recalls
type Recalls []Recall

func (rs Recalls) String() string {
	if len(rs) == 0 {
		return "Currently no recalls"
	}

	recalls := "Recalls:\n"
	for _, r := range rs {
		recalls += r.String() + "\n"
	}

	return recalls
}

// Recall is the struct for a single recall
type Recall struct {
	URL            string `json:"url"`
	SubjectProduct string `json:"subjectProduct"`
	SubjectReason  string `json:"subjectReason"`
}

func (r Recall) String() string {
	return fmt.Sprintf("%s\n%s\n%s\n", r.SubjectProduct, r.SubjectReason, r.URL)
}

// GetRecalls returns all currently ongoing recalls from Rewe
func GetRecalls() (r Recalls, err error) {
	req, err := BuildCustomRequest(clientHost, "products/recalls")
	if err != nil {
		return
	}

	var res recallsResponse
	err = DoRequest(req, &res)
	if err != nil {
		return
	}

	r = res.Data.ProductRecalls.Products
	return
}

type servicePortfolioResponse struct {
	Data struct {
		ServicePortfolio ServicePortfolio `json:"servicePortfolio"`
	} `json:"data"`
}

// ServicePortfolio contains available REWE services for a zip code
// Endpoint: GET /api/service-portfolio/{zipcode}
type ServicePortfolio struct {
	// CustomerZipCode is the queried zip code: "68199"
	CustomerZipCode string `json:"customerZipCode"`
	// DeliveryMarket contains the market that delivers to this zip code
	DeliveryMarket *struct {
		WWIdent string `json:"wwIdent"`
	} `json:"deliveryMarket"`
	// PickupMarkets contains markets offering pickup service
	PickupMarkets []PickupMarket `json:"pickupMarkets"`
}

func (sp ServicePortfolio) String() string {
	s := fmt.Sprintf("Service Portfolio for %s:\n", sp.CustomerZipCode)
	if sp.DeliveryMarket != nil {
		s += fmt.Sprintf("  Delivery available from market %s\n", sp.DeliveryMarket.WWIdent)
	} else {
		s += "  No delivery available\n"
	}
	s += fmt.Sprintf("  %d pickup markets available:\n", len(sp.PickupMarkets))
	for _, m := range sp.PickupMarkets {
		s += fmt.Sprintf("    - %s (%s): %s, %s %s\n", m.WWIdent, m.DisplayName, m.StreetWithHouseNumber, m.ZipCode, m.City)
	}
	return s
}

// PickupMarket is a market offering pickup service
type PickupMarket struct {
	// WWIdent is the market ID: "831002"
	WWIdent string `json:"wwIdent"`
	// DisplayName is the market type: "REWE Markt"
	DisplayName string `json:"displayName"`
	// CompanyName is the operating company: "REWE Hüseyin Özdemir oHG"
	CompanyName string `json:"companyName"`
	// IsPickupStation indicates if this is a pickup station (vs full store)
	IsPickupStation bool `json:"isPickupStation"`
	// SignedMapsUrl is the map URL path: "/api/markets/831002/map"
	SignedMapsUrl string `json:"signedMapsUrl"`
	// Latitude is the GPS coordinate: "49.45762"
	Latitude string `json:"latitude"`
	// Longitude is the GPS coordinate: "8.43085"
	Longitude string `json:"longitude"`
	// ZipCode is the market's zip code: "67065"
	ZipCode string `json:"zipCode"`
	// StreetWithHouseNumber is the address: "Wegelnburgstr. 33"
	StreetWithHouseNumber string `json:"streetWithHouseNumber"`
	// City is the city name: "Ludwigshafen / Mundenheim"
	City string `json:"city"`
	// PickupType is the pickup service type: "PICKUP_SERVICE"
	PickupType string `json:"pickupType"`
}

// GetServicePortfolio returns available REWE services for a zip code.
// Endpoint: GET /api/service-portfolio/{zipcode}
func GetServicePortfolio(zipcode string) (sp ServicePortfolio, err error) {
	req, err := BuildCustomRequest(clientHost, "service-portfolio/"+zipcode)
	if err != nil {
		return
	}

	var res servicePortfolioResponse
	err = DoRequest(req, &res)
	if err != nil {
		return
	}

	sp = res.Data.ServicePortfolio
	return
}
