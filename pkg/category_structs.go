package rewerse

import (
	"fmt"
	"strings"
)

type ProductSearchResults struct {
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
					Original  string `json:"original"`
					Corrected any    `json:"corrected"`
				} `json:"term"`
			} `json:"search"`
			Products []Product `json:"products"`
		} `json:"products"`
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

func (psr ProductSearchResults) String() string {
	var s strings.Builder
	s.WriteString("Products for query " + psr.Data.Products.Search.Term.Original + "\n")
	s.WriteString(fmt.Sprintf("Page %d of %d\n", psr.Data.Products.Pagination.CurrentPage, psr.Data.Products.Pagination.PageCount))
	s.WriteString(sep("") + "\n")
	for _, product := range psr.Data.Products.Products {
		s.WriteString(product.Title + "\n")
	}
	return s.String()
}

type Product struct {
	ProductID    string `json:"productId"`
	Title        string `json:"title"`
	DepositLabel any    `json:"depositLabel"`
	ImageURL     string `json:"imageURL"`
	Attributes   struct {
		IsBulkyGood     bool `json:"isBulkyGood"`
		IsOrganic       bool `json:"isOrganic"`
		IsVegan         bool `json:"isVegan"`
		IsVegetarian    bool `json:"isVegetarian"`
		IsDairyFree     bool `json:"isDairyFree"`
		IsGlutenFree    bool `json:"isGlutenFree"`
		IsBiocide       bool `json:"isBiocide"`
		IsAgeRestricted any  `json:"isAgeRestricted"`
		IsRegional      bool `json:"isRegional"`
		IsNew           bool `json:"isNew"`
	} `json:"attributes"`
	OrderLimit          int      `json:"orderLimit"`
	Categories          []string `json:"categories"`
	DetailsViewRequired bool     `json:"detailsViewRequired"`
	ArticleID           string   `json:"articleId"`
	Listing             struct {
		ListingID          string `json:"listingId"`
		ListingVersion     int    `json:"listingVersion"`
		CurrentRetailPrice int    `json:"currentRetailPrice"`
		TotalRefundPrice   any    `json:"totalRefundPrice"`
		Grammage           string `json:"grammage"`
		Discount           any    `json:"discount"`
		LoyaltyBonus       any    `json:"loyaltyBonus"`
	} `json:"listing"`
	Advertisement any `json:"advertisement"`
}
