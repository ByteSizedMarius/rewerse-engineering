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
}
