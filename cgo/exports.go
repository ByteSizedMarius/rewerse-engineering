package main

/*
#include <stdlib.h>

#ifdef _WIN32
#include <windows.h>

// Pin this DLL to prevent FreeLibrary from unloading it.
// Go's runtime doesn't support being unloaded.
// Best-effort: if pinning fails, there's no recovery path - the process
// would crash on exit if the DLL is forcibly unloaded. In practice,
// GetModuleHandleExW with PIN flag rarely fails.
static void pinDLL() {
    HMODULE self;
    GetModuleHandleExW(
        GET_MODULE_HANDLE_EX_FLAG_FROM_ADDRESS | GET_MODULE_HANDLE_EX_FLAG_PIN,
        (LPCWSTR)pinDLL,
        &self
    );
}
#else
static void pinDLL() {}
#endif
*/
import "C"
import (
	"encoding/json"
	"errors"
	"unsafe"

	rewerse "github.com/ByteSizedMarius/rewerse-engineering/pkg"
)

var errNullPointer = errors.New("required parameter is null")

func init() {
	C.pinDLL() // Pin the DLL on Windows to prevent crash at exit
}

// result wraps API responses for consistent JSON structure
type result struct {
	OK    bool   `json:"ok"`
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

func successJSON(data any) *C.char {
	b, err := json.Marshal(result{OK: true, Data: data})
	if err != nil {
		return errorJSON(err)
	}
	return C.CString(string(b))
}

func errorJSON(err error) *C.char {
	b, marshalErr := json.Marshal(result{OK: false, Error: err.Error()})
	if marshalErr != nil {
		// fallback: return a simple error string if marshal fails
		return C.CString(`{"ok":false,"error":"internal error: failed to marshal response"}`)
	}
	return C.CString(string(b))
}

// safeGoString safely converts a C string to Go string, returns empty string for NULL
func safeGoString(s *C.char) string {
	if s == nil {
		return ""
	}
	return C.GoString(s)
}

// requireString checks if a C string is non-NULL and non-empty
func requireString(s *C.char) (string, error) {
	if s == nil {
		return "", errNullPointer
	}
	str := C.GoString(s)
	if str == "" {
		return "", errors.New("required parameter is empty")
	}
	return str, nil
}

//export FreeString
func FreeString(s *C.char) {
	if s == nil {
		return
	}
	C.free(unsafe.Pointer(s))
}

//export SetCertificate
func SetCertificate(certPath, keyPath *C.char) *C.char {
	err := rewerse.SetCertificate(C.GoString(certPath), C.GoString(keyPath))
	if err != nil {
		return errorJSON(err)
	}
	return successJSON(nil)
}

// --- Markets ---

//export MarketSearch
func MarketSearch(query *C.char) *C.char {
	q, err := requireString(query)
	if err != nil {
		return errorJSON(errors.New("query: " + err.Error()))
	}
	markets, err := rewerse.MarketSearch(q)
	if err != nil {
		return errorJSON(err)
	}
	return successJSON(markets)
}

//export GetMarketDetails
func GetMarketDetails(marketID *C.char) *C.char {
	mid, err := requireString(marketID)
	if err != nil {
		return errorJSON(errors.New("marketID: " + err.Error()))
	}
	details, err := rewerse.GetMarketDetails(mid)
	if err != nil {
		return errorJSON(err)
	}
	return successJSON(details)
}

// --- Products ---

// productOptsFromJSON parses optional ProductOpts from JSON string
func productOptsFromJSON(optsJSON *C.char) (*rewerse.ProductOpts, error) {
	if optsJSON == nil {
		return nil, nil
	}
	s := C.GoString(optsJSON)
	if s == "" || s == "null" {
		return nil, nil
	}

	var opts rewerse.ProductOpts
	if err := json.Unmarshal([]byte(s), &opts); err != nil {
		return nil, err
	}
	return &opts, nil
}

//export GetProducts
func GetProducts(marketID, search, optsJSON *C.char) *C.char {
	mid, err := requireString(marketID)
	if err != nil {
		return errorJSON(errors.New("marketID: " + err.Error()))
	}
	opts, err := productOptsFromJSON(optsJSON)
	if err != nil {
		return errorJSON(errors.New("options: " + err.Error()))
	}
	results, err := rewerse.GetProducts(mid, safeGoString(search), opts)
	if err != nil {
		return errorJSON(err)
	}
	return successJSON(results)
}

//export GetCategoryProducts
func GetCategoryProducts(marketID, categorySlug, optsJSON *C.char) *C.char {
	mid, err := requireString(marketID)
	if err != nil {
		return errorJSON(errors.New("marketID: " + err.Error()))
	}
	slug, err := requireString(categorySlug)
	if err != nil {
		return errorJSON(errors.New("categorySlug: " + err.Error()))
	}
	opts, err := productOptsFromJSON(optsJSON)
	if err != nil {
		return errorJSON(errors.New("options: " + err.Error()))
	}
	results, err := rewerse.GetCategoryProducts(mid, slug, opts)
	if err != nil {
		return errorJSON(err)
	}
	return successJSON(results)
}

//export GetProductByID
func GetProductByID(marketID, productID *C.char) *C.char {
	mid, err := requireString(marketID)
	if err != nil {
		return errorJSON(errors.New("marketID: " + err.Error()))
	}
	pid, err := requireString(productID)
	if err != nil {
		return errorJSON(errors.New("productID: " + err.Error()))
	}
	product, err := rewerse.GetProductByID(mid, pid)
	if err != nil {
		return errorJSON(err)
	}
	return successJSON(product)
}

//export GetProductSuggestions
func GetProductSuggestions(query, optsJSON *C.char) *C.char {
	q, err := requireString(query)
	if err != nil {
		return errorJSON(errors.New("query: " + err.Error()))
	}
	opts, err := productOptsFromJSON(optsJSON)
	if err != nil {
		return errorJSON(errors.New("options: " + err.Error()))
	}
	suggestions, err := rewerse.GetProductSuggestions(q, opts)
	if err != nil {
		return errorJSON(err)
	}
	return successJSON(suggestions)
}

//export GetProductRecommendations
func GetProductRecommendations(marketID, listingID *C.char) *C.char {
	mid, err := requireString(marketID)
	if err != nil {
		return errorJSON(errors.New("marketID: " + err.Error()))
	}
	lid, err := requireString(listingID)
	if err != nil {
		return errorJSON(errors.New("listingID: " + err.Error()))
	}
	products, err := rewerse.GetProductRecommendations(mid, lid)
	if err != nil {
		return errorJSON(err)
	}
	return successJSON(products)
}

// --- Discounts ---

//export GetDiscountsRaw
func GetDiscountsRaw(marketID *C.char) *C.char {
	mid, err := requireString(marketID)
	if err != nil {
		return errorJSON(errors.New("marketID: " + err.Error()))
	}
	discounts, err := rewerse.GetDiscountsRaw(mid)
	if err != nil {
		return errorJSON(err)
	}
	return successJSON(discounts)
}

//export GetDiscounts
func GetDiscounts(marketID *C.char) *C.char {
	mid, err := requireString(marketID)
	if err != nil {
		return errorJSON(errors.New("marketID: " + err.Error()))
	}
	discounts, err := rewerse.GetDiscounts(mid)
	if err != nil {
		return errorJSON(err)
	}
	return successJSON(discounts)
}

// --- Recipes ---

// recipeOptsFromJSON parses optional RecipeSearchOpts from JSON string
func recipeOptsFromJSON(optsJSON *C.char) (*rewerse.RecipeSearchOpts, error) {
	if optsJSON == nil {
		return nil, nil
	}
	s := C.GoString(optsJSON)
	if s == "" || s == "null" {
		return nil, nil
	}

	var opts rewerse.RecipeSearchOpts
	if err := json.Unmarshal([]byte(s), &opts); err != nil {
		return nil, err
	}
	return &opts, nil
}

//export RecipeSearch
func RecipeSearch(optsJSON *C.char) *C.char {
	opts, err := recipeOptsFromJSON(optsJSON)
	if err != nil {
		return errorJSON(errors.New("options: " + err.Error()))
	}
	results, err := rewerse.RecipeSearch(opts)
	if err != nil {
		return errorJSON(err)
	}
	return successJSON(results)
}

//export GetRecipeDetails
func GetRecipeDetails(recipeID *C.char) *C.char {
	rid, err := requireString(recipeID)
	if err != nil {
		return errorJSON(errors.New("recipeID: " + err.Error()))
	}
	details, err := rewerse.GetRecipeDetails(rid)
	if err != nil {
		return errorJSON(err)
	}
	return successJSON(details)
}

//export GetRecipePopularTerms
func GetRecipePopularTerms() *C.char {
	terms, err := rewerse.GetRecipePopularTerms()
	if err != nil {
		return errorJSON(err)
	}
	return successJSON(terms)
}

// --- Misc ---

//export GetRecalls
func GetRecalls() *C.char {
	recalls, err := rewerse.GetRecalls()
	if err != nil {
		return errorJSON(err)
	}
	return successJSON(recalls)
}

//export GetRecipeHub
func GetRecipeHub() *C.char {
	hub, err := rewerse.GetRecipeHub()
	if err != nil {
		return errorJSON(err)
	}
	return successJSON(hub)
}

//export GetServicePortfolio
func GetServicePortfolio(zipcode *C.char) *C.char {
	zip, err := requireString(zipcode)
	if err != nil {
		return errorJSON(errors.New("zipcode: " + err.Error()))
	}
	portfolio, err := rewerse.GetServicePortfolio(zip)
	if err != nil {
		return errorJSON(err)
	}
	return successJSON(portfolio)
}

//export GetShopOverview
func GetShopOverview(marketID *C.char) *C.char {
	mid, err := requireString(marketID)
	if err != nil {
		return errorJSON(errors.New("marketID: " + err.Error()))
	}
	overview, err := rewerse.GetShopOverview(mid)
	if err != nil {
		return errorJSON(err)
	}
	return successJSON(overview)
}

// shopOverviewOptsFromJSON parses optional ShopOverviewOpts from JSON string
func shopOverviewOptsFromJSON(optsJSON *C.char) (*rewerse.ShopOverviewOpts, error) {
	if optsJSON == nil {
		return nil, nil
	}
	s := C.GoString(optsJSON)
	if s == "" || s == "null" {
		return nil, nil
	}

	var opts rewerse.ShopOverviewOpts
	if err := json.Unmarshal([]byte(s), &opts); err != nil {
		return nil, err
	}
	return &opts, nil
}

//export GetShopOverviewWithOpts
func GetShopOverviewWithOpts(marketID, optsJSON *C.char) *C.char {
	mid, err := requireString(marketID)
	if err != nil {
		return errorJSON(errors.New("marketID: " + err.Error()))
	}
	opts, err := shopOverviewOptsFromJSON(optsJSON)
	if err != nil {
		return errorJSON(errors.New("options: " + err.Error()))
	}
	overview, err := rewerse.GetShopOverviewWithOpts(mid, opts)
	if err != nil {
		return errorJSON(err)
	}
	return successJSON(overview)
}

// --- Basket ---

// basketSessionResult wraps basket creation response with session info
type basketSessionResult struct {
	BasketID    string `json:"basketId"`
	DeviceID    string `json:"deviceId"`
	Version     int    `json:"version"`
	MarketID    string `json:"marketId"`
	ZipCode     string `json:"zipCode"`
	ServiceType string `json:"serviceType"`
}

//export CreateBasket
func CreateBasket(marketID, zipCode, serviceType *C.char) *C.char {
	mid, err := requireString(marketID)
	if err != nil {
		return errorJSON(errors.New("marketID: " + err.Error()))
	}
	zip, err := requireString(zipCode)
	if err != nil {
		return errorJSON(errors.New("zipCode: " + err.Error()))
	}
	st, err := requireString(serviceType)
	if err != nil {
		return errorJSON(errors.New("serviceType: " + err.Error()))
	}
	session, err := rewerse.CreateBasket(mid, zip, rewerse.ServiceType(st))
	if err != nil {
		return errorJSON(err)
	}
	return successJSON(basketSessionResult{
		BasketID:    session.ID,
		DeviceID:    session.DeviceID,
		Version:     session.Version,
		MarketID:    session.MarketID,
		ZipCode:     session.ZipCode,
		ServiceType: string(session.ServiceType),
	})
}

// reconstructSession rebuilds a BasketSession from individual parameters
func reconstructSession(basketID, marketID, zipCode, serviceType *C.char, version int) (*rewerse.BasketSession, error) {
	bid, err := requireString(basketID)
	if err != nil {
		return nil, errors.New("basketID: " + err.Error())
	}
	mid, err := requireString(marketID)
	if err != nil {
		return nil, errors.New("marketID: " + err.Error())
	}
	zip, err := requireString(zipCode)
	if err != nil {
		return nil, errors.New("zipCode: " + err.Error())
	}
	st, err := requireString(serviceType)
	if err != nil {
		return nil, errors.New("serviceType: " + err.Error())
	}
	return &rewerse.BasketSession{
		ID:          bid,
		MarketID:    mid,
		ZipCode:     zip,
		ServiceType: rewerse.ServiceType(st),
		Version:     version,
	}, nil
}

//export GetBasket
func GetBasket(basketID, marketID, zipCode, serviceType *C.char, version C.int) *C.char {
	session, err := reconstructSession(basketID, marketID, zipCode, serviceType, int(version))
	if err != nil {
		return errorJSON(err)
	}
	basket, err := session.GetBasket()
	if err != nil {
		return errorJSON(err)
	}
	return successJSON(basket)
}

//export SetBasketItemQuantity
func SetBasketItemQuantity(basketID, marketID, zipCode, serviceType, listingID *C.char, quantity, version C.int) *C.char {
	lid, err := requireString(listingID)
	if err != nil {
		return errorJSON(errors.New("listingID: " + err.Error()))
	}
	session, err := reconstructSession(basketID, marketID, zipCode, serviceType, int(version))
	if err != nil {
		return errorJSON(err)
	}
	basket, err := session.SetItemQuantity(lid, int(quantity))
	if err != nil {
		return errorJSON(err)
	}
	return successJSON(basket)
}

//export RemoveBasketItem
func RemoveBasketItem(basketID, marketID, zipCode, serviceType, listingID *C.char) *C.char {
	session, err := reconstructSession(basketID, marketID, zipCode, serviceType, 0)
	if err != nil {
		return errorJSON(err)
	}
	lid, err := requireString(listingID)
	if err != nil {
		return errorJSON(errors.New("listingID: " + err.Error()))
	}
	basket, err := session.RemoveItem(lid)
	if err != nil {
		return errorJSON(err)
	}
	return successJSON(basket)
}

// --- Delivery ---

//export GetBulkyGoodsConfig
func GetBulkyGoodsConfig(marketID, serviceType *C.char) *C.char {
	mid, err := requireString(marketID)
	if err != nil {
		return errorJSON(errors.New("marketID: " + err.Error()))
	}
	st, err := requireString(serviceType)
	if err != nil {
		return errorJSON(errors.New("serviceType: " + err.Error()))
	}
	config, err := rewerse.GetBulkyGoodsConfig(mid, rewerse.ServiceType(st))
	if err != nil {
		return errorJSON(err)
	}
	return successJSON(config)
}

func main() {}
