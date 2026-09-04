package rewerse

import (
	"encoding/json"
	"testing"
	"time"
)

func TestParseLoyaltyBonus(t *testing.T) {
	tests := []struct {
		name string
		raw  *RawLoyaltyBonus
		want float64
	}{
		{
			name: "cent bonus",
			raw:  &RawLoyaltyBonus{BonusType: "cent", BonusValue: 60},
			want: 0.60,
		},
		{
			name: "uppercase cent",
			raw:  &RawLoyaltyBonus{BonusType: "CENT", BonusValue: 100},
			want: 1.00,
		},
		{
			name: "nil",
			raw:  nil,
			want: 0,
		},
		{
			name: "zero value",
			raw:  &RawLoyaltyBonus{BonusType: "cent", BonusValue: 0},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseLoyaltyBonus(tt.raw); got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseLoyaltyBonusFromFixture(t *testing.T) {
	var rd RawDiscounts
	if err := json.Unmarshal(loadFixture(t, "discounts.json"), &rd); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	var yfood *RawOffer
	for _, cat := range rd.Data.Offers.Current.Categories {
		if cat.ID != "rewe-bonus-produkte" {
			continue
		}
		for i := range cat.Offers {
			if cat.Offers[i].Title == "YFood Trinkmahlzeit" {
				yfood = &cat.Offers[i]
				break
			}
		}
	}
	if yfood == nil {
		t.Fatal("YFood offer not found in fixture")
	}

	if got := parseLoyaltyBonus(yfood.LoyaltyBonus); got != 0.30 {
		t.Fatalf("loyalty bonus = %v, want 0.30", got)
	}
}

func TestCleanDiscounts(t *testing.T) {
	var rd RawDiscounts
	if err := json.Unmarshal(loadFixture(t, "discounts.json"), &rd); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	ds, err := cleanDiscounts(rd)
	if err != nil {
		t.Fatalf("cleanDiscounts failed: %v", err)
	}
	if len(ds.Categories) != 2 {
		t.Fatalf("categories = %d, want 2", len(ds.Categories))
	}
	if want := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC); !ds.ValidUntil.Equal(want) {
		t.Fatalf("ValidUntil = %v, want %v", ds.ValidUntil, want)
	}

	top := ds.Categories[0]
	if top.ID != "markt-topangebote" {
		t.Fatalf("first category ID = %q, want markt-topangebote", top.ID)
	}
	if len(top.Offers) == 0 {
		t.Fatalf("first category has no offers")
	}
	if top.Offers[0].Manufacturer != "HARIBO" {
		t.Fatalf("first offer manufacturer = %q, want HARIBO", top.Offers[0].Manufacturer)
	}
	if top.Offers[0].Price <= 0 {
		t.Fatalf("first offer price = %v, want > 0", top.Offers[0].Price)
	}
	if top.Offers[0].PriceParseFail {
		t.Fatalf("first offer price parse failed for %q", top.Offers[0].PriceRaw)
	}

	var bonus *DiscountCategory
	for i := range ds.Categories {
		if ds.Categories[i].ID == "rewe-bonus-produkte" {
			bonus = &ds.Categories[i]
		}
	}
	if bonus == nil {
		t.Fatal("category rewe-bonus-produkte not found")
	}
	if len(bonus.Offers) != 1 {
		t.Fatalf("rewe-bonus-produkte offers = %d, want 1", len(bonus.Offers))
	}
	if bonus.Offers[0].LoyaltyBonus != 0.30 {
		t.Fatalf("loyalty bonus = %v, want 0.30", bonus.Offers[0].LoyaltyBonus)
	}
}
