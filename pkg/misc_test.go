package rewerse

import (
	"encoding/json"
	"testing"
)

func TestRecallsResponseUnmarshal(t *testing.T) {
	var res recallsResponse
	if err := json.Unmarshal(loadFixture(t, "recalls.json"), &res); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	recalls := res.Data.ProductRecalls.Products
	if len(recalls) == 0 {
		t.Fatal("no recalls in response")
	}
	if recalls[0].SubjectProduct == "" {
		t.Error("subjectProduct is empty")
	}
	if recalls[0].URL == "" {
		t.Error("url is empty")
	}
}

func TestRecipeHubUnmarshal(t *testing.T) {
	var hub RecipeHub
	if err := json.Unmarshal(loadFixture(t, "recipe_hub.json"), &hub); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if hub.RecipeOfTheDay.Title == "" {
		t.Error("recipeOfTheDay title is empty")
	}
	if len(hub.PopularRecipes) == 0 {
		t.Error("no popular recipes")
	}
	if len(hub.Categories) == 0 {
		t.Error("no categories")
	}
	if hub.Categories[0].Title == "" {
		t.Error("category title is empty")
	}
}

func TestServicePortfolioResponseUnmarshal(t *testing.T) {
	var res servicePortfolioResponse
	if err := json.Unmarshal(loadFixture(t, "service_portfolio.json"), &res); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	sp := res.Data.ServicePortfolio
	if sp.CustomerZipCode == "" {
		t.Error("customerZipCode is empty")
	}
	if len(sp.PickupMarkets) == 0 {
		t.Error("no pickup markets")
	}
	pm := sp.PickupMarkets[0]
	if pm.WWIdent == "" {
		t.Error("pickup market wwIdent is empty")
	}
	if pm.City == "" {
		t.Error("pickup market city is empty")
	}
}


