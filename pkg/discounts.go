package rewerse

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// priceRegex extracts German-format prices: "1,99", "12,99", "1.234,56"
var priceRegex = regexp.MustCompile(`(\d{1,3}(?:\.\d{3})*),(\d{2})`)

// GetDiscountsRaw returns the raw data from the API in a RawDiscounts struct.
// It contains links to the handout (Prospekt) as well as discount categories, which in turn contain the actual discounts.
func GetDiscountsRaw(marketID string) (rd RawDiscounts, err error) {
	req, err := BuildCustomRequest(clientHost, "stationary-offers/"+marketID)
	if err != nil {
		return
	}

	err = DoRequest(req, &rd)
	if err != nil {
		return
	}

	week := getActiveWeek(rd)
	if len(week.Categories) == 0 {
		err = fmt.Errorf("no discounts found for market %s", marketID)
		return
	}

	return
}

// getActiveWeek returns the week that should be displayed based on defaultWeek field
func getActiveWeek(rd RawDiscounts) RawOffersWeek {
	if rd.Data.Offers.DefaultWeek == "next" && rd.Data.Offers.Next.Available {
		return rd.Data.Offers.Next
	}
	return rd.Data.Offers.Current
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
func GetDiscounts(marketID string) (ds Discounts, err error) {
	rd, err := GetDiscountsRaw(marketID)
	if err != nil {
		return
	}

	week := getActiveWeek(rd)

	// Parse date from ISO format: "2026-01-10"
	ds.ValidUntil, err = time.Parse("2006-01-02", week.UntilDate)
	if err != nil {
		err = fmt.Errorf("error parsing date %s: %w", week.UntilDate, err)
		return
	}

	var foundAnyManuf bool
	for _, rawCat := range week.Categories {
		cat := DiscountCategory{
			ID:    rawCat.ID,
			Index: rawCat.Order,
			Title: rawCat.Title,
		}

		for _, rawOffer := range rawCat.Offers {
			if rawOffer.CellType != expectedCellType {
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

			// Parse Price: German format "1,99 €" or "1.234,56 €" -> float
			if match := priceRegex.FindStringSubmatch(rawOffer.PriceData.Price); match != nil {
				euros := strings.ReplaceAll(match[1], ".", "") // remove thousand separators
				cents := match[2]
				discount.Price, err = strconv.ParseFloat(euros+"."+cents, 64)
				if err != nil {
					discount.PriceParseFail = true
					err = nil // don't fail the entire call, just flag this discount
				}
			} else if rawOffer.PriceData.Price != "" {
				discount.PriceParseFail = true
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

		ds.Categories = append(ds.Categories, cat)
	}

	if !foundAnyManuf {
		err = fmt.Errorf("could not find any manufacturer in the discounts. this suggests that the format changed")
	}

	return
}
