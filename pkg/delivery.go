package rewerse

import "fmt"

type bulkyGoodsResponse struct {
	Data struct {
		BulkyGoodsConfiguration BulkyGoodsConfig `json:"bulkyGoodsConfiguration"`
	} `json:"data"`
}

// BulkyGoodsConfig contains beverage crate limits and surcharges for a market.
// Endpoint: GET /api/bulky-goods-configuration/{marketCode}/service-types/{serviceType}
type BulkyGoodsConfig struct {
	// HasBeverageSurcharge indicates if beverage surcharges apply
	HasBeverageSurcharge bool `json:"hasBeverageSurcharge"`
	// BeverageSurcharge contains the surcharge details (nil if none)
	BeverageSurcharge *BeverageSurcharge `json:"beverageSurcharge"`
}

// BeverageSurcharge contains limits and pricing for beverage crates
type BeverageSurcharge struct {
	// SoftLimit threshold before surcharge applies
	SoftLimit int `json:"softLimit"`
	// HardLimit maximum threshold
	HardLimit int `json:"hardLimit"`
	// Surcharge amount in cents: 190 = 1.90 EUR
	Surcharge int `json:"surcharge"`
	// DisplayTexts are German explanations of the surcharge rules
	DisplayTexts []string `json:"displayTexts"`
}

// GetBulkyGoodsConfig retrieves beverage crate limits and surcharges for a market.
// serviceType should be "PICKUP" or "DELIVERY".
func GetBulkyGoodsConfig(marketID string, serviceType ServiceType) (BulkyGoodsConfig, error) {
	if err := validateServiceType(serviceType); err != nil {
		return BulkyGoodsConfig{}, err
	}

	path := fmt.Sprintf("bulky-goods-configuration/%s/service-types/%s", marketID, serviceType)
	req, err := BuildCustomRequest(clientHost, path)
	if err != nil {
		return BulkyGoodsConfig{}, err
	}

	req.Header.Set("x-rd-market-id", marketID)
	req.Header.Set("x-rd-service-types", string(serviceType))
	setCommonHeaders(req)

	var res bulkyGoodsResponse
	if err := DoRequest(req, &res); err != nil {
		return BulkyGoodsConfig{}, err
	}

	return res.Data.BulkyGoodsConfiguration, nil
}
