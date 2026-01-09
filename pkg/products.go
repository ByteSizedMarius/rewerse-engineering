package rewerse

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// ProductFilter is a filter type for product search.
// Use with ProductOpts.Filters to narrow down search results.
type ProductFilter string

// Product filters for search queries. Not exhaustive.
const (
	FilterNew        ProductFilter = "attribute=new"        // Recently added products
	FilterOrganic    ProductFilter = "attribute=organic"    // Organic/Bio products
	FilterVegan      ProductFilter = "attribute=vegan"      // Vegan products
	FilterVegetarian ProductFilter = "attribute=vegetarian" // Vegetarian products
	FilterRegional   ProductFilter = "attribute=regional"   // Regional/local products
)

var defaultOpts = ProductOpts{
	Page:           1,
	ObjectsPerPage: 30,
}

// ProductOpts configures product search requests.
type ProductOpts struct {
	Page           int
	ObjectsPerPage int
	// ServiceType must match market capabilities (PICKUP or DELIVERY).
	// Check market details hasPickup field. Returns HTTP 400 if mismatched.
	ServiceType ServiceType
	Filters     []ProductFilter
}

// GetProducts searches for products by query string.
// Endpoint: GET /api/products?query={search}
// ServiceType in opts must match market capabilities - use GetMarketDetails to check hasPickup.
// Returns HTTP 400 "Invalid market selection" if service type doesn't match.
func GetProducts(marketID, search string, opts *ProductOpts) (ProductResults, error) {
	if marketID == "" {
		return ProductResults{}, fmt.Errorf("marketID: cannot be empty")
	}
	query := url.Values{}
	query.Add("query", search)
	return getProductsCommon(marketID, opts, query)
}

// GetCategoryProducts returns products in a specific category.
// Endpoint: GET /api/products?categorySlug={slug}
// ServiceType in opts must match market capabilities - use GetMarketDetails to check hasPickup.
// Returns HTTP 400 "Invalid market selection" if service type doesn't match.
func GetCategoryProducts(marketID, categorySlug string, opts *ProductOpts) (ProductResults, error) {
	if marketID == "" {
		return ProductResults{}, fmt.Errorf("marketID: cannot be empty")
	}
	query := url.Values{}
	query.Add("query", "")
	query.Add("categorySlug", categorySlug)
	return getProductsCommon(marketID, opts, query)
}

// GetProductByID returns full product details by product ID.
// Endpoint: GET /api/products/{productId}
func GetProductByID(marketID, productID string) (ProductDetail, error) {
	if marketID == "" {
		return ProductDetail{}, fmt.Errorf("marketID: cannot be empty")
	}
	req, err := BuildCustomRequest(clientHost, "products/"+productID)
	if err != nil {
		return ProductDetail{}, err
	}

	setMarketHeaders(req, marketID, "")

	var res productDetailResponse
	err = DoRequest(req, &res)
	if err != nil {
		return ProductDetail{}, err
	}

	if len(res.Data.Product) == 0 {
		return ProductDetail{}, fmt.Errorf("product %s not found", productID)
	}

	return res.Data.Product[0], nil
}

// GetProductSuggestions returns search autocomplete suggestions.
// Endpoint: GET /products/suggestion-search (note: no /api prefix)
func GetProductSuggestions(query string, opts *ProductOpts) (ProductSuggestions, error) {
	if opts == nil {
		opts = &ProductOpts{Page: 1, ObjectsPerPage: 25}
	}
	if opts.Page <= 0 {
		opts.Page = 1
	}
	if opts.ObjectsPerPage <= 0 {
		opts.ObjectsPerPage = 25
	}

	params := url.Values{}
	params.Add("query", query)
	params.Add("page", strconv.Itoa(opts.Page))
	params.Add("objectsPerPage", strconv.Itoa(opts.ObjectsPerPage))

	req, err := BuildCustomRequestRaw(apiHost, "/products/suggestion-search?"+params.Encode())
	if err != nil {
		return nil, err
	}

	setCommonHeaders(req)
	req.Header.Set("ruleversion", "2")

	var suggestions ProductSuggestions
	err = DoRequest(req, &suggestions)
	return suggestions, err
}

// GetProductRecommendations returns related product recommendations.
// Endpoint: GET /api/products/recommendations?listingIds={listingId}
func GetProductRecommendations(marketID, listingID string) ([]Product, error) {
	if marketID == "" {
		return nil, fmt.Errorf("marketID: cannot be empty")
	}
	query := url.Values{}
	query.Add("listingIds", listingID)

	req, err := BuildCustomRequest(clientHost, "products/recommendations?"+query.Encode())
	if err != nil {
		return nil, err
	}

	setMarketHeaders(req, marketID, "")

	var res productRecommendationsResponse
	err = DoRequest(req, &res)
	if err != nil {
		return nil, err
	}

	return res.Data.ProductRecommendations.Recommendations, nil
}

func getProductsCommon(marketID string, opts *ProductOpts, queryParams url.Values) (ProductResults, error) {
	if opts == nil {
		opts = &defaultOpts
	} else {
		if opts.Page <= 0 {
			opts.Page = defaultOpts.Page
		}
		if opts.ObjectsPerPage <= 0 {
			opts.ObjectsPerPage = defaultOpts.ObjectsPerPage
		}
	}

	queryParams.Add("page", strconv.Itoa(opts.Page))
	queryParams.Add("objectsPerPage", strconv.Itoa(opts.ObjectsPerPage))

	for _, filter := range opts.Filters {
		queryParams.Add("filters", string(filter))
	}

	req, err := BuildCustomRequest(clientHost, "products?"+queryParams.Encode())
	if err != nil {
		return ProductResults{}, err
	}

	setMarketHeaders(req, marketID, opts.ServiceType)

	var raw productSearchResponse
	err = DoRequest(req, &raw)
	if err != nil {
		return ProductResults{}, err
	}

	return ProductResults{
		Products: raw.Data.Products.Products,
		Pagination: struct {
			ObjectsPerPage int
			CurrentPage    int
			PageCount      int
			ObjectCount    int
		}{
			ObjectsPerPage: raw.Data.Products.Pagination.ObjectsPerPage,
			CurrentPage:    raw.Data.Products.Pagination.CurrentPage,
			PageCount:      raw.Data.Products.Pagination.PageCount,
			ObjectCount:    raw.Data.Products.Pagination.ObjectCount,
		},
		SearchTerm: struct {
			Original  string
			Corrected *string
		}{
			Original:  raw.Data.Products.Search.Term.Original,
			Corrected: raw.Data.Products.Search.Term.Corrected,
		},
	}, nil
}

// setMarketHeaders adds market-specific headers for product API requests.
// Uses default zip 67065 and defaults to PICKUP if serviceType is empty.
func setMarketHeaders(req *http.Request, marketID string, serviceType ServiceType) {
	if serviceType == "" {
		serviceType = ServicePickup
	}
	setCommonHeaders(req)
	setDualHeader(req, "service-types", string(serviceType))
	setDualHeader(req, "customer-zip", "67065")
	setDualHeader(req, "market-id", marketID)
	req.Header.Set("rd-postcode", "67065")
	req.Header.Set("rd-is-lsfk", "false")
}
