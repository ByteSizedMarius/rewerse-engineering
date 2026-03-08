package rewerse

import (
	"encoding/json"
	"testing"
)

func TestBasketResponseUnmarshal(t *testing.T) {
	var res basketResponse
	if err := json.Unmarshal(loadFixture(t, "basket_create.json"), &res); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	b := res.Data.Basket
	if b.ID != "b4a5c6d7-e8f9-4a0b-1c2d-3e4f5a6b7c8d" {
		t.Errorf("basket id: expected b4a5c6d7-e8f9-4a0b-1c2d-3e4f5a6b7c8d, got %s", b.ID)
	}
	if b.DeviceID != "d1e2v3c4-a5b6-4c7d-8e9f-0a1b2c3d4e5f" {
		t.Errorf("deviceId: expected d1e2v3c4-a5b6-4c7d-8e9f-0a1b2c3d4e5f, got %s", b.DeviceID)
	}
	if b.Version != 0 {
		t.Errorf("version: expected 0, got %d", b.Version)
	}
	if b.ServiceSelection.WWIdent != "831002" {
		t.Errorf("serviceSelection wwIdent: expected 831002, got %s", b.ServiceSelection.WWIdent)
	}
	if b.ServiceSelection.ServiceType != "PICKUP" {
		t.Errorf("serviceSelection serviceType: expected PICKUP, got %s", b.ServiceSelection.ServiceType)
	}
	if b.ServiceSelection.ZipCode != "67065" {
		t.Errorf("serviceSelection zipCode: expected 67065, got %s", b.ServiceSelection.ZipCode)
	}
	if b.Staggerings.ReachedStaggering == nil {
		t.Error("reachedStaggering: expected non-nil")
	}
	if b.Staggerings.NextStaggering == nil {
		t.Error("nextStaggering: expected non-nil")
	}
}
