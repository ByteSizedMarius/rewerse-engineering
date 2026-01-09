package rewerse

import (
	"fmt"
	"strings"
	"time"
)

// RawDiscounts is the struct that holds the raw discount data from the rewe API.
// Endpoint: GET /api/stationary-offers/{marketId}
type RawDiscounts struct {
	Data struct {
		Offers struct {
			// DefaultWeek indicates which week to show by default: "current" or "next"
			DefaultWeek string `json:"defaultWeek"`
			// NextWeekAvailableFrom indicates when next week's offers become visible: "saturday"
			NextWeekAvailableFrom string `json:"nextWeekAvailableFrom"`
			// Current contains this week's offers
			Current RawOffersWeek `json:"current"`
			// Next contains next week's offers (if available)
			Next RawOffersWeek `json:"next"`
		} `json:"offers"`
	} `json:"data"`
}

// RawOffersWeek contains offers for a specific week
type RawOffersWeek struct {
	// Available indicates if offers are available for this week
	Available bool `json:"available"`
	// FromDate is the start date in ISO format: "2026-01-05"
	FromDate string `json:"fromDate"`
	// UntilDate is the end date in ISO format: "2026-01-10"
	UntilDate string `json:"untilDate"`
	// HasOnlineOffers indicates if online ordering is available
	HasOnlineOffers bool `json:"hasOnlineOffers"`
	// Handout contains the printed flyer/prospekt images
	Handout struct {
		Images []struct {
			// Original is the full-size image URL
			Original string `json:"original"`
			// Thumbnail is a smaller preview
			Thumbnail string `json:"thumbnail"`
		} `json:"images"`
	} `json:"handout"`
	// Categories contains grouped offers
	Categories []RawOfferCategory `json:"categories"`
}

// RawOfferCategory is a grouping of offers as defined by REWE
type RawOfferCategory struct {
	// ID is the category identifier: "markt-topangebote", "suesses-und-salziges"
	ID string `json:"id"`
	// Title is the display name: "Top-Angebote in deinem Markt"
	Title string `json:"title"`
	// MoodURL is an optional header image for the category
	MoodURL any `json:"moodUrl"`
	// Order is the display order (lower = earlier)
	Order int `json:"order"`
	// BackgroundColor is the hex color for the category header: "#EDF9FA"
	BackgroundColor string `json:"backgroundColor"`
	// ForegroundColor is the hex text color: "#2E7B85"
	ForegroundColor string `json:"foregroundColor"`
	// Offers contains the actual discounted products
	Offers []RawOffer `json:"offers"`
}

// RawOffer is a single discounted product
type RawOffer struct {
	// CellType is the display type: "DEFAULT" for regular offers
	CellType string `json:"cellType"`
	// Overline is optional text above the title
	Overline string `json:"overline"`
	// Title is the product name: "Haribo Goldbären oder Color-Rado"
	Title string `json:"title"`
	// Subtitle contains weight/quantity info: "je 175-g-Btl. (1 kg = 4.40)"
	Subtitle string `json:"subtitle"`
	// Images contains product image URLs
	Images []string `json:"images"`
	// Biozid indicates if the product is a biocide
	Biozid bool `json:"biozid"`
	// PriceData contains pricing information
	PriceData struct {
		// Price is the discounted price: "0,77 €"
		Price string `json:"price"`
		// RegularPrice is the original price or a label like "Knaller"
		RegularPrice string `json:"regularPrice"`
	} `json:"priceData"`
	// LoyaltyBonus contains bonus points info if applicable
	LoyaltyBonus any `json:"loyaltyBonus"`
	// Stock contains availability info
	Stock any `json:"stock"`
	// Detail contains additional product information
	Detail struct {
		Sensational any    `json:"sensational"`
		PitchIn     string `json:"pitchIn"`
		Tags        []any  `json:"tags"`
		// Contents contains structured product details
		Contents []struct {
			// Header is the section name: "Produktdetails", "Hinweise"
			Header string `json:"header"`
			// Titles contains the detail lines: "Art.-Nr.: 7772669", "Hersteller: HARIBO"
			Titles []string `json:"titles"`
		} `json:"contents"`
		Biocide bool `json:"biocide"`
		// NutriScore is the nutrition score letter if available: "A", "B", etc.
		NutriScore string `json:"nutriScore"`
	} `json:"detail"`
	// RawValues contains internal tracking data
	RawValues struct {
		// CategoryTitle is the product category slug: "suesses-und-salziges"
		CategoryTitle string `json:"categoryTitle"`
		PriceAverage  float64 `json:"priceAverage"`
		FlyerPage     int     `json:"flyerPage"`
		// Nan is the article number (German: Artikelnummer): "7772669"
		Nan string `json:"nan"`
	} `json:"rawValues"`
}

// Discounts is the struct that holds cleaned up discount data.
type Discounts struct {
	Categories []DiscountCategory
	ValidUntil time.Time
}

// GroupByProductCategory groups the discounts by their product category and returns the results.
func (d Discounts) GroupByProductCategory() Discounts {
	catMap := make(map[string][]Discount)
	for _, cat := range d.Categories {
		for _, offer := range cat.Offers {
			if _, ok := catMap[offer.ProductCategory]; !ok {
				catMap[offer.ProductCategory] = []Discount{offer}
			} else {
				catMap[offer.ProductCategory] = append(catMap[offer.ProductCategory], offer)
			}
		}
	}

	i := 0
	d.Categories = make([]DiscountCategory, len(catMap))
	for cat, offers := range catMap {
		d.Categories[i] = DiscountCategory{
			ID:     strings.ToLower(strings.ReplaceAll(cat, " ", "-")),
			Title:  cat,
			Offers: offers,
		}
		i++
	}

	return d
}

func (d Discounts) String() string {
	var sb strings.Builder
	for _, cat := range d.Categories {
		sb.WriteString(cat.Title)
		sb.WriteByte('\n')
		for _, offer := range cat.Offers {
			sb.WriteString(fmt.Sprintf("\t%s, %.2f€\n", offer.Title, offer.Price))
		}
	}
	return sb.String()
}

// DiscountCategory is the category defined by rewe for the presentation of the discounts.
// It contains an index for sorting the categories in their intended order and the actual discounts.
// Calling GroupByProductCategory reorders the discounts by their product category (z. B. "Nahrungsmittel").
type DiscountCategory struct {
	ID     string
	Title  string
	Index  int
	Offers []Discount
}

// Discount is the actual discount with some of the information provided by rewe.
type Discount struct {
	Title           string
	Subtitle        string
	Images          []string
	PriceRaw string
	// Price is the parsed price in euros. Check PriceParseFail before using -
	// if true, Price is 0.0 due to parse failure, not because item is free.
	Price          float64
	PriceParseFail bool
	Manufacturer    string
	ArticleNo       string
	NutriScore      string
	ProductCategory string
}
