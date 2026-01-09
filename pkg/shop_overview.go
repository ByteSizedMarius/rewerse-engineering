package rewerse

import (
	"net/url"
)

// ShopOverviewOpts configures the shop overview request
type ShopOverviewOpts struct {
	// ServiceType is "PICKUP" or "DELIVERY" (default: PICKUP)
	ServiceType ServiceType
	// ZipCode is required for DELIVERY, optional for PICKUP
	ZipCode string
}

// GetShopOverview retrieves the product categories for a market (PICKUP mode).
// For delivery mode, use GetShopOverviewWithOpts.
func GetShopOverview(marketID string) (so ShopOverview, err error) {
	return GetShopOverviewWithOpts(marketID, nil)
}

// GetShopOverviewWithOpts retrieves the product categories with configurable service type.
// For DELIVERY mode, opts.ZipCode is required.
func GetShopOverviewWithOpts(marketID string, opts *ShopOverviewOpts) (so ShopOverview, err error) {
	serviceType := ServicePickup
	zipCode := "67065"

	if opts != nil {
		if opts.ServiceType != "" {
			serviceType = opts.ServiceType
		}
		if opts.ZipCode != "" {
			zipCode = opts.ZipCode
		}
	}

	query := url.Values{}
	query.Add("serviceTypes", string(serviceType))
	query.Add("marketCode", marketID)
	if serviceType == ServiceDelivery {
		query.Add("deliveryZipCode", zipCode)
	}

	req, err := BuildCustomRequest(apiHost, "v3/shop-overview?"+query.Encode())
	if err != nil {
		return
	}

	setDualHeader(req, "service-types", string(serviceType))
	setDualHeader(req, "customer-zip", zipCode)
	req.Header.Set("rd-postcode", zipCode)
	req.Header.Set("x-rd-market-id", marketID)
	setCommonHeaders(req)

	err = DoRequest(req, &so)
	return
}
