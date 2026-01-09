package rewerse

import (
	"fmt"
	"strings"
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

// RecipeHub is the struct for the Data returned by the Rewe Recipe-Page
type RecipeHub struct {
	RecipeOfTheDay Recipe   `json:"recipeOfTheDay"`
	PopularRecipes []Recipe `json:"popularRecipes"`
	Categories     []struct {
		Type        string `json:"type"`
		Title       string `json:"title"`
		SearchQuery string `json:"searchQuery"`
	} `json:"categories"`
}

func (rh RecipeHub) String() string {
	var sb strings.Builder

	sb.WriteString(sep("Recipe of the Day"))
	sb.WriteByte('\n')
	sb.WriteString(rh.RecipeOfTheDay.String())
	sb.WriteByte('\n')

	sb.WriteString(sep(fmt.Sprintf("Popular Recipes (%d)", len(rh.PopularRecipes))))
	sb.WriteByte('\n')
	for _, r := range rh.PopularRecipes {
		sb.WriteString(r.String())
	}

	sb.WriteByte('\n')
	sb.WriteString(sep(fmt.Sprintf("Categories (%d)", len(rh.Categories))))
	sb.WriteByte('\n')
	for _, c := range rh.Categories {
		sb.WriteString("   ")
		sb.WriteString(c.Title)
		sb.WriteByte('\n')
	}

	return sb.String()
}

// Recipe is the struct for a single recipe
type Recipe struct {
	ID                    string `json:"id"`
	Title                 string `json:"title"`
	DetailURL             string `json:"detailUrl"`
	ImageURL              string `json:"imageUrl"`
	Duration              string `json:"duration"`
	DifficultyLevel       int    `json:"difficultyLevel"`
	DifficultyDescription string `json:"difficultyDescription"`
}

func (r Recipe) String() string {
	return fmt.Sprintf("   %s (%s, %s)\n   %s\n", r.Title, r.Duration, r.DifficultyDescription, r.DetailURL)
}

// GetRecipeHub returns the Data from the RecipeHub
func GetRecipeHub() (r RecipeHub, err error) {
	req, err := BuildCustomRequest(apiHost, "v3/recipe-hub")
	if err != nil {
		return
	}

	err = DoRequest(req, &r)
	if err != nil {
		return
	}

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
