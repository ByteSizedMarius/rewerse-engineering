package rewerse

import (
	"encoding/json"
	"testing"
)

func TestRecipeSearchResponseUnmarshal(t *testing.T) {
	var res recipeSearchResponse
	if err := json.Unmarshal(loadFixture(t, "recipe_search.json"), &res); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	r := res.Data.Recipes
	if len(r.Recipes) != 8 {
		t.Fatalf("expected 8 recipes, got %d", len(r.Recipes))
	}
	if r.Metadata.TotalRecipeCount != 42 {
		t.Errorf("expected totalRecipeCount 42, got %d", r.Metadata.TotalRecipeCount)
	}
	if len(r.Metadata.Collections) == 0 {
		t.Error("collections is empty")
	}
	if len(r.Metadata.Difficulties) == 0 {
		t.Error("difficulties is empty")
	}

	first := r.Recipes[0]
	if first.Title != "Nudeln mit Tomatensoße" {
		t.Errorf("expected first title \"Nudeln mit Tomatensoße\", got %q", first.Title)
	}
	if first.ID == "" {
		t.Error("first recipe id is empty")
	}
	if first.Difficulty != DifficultyEasy {
		t.Errorf("expected difficulty %d, got %d", DifficultyEasy, first.Difficulty)
	}
	if len(first.ImageURLs) == 0 {
		t.Error("first recipe has no image urls")
	}

	nilTotal := false
	for _, rec := range r.Recipes {
		if rec.TimeTotal == nil {
			nilTotal = true
			break
		}
	}
	if !nilTotal {
		t.Error("expected at least one recipe with null timeTotal")
	}

	if res.Data.RecipeOfTheDay != nil {
		t.Error("expected no recipeOfTheDay")
	}
}

func TestRecipeSearchEmptyUnmarshal(t *testing.T) {
	var res recipeSearchResponse
	if err := json.Unmarshal(loadFixture(t, "recipe_search_empty.json"), &res); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	r := res.Data.Recipes
	if len(r.Recipes) != 0 {
		t.Errorf("expected 0 recipes, got %d", len(r.Recipes))
	}
	if r.Metadata.TotalRecipeCount != 0 {
		t.Errorf("expected totalRecipeCount 0, got %d", r.Metadata.TotalRecipeCount)
	}
	if len(res.Errors) != 0 {
		t.Errorf("expected no errors, got %v", res.Errors)
	}
}

func TestRecipeSearchRotdUnmarshal(t *testing.T) {
	var res recipeSearchResponse
	if err := json.Unmarshal(loadFixture(t, "recipe_search_rotd.json"), &res); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	rotd := res.Data.RecipeOfTheDay
	if rotd == nil {
		t.Fatal("recipeOfTheDay is nil")
	}
	if rotd.Title != "Ofengemüse mit Kräuterquark" {
		t.Errorf("expected title \"Ofengemüse mit Kräuterquark\", got %q", rotd.Title)
	}
	if len(rotd.Ingredients) == 0 {
		t.Error("recipeOfTheDay has no ingredients")
	}
	if len(rotd.Preparation.Steps) == 0 {
		t.Error("recipeOfTheDay has no steps")
	}
}

func TestRecipeDetailsResponseUnmarshal(t *testing.T) {
	var res recipeDetailsResponse
	if err := json.Unmarshal(loadFixture(t, "recipe_details.json"), &res); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	r := res.Data.RecipeByID
	if r == nil {
		t.Fatal("recipeById is nil")
	}
	if r.Title != "Blechkuchen mit Äpfeln" {
		t.Errorf("expected title \"Blechkuchen mit Äpfeln\", got %q", r.Title)
	}
	if r.Serving.Type != "pieces" || r.Serving.Quantity != 20 {
		t.Errorf("expected serving pieces/20, got %s/%d", r.Serving.Type, r.Serving.Quantity)
	}
	if r.TimeTotal == nil || *r.TimeTotal != 75 {
		t.Errorf("expected timeTotal 75, got %v", r.TimeTotal)
	}
	if len(r.Tags) != 6 {
		t.Errorf("expected 6 tags, got %d", len(r.Tags))
	}
	if len(r.ImageURLs) == 0 || r.ImageURLs[0].URLs.Default == "" {
		t.Error("image url default is empty")
	}
	if len(r.Preparation.Steps) == 0 {
		t.Error("preparation steps are empty")
	}
	if len(r.Nutrients) == 0 {
		t.Error("nutrients are empty")
	}
	if len(r.VitalSubstances) == 0 {
		t.Error("vitalSubstances are empty")
	}

	nilUnit := false
	for _, i := range r.Ingredients {
		if i.Unit == nil {
			nilUnit = true
			break
		}
	}
	if !nilUnit {
		t.Error("expected at least one ingredient with null unit")
	}
}

func TestRecipePopularTermsResponseUnmarshal(t *testing.T) {
	var res recipePopularTermsResponse
	if err := json.Unmarshal(loadFixture(t, "recipe_popular_terms.json"), &res); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	tags := res.Data.SearchTags.Tags
	if len(tags) != 13 {
		t.Fatalf("expected 13 tags, got %d", len(tags))
	}
	if tags[0].Title != "Lachs" {
		t.Errorf("expected first title \"Lachs\", got %q", tags[0].Title)
	}
}
