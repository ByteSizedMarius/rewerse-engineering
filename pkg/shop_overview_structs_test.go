package rewerse

import (
	"encoding/json"
	"testing"
)

func TestShopOverviewUnmarshal(t *testing.T) {
	var so ShopOverview
	if err := json.Unmarshal(loadFixture(t, "shop_overview.json"), &so); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if len(so.ProductCategories) == 0 {
		t.Fatal("no product categories")
	}

	root := so.ProductCategories[0]
	if root.Name == "" {
		t.Error("root category name is empty")
	}
	if root.Slug == "" {
		t.Error("root category slug is empty")
	}

	// Verify recursive child categories unmarshal
	if len(root.ChildCategories) == 0 {
		t.Error("root has no child categories")
	}
	if len(root.ChildCategories) > 0 {
		child := root.ChildCategories[0]
		if child.Name == "" {
			t.Error("child category name is empty")
		}
	}
}
