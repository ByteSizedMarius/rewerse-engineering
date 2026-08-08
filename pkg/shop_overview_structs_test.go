package rewerse

import (
	"encoding/json"
	"testing"
)

func TestShopOverviewUnmarshal(t *testing.T) {
	var res shopOverviewResponse
	if err := json.Unmarshal(loadFixture(t, "shop_overview.json"), &res); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if len(res.Data.Categories) == 0 {
		t.Fatal("no product categories")
	}
	if len(res.Data.ProductRecalls.Products) == 0 {
		t.Error("no product recalls")
	}

	root := res.Data.Categories[0]
	if root.Name == "" {
		t.Error("root category name is empty")
	}
	if root.Slug == "" {
		t.Error("root category slug is empty")
	}
	if len(root.CategoryTags) == 0 {
		t.Error("root category has no categoryTags")
	}

	// Verify recursive child categories unmarshal
	if len(root.ChildCategories) == 0 {
		t.Fatal("root has no child categories")
	}
	child := root.ChildCategories[0]
	if child.Name == "" {
		t.Error("child category name is empty")
	}
	// Leaf categories have childCategories: null
	if child.ChildCategories != nil {
		t.Errorf("expected leaf child to have nil childCategories, got %d", len(child.ChildCategories))
	}
}
