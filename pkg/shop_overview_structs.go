package rewerse

// This sucks and was very little fun to write.

import (
	"fmt"
	"strings"
)

const (
	noCategoriesStr = "No Categories found.\nMake sure the market has the service PICKUP (query market-details, services should contain type:PICKUP)"
)

//===================================================
//
// Shop Overview
//
//===================================================

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

//===================================================
//
// Category
//
//===================================================

type ShopCategory struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Slug            string `json:"slug"`
	ProductCount    int    `json:"productCount"`
	ImageURL        string `json:"imageUrl"`
	ChildCategories []ShopCategory
}

func (sc ShopCategory) String() string {
	return fmt.Sprintf(
		"%s%s\n"+ // ID
			"%s%s\n"+ // Name
			"%s%s\n"+ // Slug
			"%s%d\n"+ // ProductCount
			"%s%s\n%s", // ImageURL
		align("ID"),
		sc.ID,
		align("Name"),
		sc.Name,
		align("Slug"),
		sc.Slug,
		align("ProductCount"),
		sc.ProductCount,
		align("ImageURL"),
		sc.ImageURL,
		shopCatNestToString(len(sc.ChildCategories)),
	)
}

func (sc ShopCategory) StringAll() string {
	return generateStringAll(sc, 0)
}

//===================================================
//
// Internals
//
//===================================================

func shopCatNestToString(amnt int) string {
	return fmt.Sprintf("%s%d", align("Child-Amnt"), amnt)
}

func generateStringAll(sc ShopCategory, depth int) string {
	var s strings.Builder
	scs := sc.String()
	for i, scsl := range strings.Split(scs, "\n") {
		if i != 0 {
			s.WriteString("\n")
		}
		s.WriteString(alignL(scsl, depth))
	}

	children := sc.ChildCategories
	if len(children) != 0 {
		s.WriteString("\n" + alignL("Children:", depth+1))
		for _, child := range children {
			s.WriteString("\n" + alignL(child.Name, depth+2))
			childString := generateStringAll(child, depth)
			for _, line := range strings.Split(childString, "\n") {
				s.WriteString("\n" + alignL(line, depth+2))
			}
		}
	}

	return s.String()
}
