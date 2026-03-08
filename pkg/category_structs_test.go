package rewerse

import (
	"encoding/json"
	"testing"
)

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
