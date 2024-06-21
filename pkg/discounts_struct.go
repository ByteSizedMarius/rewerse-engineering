package rewerse

import (
	"fmt"
	"strings"
	"time"
)

// RawDiscounts is the struct that holds the raw discount data from the rewe API.
type RawDiscounts struct {
	Data struct {
		Offers struct {
			Handout struct {
				Width  int `json:"width"`
				Height int `json:"height"`
				Images []struct {
					Original  string `json:"original"`
					Thumbnail string `json:"thumbnail"`
				} `json:"images"`
			} `json:"handout"`
			Categories []struct {
				ID              string `json:"id"`
				Title           string `json:"title"`
				MoodURL         any    `json:"moodUrl"`
				Order           int    `json:"order"`
				BackgroundColor string `json:"backgroundColor"`
				ForegroundColor string `json:"foregroundColor"`
				Offers          []struct {
					CellType  string   `json:"cellType"`
					Overline  string   `json:"overline"`
					Title     string   `json:"title"`
					Subtitle  string   `json:"subtitle"`
					Images    []string `json:"images"`
					Biozid    bool     `json:"biozid"`
					PriceData struct {
						Price        string `json:"price"`
						RegularPrice string `json:"regularPrice"`
					} `json:"priceData"`
					LoyaltyBonus any `json:"loyaltyBonus"`
					Detail       struct {
						Sensational any    `json:"sensational"`
						PitchIn     string `json:"pitchIn"`
						Tags        []any  `json:"tags"`
						Contents    []struct {
							Header string   `json:"header"`
							Titles []string `json:"titles"`
						} `json:"contents"`
						Biocide    bool   `json:"biocide"`
						NutriScore string `json:"nutriScore"`
					} `json:"detail"`
					RawValues struct {
						CategoryTitle string  `json:"categoryTitle"`
						PriceAverage  float64 `json:"priceAverage"`
						FlyerPage     int     `json:"flyerPage"`
						Nan           string  `json:"nan"`
					} `json:"rawValues"`
				} `json:"offers"`
			} `json:"categories"`
			UntilDate       float64 `json:"untilDate"`
			HasOnlineOffers bool    `json:"hasOnlineOffers"`
		} `json:"offers"`
	} `json:"data"`
	Extensions struct {
		HTTP []struct {
			Path         []string `json:"path"`
			Message      string   `json:"message"`
			StatusCode   int      `json:"statusCode"`
			ResponseBody any      `json:"responseBody"`
		} `json:"http"`
	} `json:"extensions"`
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
	var s string
	for _, cat := range d.Categories {
		s += cat.Title + "\n"
		for _, offer := range cat.Offers {
			s += fmt.Sprintf("\t%s, %.2f€\n", offer.Title, offer.Price)
		}
	}
	return s
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
	PriceRaw        string
	Price           float64
	Manufacturer    string
	ArticleNo       string
	NutriScore      string
	ProductCategory string
}
