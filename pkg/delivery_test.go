package rewerse

import (
	"encoding/json"
	"testing"
)

func TestBulkyGoodsResponseUnmarshal(t *testing.T) {
	var res bulkyGoodsResponse
	if err := json.Unmarshal(loadFixture(t, "bulky_goods.json"), &res); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	cfg := res.Data.BulkyGoodsConfiguration
	if !cfg.HasBeverageSurcharge {
		t.Error("expected hasBeverageSurcharge to be true")
	}
	if cfg.BeverageSurcharge == nil {
		t.Fatal("beverageSurcharge is nil")
	}
	if cfg.BeverageSurcharge.Surcharge == 0 {
		t.Error("surcharge is zero")
	}
	if cfg.BeverageSurcharge.HardLimit <= cfg.BeverageSurcharge.SoftLimit {
		t.Errorf("expected hardLimit > softLimit, got %d and %d",
			cfg.BeverageSurcharge.HardLimit, cfg.BeverageSurcharge.SoftLimit)
	}
	if len(cfg.BeverageSurcharge.DisplayTexts) == 0 {
		t.Error("no displayTexts")
	}
}
