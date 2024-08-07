package rewerse

import (
	"net/url"
	"strconv"
)

var defaultOpts = ProductOpts{
	Page:           1,
	ObjectsPerPage: 30,
}

type ProductOpts struct {
	Page           int
	ObjectsPerPage int
}

func GetCategoryProducts(marketID, categorySlug string, opts *ProductOpts) (ProductSearchResults, error) {
	query := url.Values{}
	query.Add("query", "")
	query.Add("categorySlug", categorySlug)
	return getProductsCommon(marketID, opts, query)
}

func GetProducts(marketID, search string, opts *ProductOpts) (ProductSearchResults, error) {
	query := url.Values{}
	query.Add("query", search)
	return getProductsCommon(marketID, opts, query)
}

//===================================================
//
// Internals
//
//===================================================

func getProductsCommon(marketID string, opts *ProductOpts, queryParams url.Values) (ProductSearchResults, error) {
	if opts == nil {
		opts = &defaultOpts
	} else {
		if opts.Page == 0 {
			opts.Page = defaultOpts.Page
		}
		if opts.ObjectsPerPage == 0 {
			opts.ObjectsPerPage = defaultOpts.ObjectsPerPage
		}
	}

	queryParams.Add("page", strconv.Itoa(opts.Page))
	queryParams.Add("objectsPerPage", strconv.Itoa(opts.ObjectsPerPage))

	req, err := BuildCustomRequest(clientHost, "products?"+queryParams.Encode())
	if err != nil {
		return ProductSearchResults{}, err
	}

	zipCode := "67065"
	headers := map[string]string{
		"rd-service-types": "PICKUP",
		"rd-customer-zip":  zipCode,
		"rd-postcode":      zipCode,
		"rd-market-id":     marketID,
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	var psr ProductSearchResults
	err = DoRequest(req, &psr)
	if err != nil {
		return ProductSearchResults{}, err
	}

	return psr, nil
}
