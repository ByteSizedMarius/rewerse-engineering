package rewerse

import (
	"encoding/json"
	"testing"
)

func TestRecipeSearchResultsUnmarshal(t *testing.T) {
	var res RecipeSearchResults
	if err := json.Unmarshal(loadFixture(t, "recipe_search.json"), &res); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if res.TotalCount == 0 {
		t.Error("totalCount is 0")
	}
	if len(res.Recipes) == 0 {
		t.Fatal("no recipes in response")
	}
	r := res.Recipes[0]
	if r.ID == "" {
		t.Error("recipe id is empty")
	}
	if r.Title == "" {
		t.Error("recipe title is empty")
	}
	if len(res.Metadata.Difficulties) == 0 {
		t.Error("metadata difficulties is empty")
	}
}

func TestRecipeDetailsUnmarshal(t *testing.T) {
	var res RecipeDetails
	if err := json.Unmarshal(loadFixture(t, "recipe_details.json"), &res); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	r := res.Recipe
	if r.ID == "" {
		t.Error("recipe id is empty")
	}
	if r.Title == "" {
		t.Error("recipe title is empty")
	}
	if r.Ingredients.Portions == 0 {
		t.Error("portions is 0")
	}
	if len(r.Ingredients.Items) == 0 {
		t.Error("no ingredients")
	}
	if len(r.Steps) == 0 {
		t.Error("no steps")
	}
}

func TestRecipePopularTermsUnmarshal(t *testing.T) {
	var terms PopularSearchTerms
	if err := json.Unmarshal(loadFixture(t, "recipe_popular_terms.json"), &terms); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if len(terms) == 0 {
		t.Fatal("no popular terms")
	}
	if terms[0].Title == "" {
		t.Error("term title is empty")
	}
}
