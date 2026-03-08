package rewerse

import (
	"encoding/json"
	"os"
	"testing"
)

func loadFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatalf("failed to read fixture %s: %v", name, err)
	}
	return data
}

func TestProductSearchResponseUnmarshal(t *testing.T) {
	var res productSearchResponse
	if err := json.Unmarshal(loadFixture(t, "product_search.json"), &res); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	p := res.Data.Products
	if p.Pagination.ObjectsPerPage == 0 {
		t.Error("objectsPerPage is 0")
	}
	if p.Pagination.ObjectCount == 0 {
		t.Error("objectCount is 0")
	}
	if p.Search.Term.Original == "" {
		t.Error("search term original is empty")
	}
	if len(p.Products) == 0 {
		t.Fatal("no products in response")
	}

	prod := p.Products[0]
	if prod.ProductID == "" {
		t.Error("productId is empty")
	}
	if prod.Title == "" {
		t.Error("title is empty")
	}
	if prod.Listing.ListingID == "" {
		t.Error("listingId is empty")
	}
	if prod.Listing.CurrentRetailPrice == 0 {
		t.Error("currentRetailPrice is 0")
	}
}

func TestProductDetailResponseUnmarshal(t *testing.T) {
	var res productDetailResponse
	if err := json.Unmarshal(loadFixture(t, "product_detail.json"), &res); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if len(res.Data.Product) == 0 {
		t.Fatal("no products in detail response")
	}

	pd := res.Data.Product[0]
	if pd.ProductID == "" {
		t.Error("productId is empty")
	}
	if pd.Brand == "" {
		t.Error("brand is empty")
	}
	if len(pd.NutritionFacts) == 0 {
		t.Error("nutritionFacts is empty")
	}
	if len(pd.NutritionFacts) > 0 {
		nf := pd.NutritionFacts[0]
		if len(nf.NutrientInformation) == 0 {
			t.Error("nutrientInformation is empty")
		}
		if len(nf.NutrientInformation) > 0 {
			ni := nf.NutrientInformation[0]
			if ni.NutrientType == "" {
				t.Error("nutrientType is empty")
			}
			if ni.QuantityContained.UomShortText == "" {
				t.Error("uomShortText is empty")
			}
		}
	}
}

func TestProductDetailFeatureBenefitNullUnmarshal(t *testing.T) {
	payload := []byte(`{"data":{"product":[{"productId":"1","title":"Test","featureBenefit":null}]}}`)

	var res productDetailResponse
	if err := json.Unmarshal(payload, &res); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	if len(res.Data.Product) != 1 {
		t.Fatalf("expected 1 product, got %d", len(res.Data.Product))
	}

	if res.Data.Product[0].FeatureBenefit != nil {
		t.Fatalf("expected nil featureBenefit for null input, got %#v", res.Data.Product[0].FeatureBenefit)
	}
}

func TestProductDetailFeatureBenefitArrayUnmarshal(t *testing.T) {
	payload := []byte(`{"data":{"product":[{"productId":"1","title":"Test","featureBenefit":["A","B"]}]}}`)

	var res productDetailResponse
	if err := json.Unmarshal(payload, &res); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	got := res.Data.Product[0].FeatureBenefit
	if len(got) != 2 || got[0] != "A" || got[1] != "B" {
		t.Fatalf("unexpected featureBenefit values: %#v", got)
	}
}

func TestProductDetailFeatureBenefitStringFailsUnmarshal(t *testing.T) {
	payload := []byte(`{"data":{"product":[{"productId":"1","title":"Test","featureBenefit":"legacy"}]}}`)

	var res productDetailResponse
	err := json.Unmarshal(payload, &res)
	if err == nil {
		t.Fatal("expected unmarshal error for string featureBenefit, got nil")
	}
}

func TestProductSuggestionsUnmarshal(t *testing.T) {
	var suggestions ProductSuggestions
	if err := json.Unmarshal(loadFixture(t, "product_suggestions.json"), &suggestions); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if len(suggestions) == 0 {
		t.Fatal("no suggestions in response")
	}
	s := suggestions[0]
	if s.Title == "" {
		t.Error("title is empty")
	}
	if s.RawValues.ProductID == "" {
		t.Error("rawValues.productId is empty")
	}
}
