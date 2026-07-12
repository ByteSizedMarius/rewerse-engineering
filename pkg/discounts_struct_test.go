package rewerse

import (
	"encoding/json"
	"testing"
)

func TestRawDiscountsUnmarshal(t *testing.T) {
	var res RawDiscounts
	if err := json.Unmarshal(loadFixture(t, "discounts.json"), &res); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	offers := res.Data.Offers
	if offers.DefaultWeek == "" {
		t.Error("defaultWeek is empty")
	}

	current := offers.Current
	if current.UntilDate == "" {
		t.Error("current untilDate is empty")
	}
	if len(current.Categories) == 0 {
		t.Fatal("no categories in current week")
	}

	cat := current.Categories[0]
	if cat.ID == "" {
		t.Error("category id is empty")
	}
	if cat.Title == "" {
		t.Error("category title is empty")
	}
	if len(cat.Offers) == 0 {
		t.Fatal("no offers in category")
	}

	offer := cat.Offers[0]
	if offer.Title == "" {
		t.Error("offer title is empty")
	}
	if offer.PriceData.Price == "" {
		t.Error("offer price is empty")
	}
	if len(offer.Images) == 0 {
		t.Error("offer images is empty")
	}
	if offer.RawValues.Nan == "" {
		t.Error("offer nan is empty")
	}
	if offer.LoyaltyBonus != nil && offer.LoyaltyBonus.BonusType == "" {
		t.Error("loyalty bonus type is empty")
	}
}

func TestDiscountFormatPriceLine(t *testing.T) {
	tests := []struct {
		name string
		d    Discount
		want string
	}{
		{
			name: "shelf price only",
			d:    Discount{Title: "Haribo", Price: 0.77},
			want: "Haribo, 0.77€",
		},
		{
			name: "price and bonus",
			d:    Discount{Title: "Dr. Oetker Pizza", Price: 1.99, LoyaltyBonus: 0.10},
			want: "Dr. Oetker Pizza, 1.99€ (0.10€ Bonus)",
		},
		{
			name: "bonus only",
			d:    Discount{Title: "Dunkle Pflaumen", LoyaltyBonus: 0.60},
			want: "Dunkle Pflaumen, 0.60€ Bonus",
		},
		{
			name: "parse failure fallback",
			d:    Discount{Title: "Broken", PriceRaw: "ab 1,99 €", PriceParseFail: true},
			want: "Broken, ab 1,99 €",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.formatPriceLine(); got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
