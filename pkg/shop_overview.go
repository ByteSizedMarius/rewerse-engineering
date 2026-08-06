package rewerse

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
//
// Endpoint: GET /api/shop-overview. The endpoint takes no query parameters at all -
// market, service type and zip code are only read from the request headers.
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

	req, err := BuildCustomRequest(clientHost, "shop-overview")
	if err != nil {
		return
	}

	setDualHeader(req, "service-types", string(serviceType))
	setDualHeader(req, "customer-zip", zipCode)
	setDualHeader(req, "market-id", marketID)
	req.Header.Set("rd-postcode", zipCode)
	setCommonHeaders(req)

	var res shopOverviewResponse
	err = DoRequest(req, &res)
	if err != nil {
		return
	}

	so = ShopOverview{
		ProductRecalls:    res.Data.ProductRecalls.Products,
		ProductCategories: res.Data.Categories,
	}
	return
}
