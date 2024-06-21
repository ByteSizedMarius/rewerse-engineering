package rewerse

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// GetDiscountsRaw returns the raw data from the API in a RawDiscounts struct.
// It contains links to the handout (Prospekt) as well as discount categories, which in turn contain the actual discounts.
func GetDiscountsRaw(marketID string) (rd RawDiscounts, err error) {
	// Build Request
	req, err := BuildCustomRequest(clientHost, "stationary-app-offers/"+marketID)
	if err != nil {
		return
	}

	// Send Request & Parse Response
	err = DoRequest(req, &rd)
	if err != nil {
		return
	}

	// Somewhat validate Result
	if len(rd.Data.Offers.Categories) == 0 {
		err = fmt.Errorf("no discounts found for market %s", marketID)
		return
	}

	return
}

const (
	expectedHeader   = "Produktdetails"
	expectedTitle    = "Hersteller: "
	expectedCellType = "DEFAULT"
)

// GetDiscounts returns the discounts from the API with less bloat.
// It contains the discount categories, which in turn contain the actual discounts.
// I removed various parameters I deemed unnecessary and parsed into different datatypes where it made sense to me.
// The struct Discounts also provides some helper methods.
func GetDiscounts(marketID string) (d Discounts, err error) {
	rd, err := GetDiscountsRaw(marketID)
	if err != nil {
		return
	}

	d.ValidUntil = time.Unix(int64(rd.Data.Offers.UntilDate)/1000, 0)
	var foundAnyManuf bool // detect if the manufacturer format changed
	for _, rawCat := range rd.Data.Offers.Categories {
		cat := DiscountCategory{
			ID:    rawCat.ID,
			Index: rawCat.Order,
			Title: rawCat.Title,
		}

		for _, rawOffer := range rawCat.Offers {
			if rawOffer.CellType != expectedCellType {
				//fmt.Println("skipping non-default cell type", rawOffer.CellType)
				continue
			}

			discount := Discount{
				Title:           rawOffer.Title,
				Subtitle:        rawOffer.Subtitle,
				Images:          rawOffer.Images,
				PriceRaw:        rawOffer.PriceData.Price,
				NutriScore:      rawOffer.Detail.NutriScore,
				ArticleNo:       rawOffer.RawValues.Nan,
				ProductCategory: rawOffer.RawValues.CategoryTitle,
			}

			// Parse Price
			p := strings.Replace(strings.TrimSpace(strings.Trim(rawOffer.PriceData.Price, "€")), ",", ".", 1)
			if p != "" && !strings.Contains(p, " ") {
				discount.Price, err = strconv.ParseFloat(p, 64)
				if err != nil {
					fmt.Printf("error parsing price: %v. please report.\n", err)
					return
				}
			}

			// Parse Manufacturer
			for _, content := range rawOffer.Detail.Contents {
				if content.Header == expectedHeader {
					for _, title := range content.Titles {
						if strings.HasPrefix(title, expectedTitle) {
							discount.Manufacturer = strings.TrimPrefix(title, expectedTitle)
							foundAnyManuf = true
							break
						}
					}
				}
				if discount.Manufacturer != "" {
					break
				}
			}

			cat.Offers = append(cat.Offers, discount)
		}

		d.Categories = append(d.Categories, cat)
	}

	if !foundAnyManuf {
		err = fmt.Errorf("could not find any manufacturer in the discounts. this suggests that the format changed")
	}

	return
}
