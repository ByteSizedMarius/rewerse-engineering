package rewerse

import (
	"net/url"
)

func GetShopOverview(marketID string) (so ShopOverview, err error) {
	query := url.Values{}
	query.Add("serviceTypes", "PICKUP")
	query.Add("marketCode", marketID)

	req, err := BuildCustomRequest(apiHost, "v3/shop-overview?"+query.Encode())
	if err != nil {
		return
	}

	zipCode := "67065"
	headers := map[string]string{
		"rd-service-types": "PICKUP",
		"rd-customer-zip":  zipCode,
		"rd-postcode":      zipCode,
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	err = DoRequest(req, &so)
	if err != nil {
		return
	}

	return
}
