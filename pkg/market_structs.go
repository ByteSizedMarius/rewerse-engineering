package rewerse

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Markets is a struct that holds a list of Market structs
// It is returned by the MarketSearch
type Markets []Market

func (ms Markets) String() string {
	s := "ID      Standort\n"
	for _, m := range ms {
		s += m.String() + "\n"
	}
	return s
}

// Market represents a single search result
type Market struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	AddressLine1 string `json:"addressLine1"`
	RawValues    struct {
		PostalCode string `json:"postalCode"`
		City       string `json:"city"`
	} `json:"rawValues"`
}

func (m Market) String() string {
	return fmt.Sprintf("%s: %s, %s, %s %s", m.ID, m.Name, m.AddressLine1, m.RawValues.PostalCode, m.RawValues.City)
}

// MarketDetails is a struct that holds detailed information about a single market
// It is returned by the GetMarketDetails function
type MarketDetails struct {
	MarketItem struct {
		ID                string `json:"id"`
		Name              string `json:"name"`
		TypeID            string `json:"typeId"`
		AddressLine1      string `json:"addressLine1"`
		AddressLine2      string `json:"addressLine2"`
		OpeningInfo       string `json:"openingInfo"`
		OpeningInfoPrefix string `json:"openingInfoPrefix"`
		OpeningType       string `json:"openingType"`
		FeatureTypes      []any  `json:"featureTypes"`
		Location          struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"location"`
		RawValues struct {
			Attributes []any  `json:"attributes"`
			PostalCode string `json:"postalCode"`
			City       string `json:"city"`
		} `json:"rawValues"`
	} `json:"marketItem"`
	Phone        string `json:"phone"`
	RatingURL    string `json:"ratingUrl"`
	IsLSFK       bool   `json:"isLSFK"`
	OpeningTimes []struct {
		Days  string `json:"days"`
		Hours string `json:"hours"`
	} `json:"openingTimes"`
	SpecialOpeningTimes []any `json:"specialOpeningTimes"`
	Services            []any `json:"services"`
	FeatureCategories   []any `json:"featureCategories"`
	Actions             []struct {
		Type  string `json:"type"`
		Title string `json:"title"`
	} `json:"actions"`
}

func (md MarketDetails) String() string {
	raw, _ := json.Marshal(md.MarketItem.RawValues)
	ot, _ := json.Marshal(md.OpeningTimes)
	return fmt.Sprintf(
		"%s\n"+ // titel:general
			"%s%s\n"+ // id
			"%s%s\n"+ // name
			"%s%s\n"+ // type
			"%s%s, %s (%.4f, %.4f)\n"+ // standort
			"%s%s\n"+ // telefon
			"%s%s\n"+ // raw values
			"\n%s\n"+ // titel:öffnungszeiten
			"%s%s\n"+
			"%s  %s (Typ: %s)\n"+ // opening info
			"%s %s\n"+
			"%s%s\n"+ // besonderheiten
			"\n%s\n"+ // titel:verschiedenes
			"%s%s %s\n"+ // features (??)
			"%s%s\n"+ // services (??)
			"%s%s\n"+ // actions (??)
			"%s%s\n"+ // rating url
			"%s%t\n", // islsfk (??)
		sep("Allgemein"),
		align("ID"),
		md.MarketItem.ID,
		align("Name"),
		md.MarketItem.Name,
		align("Type-ID"),
		md.MarketItem.TypeID,
		align("Standort"),
		md.MarketItem.AddressLine1,
		md.MarketItem.AddressLine2,
		md.MarketItem.Location.Latitude,
		md.MarketItem.Location.Longitude,
		align("Tel-Nr"),
		md.Phone,
		align("Rohdaten"),
		raw,
		sep("Öffnungszeiten"),
		align("Aktell"),
		md.MarketItem.OpeningInfoPrefix,
		align("Nächste Änderung"),
		md.MarketItem.OpeningInfo,
		md.MarketItem.OpeningType,
		align("Wochenübersicht"),
		ot,
		align("Besonderheiten"),
		md.SpecialOpeningTimes,
		sep("Verschiedenes"),
		align("Features"),
		md.MarketItem.FeatureTypes,
		md.FeatureCategories,
		align("Services"),
		md.Services,
		align("Aktionen"),
		md.Actions,
		align("Bewertungs-URL"),
		md.RatingURL,
		align("Lieferservice"),
		md.IsLSFK,
	)
}

func sep(base string) string {
	return "---" + base + strings.Repeat("-", 60-len(base))
}

func align(s string) string {
	return "   " + s + ":" + strings.Repeat(" ", 22-len(s))
}
