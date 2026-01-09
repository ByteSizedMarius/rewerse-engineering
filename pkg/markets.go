package rewerse

import (
	"fmt"
	"net/url"
)

// GetMarketDetails returns the details of the market with the given ID.
func GetMarketDetails(marketID string) (md MarketDetails, err error) {
	req, err := BuildCustomRequest(clientHost, "stationary-markets/"+marketID)
	if err != nil {
		return
	}

	var res marketDetailsResponse
	err = DoRequest(req, &res)
	if err != nil {
		return
	}

	md = MarketDetails{
		Market:  res.Data.Market,
		Content: res.Data.Content,
	}
	return
}

// MarketSearch searches for markets based on the given query.
// Fuzzy search; accepts PLZ, city, marketname, street, etc.
func MarketSearch(searchQuery string) (markets Markets, err error) {
	query := url.Values{}
	query.Add("search", searchQuery)

	req, err := BuildCustomRequest(clientHost, "stationary-markets?"+query.Encode())
	if err != nil {
		return
	}

	var res marketSearchResponse
	err = DoRequest(req, &res)
	if err != nil {
		return
	}

	if len(res.Data.MarketSearch.Markets) == 0 {
		err = fmt.Errorf("no markets found for query \"%s\"", searchQuery)
		return
	}
	markets = res.Data.MarketSearch.Markets

	return
}
