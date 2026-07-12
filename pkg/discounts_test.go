package rewerse

import (
	"encoding/json"
	"testing"
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
