package rewerse

import (
	"fmt"
	"strings"
)

type marketSearchResponse struct {
	Data struct {
		MarketSearch struct {
			Markets Markets `json:"markets"`
		} `json:"marketSearch"`
	} `json:"data"`
}

type marketDetailsResponse struct {
	Data struct {
		Market  Market        `json:"market"`
		Content MarketContent `json:"content"`
	} `json:"data"`
}

type Markets []Market

func (ms Markets) String() string {
	var sb strings.Builder
	sb.WriteString("ID      Location\n")
	for _, m := range ms {
		sb.WriteString(m.String())
		sb.WriteByte('\n')
	}
	return sb.String()
}

type Market struct {
	// WWIdent is the market ID: "840174", "831002"
	WWIdent string `json:"wwIdent"`
	// Name is the store type: "REWE Markt", "REWE Center"
	Name string `json:"name"`
	// CompanyName is the operating company: "REWE Markt Vuthaj oHG", "REWE Thomas Viering oHG"
	CompanyName string `json:"companyName"`
	// Phone is the store phone number: "0621-8414498"
	Phone string `json:"phone"`
	// TypeID is the market type code: "REWE", "CENTER"
	TypeID string `json:"typeId"`
	// Street is the street address: "Rheingoldstr. 18-20", "Lindenhofstr. 91"
	Street string `json:"street"`
	// ZipCode is the postal code: "68199", "68161"
	ZipCode string `json:"zipCode"`
	// City is the city name: "Mannheim / Neckarau", "Mannheim"
	City string `json:"city"`
	// Location contains GPS coordinates
	Location struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
	// OpeningStatus contains current open/closed state
	OpeningStatus struct {
		// OpenState is the status: "OPEN", "CLOSED"
		OpenState string `json:"openState"`
		// InfoText is the closing time: "bis 22:00 Uhr"
		InfoText string `json:"infoText"`
		// StatusText is human-readable status: "Geöffnet", "Geschlossen"
		StatusText string `json:"statusText"`
	} `json:"openingStatus"`
	// OpeningInfo contains weekly opening hours
	OpeningInfo []struct {
		// Days is the day range: "Mo - Sa"
		Days string `json:"days"`
		// Hours is the time range: "07:00 - 22:00"
		Hours string `json:"hours"`
	} `json:"openingInfo"`
	// Category contains market type classification
	Category struct {
		// Template is the layout template code: "RN", "RC"
		Template string `json:"template"`
		// MarketTypeDisplayName is the display name: "REWE Markt", "REWE Center"
		MarketTypeDisplayName string `json:"marketTypeDisplayName"`
		// IsWarehouse indicates if this is a warehouse store
		IsWarehouse bool `json:"isWarehouse"`
		// Type is the category type: "MARKET"
		Type string `json:"type"`
	} `json:"category"`
	// Distance is the distance from search location in km (null when not searching by location)
	Distance *float64 `json:"distance,omitempty"`
	// ServiceFlags contains available services
	ServiceFlags struct {
		// HasPickup indicates if pickup service is available
		HasPickup bool `json:"hasPickup"`
	} `json:"serviceFlags"`
}

func (m Market) String() string {
	return fmt.Sprintf("%s: %s, %s, %s %s", m.WWIdent, m.Name, m.Street, m.ZipCode, m.City)
}

// MarketContent contains additional market details from the details endpoint
type MarketContent struct {
	// Components contains UI components (usually null)
	Components any `json:"components"`
	// MarketData contains market-specific info
	MarketData struct {
		// MarketName is the branded name: "REWE Vuthaj"
		MarketName string `json:"marketName"`
		// MarketManager is the manager's name (often empty)
		MarketManager string `json:"marketManager"`
		// ImageMarket is the store image URL (often empty)
		ImageMarket string `json:"imageMarket"`
		// ImageManager is the manager's photo URL (often empty)
		ImageManager string `json:"imageManager"`
		// MarketLink is the web URL (often empty)
		MarketLink string `json:"marketLink"`
	} `json:"marketData"`
	// Services contains available store services
	Services struct {
		// Fixed are standard services: Fleischtheke, Fischtheke, Parkplätze, etc.
		Fixed []MarketService `json:"fixed"`
		// Editable are custom services set by the store
		Editable []MarketService `json:"editable"`
	} `json:"services"`
}

// MarketService represents a service available at the market
type MarketService struct {
	// Text is the service name: "Fleisch- und Wursttheke", "Parkplätze", "WLAN"
	Text string `json:"text"`
	// Icon is the icon identifier: "sausagemeat", "parking", "wlan", "default"
	Icon string `json:"icon"`
	// Tooltip is additional info: "Große Auswahl deutsche Weine"
	Tooltip string `json:"tooltip"`
	// Active indicates if this service is available at the store
	Active bool `json:"active"`
	// IconURL is the full icon image URL
	IconURL string `json:"iconUrl"`
}

type MarketDetails struct {
	Market  Market
	Content MarketContent
}

func (md MarketDetails) String() string {
	var sb strings.Builder

	// General info
	sb.WriteString(sep("Allgemein"))
	sb.WriteByte('\n')
	sb.WriteString(align("ID"))
	sb.WriteString(md.Market.WWIdent)
	sb.WriteByte('\n')
	sb.WriteString(align("Name"))
	sb.WriteString(md.Market.Name)
	sb.WriteString(" (")
	sb.WriteString(md.Market.CompanyName)
	sb.WriteString(")\n")
	sb.WriteString(align("Type"))
	sb.WriteString(md.Market.Category.MarketTypeDisplayName)
	sb.WriteByte('\n')
	sb.WriteString(align("Standort"))
	sb.WriteString(fmt.Sprintf("%s, %s %s (%.4f, %.4f)", md.Market.Street, md.Market.ZipCode, md.Market.City, md.Market.Location.Latitude, md.Market.Location.Longitude))
	sb.WriteByte('\n')
	sb.WriteString(align("Tel-Nr"))
	sb.WriteString(md.Market.Phone)
	sb.WriteByte('\n')

	// Opening hours
	sb.WriteByte('\n')
	sb.WriteString(sep("Öffnungszeiten"))
	sb.WriteByte('\n')
	sb.WriteString(align("Aktuell"))
	sb.WriteString(md.Market.OpeningStatus.StatusText)
	sb.WriteString(" (")
	sb.WriteString(md.Market.OpeningStatus.InfoText)
	sb.WriteString(")\n")
	for _, oi := range md.Market.OpeningInfo {
		sb.WriteString(align(oi.Days))
		sb.WriteString(oi.Hours)
		sb.WriteByte('\n')
	}

	// Services
	sb.WriteByte('\n')
	sb.WriteString(sep("Services"))
	sb.WriteByte('\n')

	// Collect active services from both fixed and editable
	var active []MarketService
	for _, s := range md.Content.Services.Fixed {
		if s.Active {
			active = append(active, s)
		}
	}
	for _, s := range md.Content.Services.Editable {
		if s.Active {
			active = append(active, s)
		}
	}

	if len(active) == 0 {
		sb.WriteString("   (keine)\n")
	} else {
		for _, s := range active {
			sb.WriteString("   ")
			sb.WriteString(s.Text)
			if s.Tooltip != "" {
				sb.WriteString(" - ")
				sb.WriteString(s.Tooltip)
			}
			sb.WriteByte('\n')
		}
	}

	sb.WriteByte('\n')
	sb.WriteString(align("Pickup"))
	if md.Market.ServiceFlags.HasPickup {
		sb.WriteString("Ja")
	} else {
		sb.WriteString("Nein")
	}
	sb.WriteByte('\n')

	return sb.String()
}

func sep(base string) string {
	return "---" + base + strings.Repeat("-", 60-len(base))
}

func align(s string) string {
	al := 22
	if len(s) > al {
		al = 0
	} else {
		al -= len(s)
	}

	return "   " + s + ":" + strings.Repeat(" ", al)
}

func alignL(s string, n int) string {
	return strings.Repeat(" ", n*3) + s
}
