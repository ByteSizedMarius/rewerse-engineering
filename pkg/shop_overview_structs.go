package rewerse

import (
	"fmt"
	"strings"
)

const (
	noCategoriesStr  = "No Categories found.\nMake sure the market supports the service type you're using (PICKUP or DELIVERY)."
	maxCategoryDepth = 10
)

type ShopOverview struct {
	ProductRecalls    Recalls        `json:"productRecalls"`
	ProductCategories []ShopCategory `json:"productCategories"`
}

func (so ShopOverview) String() string {
	if len(so.ProductCategories) == 0 {
		return noCategoriesStr
	}

	var s strings.Builder
	for _, pc := range so.ProductCategories {
		s.WriteString("\n" + sep(pc.Name) + "\n")
		s.WriteString(pc.StringAll())
	}

	return s.String()
}

func (so ShopOverview) GetName() string {
	return "Shop Overview"
}

func (so ShopOverview) StringAll() string {
	if len(so.ProductCategories) == 0 {
		return noCategoriesStr
	}

	var s strings.Builder
	for _, pc := range so.ProductCategories {
		s.WriteString("\n" + sep(pc.Name) + "\n")
		s.WriteString(generateStringAll(pc, 0))
	}

	return s.String()
}

type ShopCategory struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Slug            string `json:"slug"`
	ProductCount    int    `json:"productCount"`
	ImageURL        string `json:"imageUrl"`
	ChildCategories []ShopCategory `json:"childCategories"`
}

func (sc ShopCategory) String() string {
	return fmt.Sprintf("%s (%s) - %d products", sc.Name, sc.Slug, sc.ProductCount)
}

func (sc ShopCategory) StringAll() string {
	return generateStringAll(sc, 0)
}

func generateStringAll(sc ShopCategory, depth int) string {
	if depth > maxCategoryDepth {
		return alignL("... (max depth reached)", depth)
	}

	var s strings.Builder
	s.WriteString(alignL(sc.String(), depth))

	for _, child := range sc.ChildCategories {
		s.WriteString("\n")
		s.WriteString(generateStringAll(child, depth+1))
	}

	return s.String()
}
