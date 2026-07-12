package rewerse

import (
	"encoding/json"
	"testing"
)

func TestParseLoyaltyBonus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		raw   *RawLoyaltyBonus
		want  *LoyaltyBonus
	}{
		{
			name: "cent bonus",
			raw:  &RawLoyaltyBonus{BonusType: "cent", BonusValue: 60},
			want: &LoyaltyBonus{BonusType: "cent", BonusValue: 0.60},
		},
		{
			name: "uppercase cent",
			raw:  &RawLoyaltyBonus{BonusType: "CENT", BonusValue: 100},
			want: &LoyaltyBonus{BonusType: "CENT", BonusValue: 1.00},
		},
		{
			name: "nil",
			raw:  nil,
			want: nil,
		},
		{
			name: "zero value",
			raw:  &RawLoyaltyBonus{BonusType: "cent", BonusValue: 0},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := parseLoyaltyBonus(tt.raw)
			if tt.want == nil {
				if got != nil {
					t.Fatalf("got %+v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("got nil, want bonus")
			}
			if got.BonusType != tt.want.BonusType || got.BonusValue != tt.want.BonusValue {
				t.Fatalf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestDiscountFormatPriceLine(t *testing.T) {
	t.Parallel()

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
			d: Discount{
				Title:        "Dr. Oetker Pizza",
				Price:        1.99,
				LoyaltyBonus: &LoyaltyBonus{BonusType: "cent", BonusValue: 0.10},
			},
			want: "Dr. Oetker Pizza, 1.99€ (+0.10€ Bonus)",
		},
		{
			name: "bonus only",
			d: Discount{
				Title:        "Dunkle Pflaumen",
				LoyaltyBonus: &LoyaltyBonus{BonusType: "cent", BonusValue: 0.60},
			},
			want: "Dunkle Pflaumen, 0.60€ Bonus",
		},
		{
			name: "parse failure fallback",
			d: Discount{
				Title:          "Broken",
				PriceRaw:       "ab 1,99 €",
				PriceParseFail: true,
			},
			want: "Broken, ab 1,99 €",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.d.formatPriceLine(); got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseDiscountFromFixture(t *testing.T) {
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

	discount := Discount{
		Title:        yfood.Title,
		Price:        3.29,
		PriceRaw:     yfood.PriceData.Price,
		LoyaltyBonus: parseLoyaltyBonus(yfood.LoyaltyBonus),
	}

	if discount.Price != 3.29 {
		t.Fatalf("price = %v, want 3.29", discount.Price)
	}
	if discount.LoyaltyBonus == nil || discount.LoyaltyBonus.BonusValue != 0.30 {
		t.Fatalf("loyalty bonus = %+v, want 0.30", discount.LoyaltyBonus)
	}
	if got := discount.formatPriceLine(); got != "YFood Trinkmahlzeit, 3.29€ (+0.30€ Bonus)" {
		t.Fatalf("formatPriceLine = %q", got)
	}
}
