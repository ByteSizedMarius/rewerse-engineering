package rewerse

import (
	"encoding/json"
	"testing"
)

func TestMarketSearchResponseUnmarshal(t *testing.T) {
	var res marketSearchResponse
	if err := json.Unmarshal(loadFixture(t, "market_search.json"), &res); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	markets := res.Data.MarketSearch.Markets
	if len(markets) == 0 {
		t.Fatal("no markets in response")
	}

	m := markets[0]
	if m.WWIdent == "" {
		t.Error("wwIdent is empty")
	}
	if m.Name == "" {
		t.Error("name is empty")
	}
	if m.Street == "" {
		t.Error("street is empty")
	}
	if m.Location.Latitude == 0 && m.Location.Longitude == 0 {
		t.Error("location is zero")
	}
	if m.OpeningStatus.OpenState == "" {
		t.Error("openState is empty")
	}
	if len(m.OpeningInfo) == 0 {
		t.Error("openingInfo is empty")
	}
	if m.Category.Template == "" {
		t.Error("category template is empty")
	}
}

func TestMarketDetailsResponseUnmarshal(t *testing.T) {
	var res marketDetailsResponse
	if err := json.Unmarshal(loadFixture(t, "market_details.json"), &res); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if res.Data.Market.WWIdent == "" {
		t.Error("market wwIdent is empty")
	}
	if res.Data.Market.Street == "" {
		t.Error("market street is empty")
	}
	if len(res.Data.Content.Services.Fixed) == 0 {
		t.Error("no fixed services")
	}
}
