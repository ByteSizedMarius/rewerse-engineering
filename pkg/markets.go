package rewerse

import (
	"fmt"
	"net/url"
)

// GetMarketDetails returns the details of the market with the given ID.
func GetMarketDetails(marketID string) (md MarketDetails, err error) {
	// Build Request
	query := url.Values{}
	query.Add("marketId", marketID)

	req, err := BuildCustomRequest(apiHost, "v3/market/details?"+query.Encode())
	if err != nil {
		return
	}

	// Send Request & Parse Response
	err = DoRequest(req, &md)
	if err != nil {
		return
	}

	return
}

// MarketSearch searches for markets based on the given query.
// Fuzzy search; accepts PLZ, city, marketname, street, etc.
func MarketSearch(searchQuery string) (markets Markets, err error) {
	// Build Request
	query := url.Values{}
	query.Add("search", searchQuery)

	req, err := BuildCustomRequest(apiHost, "v3/market/search?"+query.Encode())
	if err != nil {
		err = fmt.Errorf("error building request: %v", err)
		return
	}

	// Send Request & Parse Response
	var res struct {
		Markets Markets `json:"markets"`
	}
	err = DoRequest(req, &res)
	if err != nil {
		return
	}

	// Somewhat validate Result
	if len(res.Markets) == 0 {
		err = fmt.Errorf("no markets found for query \"%s\"", searchQuery)
		return
	}
	markets = res.Markets

	return
}
