package rewerse

import (
	"fmt"
	"strings"
)

// ProductResults is the flattened result from product search endpoints.
// Endpoint: GET /api/products
type ProductResults struct {
	Products   []Product
	Pagination struct {
		ObjectsPerPage int
		CurrentPage    int
		PageCount      int
		ObjectCount    int
	}
	SearchTerm struct {
		Original  string
		Corrected *string // non-nil if API corrected a typo
	}
}

func (pr ProductResults) String() string {
	var s strings.Builder
	s.WriteString("Products for query " + pr.SearchTerm.Original + "\n")
	s.WriteString(fmt.Sprintf("Page %d of %d\n", pr.Pagination.CurrentPage, pr.Pagination.PageCount))
	s.WriteString(sep("") + "\n")
	for _, product := range pr.Products {
		s.WriteString(product.Title + "\n")
	}
	return s.String()
}

// productSearchResponse is the raw API response structure (internal)
type productSearchResponse struct {
	Data struct {
		Products struct {
			Pagination struct {
				ObjectsPerPage int `json:"objectsPerPage"`
				CurrentPage    int `json:"currentPage"`
				PageCount      int `json:"pageCount"`
				ObjectCount    int `json:"objectCount"`
			} `json:"pagination"`
			Search struct {
				Term struct {
					Original  string  `json:"original"`
					Corrected *string `json:"corrected"`
				} `json:"term"`
			} `json:"search"`
			Products []Product `json:"products"`
		} `json:"products"`
	} `json:"data"`
}

type Product struct {
	ProductID    string `json:"productId"`
	Title        string `json:"title"`
	DepositLabel *string `json:"depositLabel"`
	ImageURL     string `json:"imageURL"`
	Attributes   struct {
		IsBulkyGood     bool `json:"isBulkyGood"`
		IsOrganic       bool `json:"isOrganic"`
		IsVegan         bool `json:"isVegan"`
		IsVegetarian    bool `json:"isVegetarian"`
		IsDairyFree     bool `json:"isDairyFree"`
		IsGlutenFree    bool `json:"isGlutenFree"`
		IsBiocide       bool `json:"isBiocide"`
		IsAgeRestricted *bool `json:"isAgeRestricted"`
		IsRegional      bool `json:"isRegional"`
		IsNew           bool `json:"isNew"`
		IsLowestPrice   bool `json:"isLowestPrice"`
	} `json:"attributes"`
	OrderLimit          int      `json:"orderLimit"`
	Categories          []string `json:"categories"`
	DetailsViewRequired bool     `json:"detailsViewRequired"`
	ArticleID           string   `json:"articleId"`
	Listing             struct {
		ListingID          string `json:"listingId"`
		ListingVersion     int    `json:"listingVersion"`
		CurrentRetailPrice int    `json:"currentRetailPrice"`
		TotalRefundPrice   *int   `json:"totalRefundPrice"`
		Grammage           string `json:"grammage"`
		Discount           any    `json:"discount"`
		LoyaltyBonus       any    `json:"loyaltyBonus"`
	} `json:"listing"`
	Advertisement any `json:"advertisement"`
}

// ProductDetail is the full product information from single-product lookup
// Endpoint: GET /api/products/{productId}
type ProductDetail struct {
	Product
	// RegulatedProductName is the legal product name: "Mandel-Cashew-Edamame-Mix geröstet, mit Chili-Würzung."
	RegulatedProductName string `json:"regulatedProductName"`
	// AdditionalImageURLs contains extra product images
	AdditionalImageURLs []string `json:"additionalImageURLs"`
	// Description is the product description
	Description string `json:"description"`
	// FeatureBenefit contains marketing benefits
	FeatureBenefit string `json:"featureBenefit"`
	// TradeItemMarketingMessage contains marketing text
	TradeItemMarketingMessage string `json:"tradeItemMarketingMessage"`
	// QSCertificationMark indicates QS certification
	QSCertificationMark bool `json:"qsCertificationMark"`
	// Brand is the manufacturer/brand name: "Bonduelle", "REWE Bio", "Hipp"
	Brand string `json:"brand"`
	// NutritionFacts contains nutritional information
	NutritionFacts []NutritionFact `json:"nutritionFacts"`
}

func (pd ProductDetail) String() string {
	var sb strings.Builder

	sb.WriteString(sep("Produkt"))
	sb.WriteByte('\n')
	sb.WriteString(align("ID"))
	sb.WriteString(pd.ProductID)
	sb.WriteByte('\n')
	sb.WriteString(align("Titel"))
	sb.WriteString(pd.Title)
	sb.WriteByte('\n')
	if pd.Brand != "" {
		sb.WriteString(align("Marke"))
		sb.WriteString(pd.Brand)
		sb.WriteByte('\n')
	}
	if pd.RegulatedProductName != "" {
		sb.WriteString(align("Bez."))
		sb.WriteString(pd.RegulatedProductName)
		sb.WriteByte('\n')
	}

	sb.WriteByte('\n')
	sb.WriteString(sep("Preis"))
	sb.WriteByte('\n')
	sb.WriteString(align("Preis"))
	sb.WriteString(fmt.Sprintf("%.2f €", float64(pd.Listing.CurrentRetailPrice)/100))
	sb.WriteByte('\n')
	sb.WriteString(align("Grammage"))
	sb.WriteString(pd.Listing.Grammage)
	sb.WriteByte('\n')
	if pd.DepositLabel != nil {
		sb.WriteString(align("Pfand"))
		sb.WriteString(*pd.DepositLabel)
		sb.WriteByte('\n')
	}

	var attrs []string
	if pd.Attributes.IsOrganic {
		attrs = append(attrs, "Bio")
	}
	if pd.Attributes.IsVegan {
		attrs = append(attrs, "Vegan")
	}
	if pd.Attributes.IsVegetarian {
		attrs = append(attrs, "Vegetarisch")
	}
	if pd.Attributes.IsDairyFree {
		attrs = append(attrs, "Laktosefrei")
	}
	if pd.Attributes.IsGlutenFree {
		attrs = append(attrs, "Glutenfrei")
	}
	if pd.Attributes.IsRegional {
		attrs = append(attrs, "Regional")
	}
	if pd.Attributes.IsNew {
		attrs = append(attrs, "Neu")
	}
	if pd.Attributes.IsLowestPrice {
		attrs = append(attrs, "Niedrigster Preis")
	}
	if pd.Attributes.IsAgeRestricted != nil && *pd.Attributes.IsAgeRestricted {
		attrs = append(attrs, "Altersbeschränkt")
	}

	if len(attrs) > 0 {
		sb.WriteByte('\n')
		sb.WriteString(sep("Eigenschaften"))
		sb.WriteByte('\n')
		sb.WriteString("   ")
		sb.WriteString(strings.Join(attrs, ", "))
		sb.WriteByte('\n')
	}

	if pd.Description != "" {
		sb.WriteByte('\n')
		sb.WriteString(sep("Beschreibung"))
		sb.WriteByte('\n')
		sb.WriteString("   ")
		sb.WriteString(pd.Description)
		sb.WriteByte('\n')
	}

	if len(pd.NutritionFacts) > 0 {
		sb.WriteByte('\n')
		sb.WriteString(sep("Nährwerte"))
		sb.WriteByte('\n')
		for _, nf := range pd.NutritionFacts {
			if nf.PreparationState != "" {
				sb.WriteString("   ")
				sb.WriteString(nf.PreparationState)
				sb.WriteString(":\n")
			}
			for _, ni := range nf.NutrientInformation {
				sb.WriteString(align(ni.NutrientType))
				sb.WriteString(fmt.Sprintf("%.1f %s", ni.QuantityContained.Value, ni.QuantityContained.UomShortText))
				sb.WriteByte('\n')
			}
		}
	}

	sb.WriteByte('\n')
	sb.WriteString(sep("IDs"))
	sb.WriteByte('\n')
	sb.WriteString(align("Artikel-ID"))
	sb.WriteString(pd.ArticleID)
	sb.WriteByte('\n')
	sb.WriteString(align("Listing-ID"))
	sb.WriteString(pd.Listing.ListingID)
	sb.WriteByte('\n')

	return sb.String()
}

// NutritionFact contains nutritional information for a preparation state
type NutritionFact struct {
	// PreparationState is: "Unzubereitet" or "Zubereitet"
	PreparationState string `json:"preparationState"`
	// NutrientInformation contains the individual nutrient values
	NutrientInformation []NutrientInfo `json:"nutrientInformation"`
}

// NutrientInfo is a single nutrient value
type NutrientInfo struct {
	// NutrientType is the nutrient name: "Energie", "Fett", "Kohlenhydrate", etc.
	NutrientType string `json:"nutrientType"`
	// MeasurementPrecision indicates accuracy: "ungefähr"
	MeasurementPrecision string `json:"measurementPrecision"`
	// QuantityContained has the value and unit
	QuantityContained struct {
		Value       float64 `json:"value"`
		UomShortText string  `json:"uomShortText"`
		UomLongText  string  `json:"uomLongText"`
	} `json:"quantityContained"`
}

type productDetailResponse struct {
	Data struct {
		Product []ProductDetail `json:"product"`
	} `json:"data"`
}

// ProductSuggestion is a search autocomplete suggestion
// Endpoint: GET /products/suggestion-search
type ProductSuggestion struct {
	// Title is the product name: "Hof Alpermühle Bio Eier 6 Stück"
	Title string `json:"title"`
	// ImageURL is a small product image (150x150)
	ImageURL string `json:"imageURL"`
	// RawValues contains product identifiers
	RawValues struct {
		// CategoryID is the category: "3523"
		CategoryID string `json:"categoryId"`
		// ProductID is the product ID: "7828199"
		ProductID string `json:"productId"`
		// ArticleID is the article ID: "T4KNENOV"
		ArticleID string `json:"articleId"`
		// Nan is the article number (German: Artikelnummer)
		Nan string `json:"nan"`
	} `json:"rawValues"`
}

type ProductSuggestions []ProductSuggestion

func (ps ProductSuggestions) String() string {
	s := "Product suggestions:\n"
	for _, p := range ps {
		s += fmt.Sprintf("  - %s (%s)\n", p.Title, p.RawValues.ProductID)
	}
	return s
}

// ProductRecommendations contains related product recommendations
// Endpoint: GET /api/products/recommendations?listingIds={listingId}
type productRecommendationsResponse struct {
	Data struct {
		ProductRecommendations struct {
			Recommendations []Product `json:"recommendations"`
		} `json:"productRecommendations"`
	} `json:"data"`
}
