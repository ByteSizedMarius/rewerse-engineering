package rewerse

import (
	"os"
	"testing"
)

var certLoaded = false

const (
	testMarketID = "831002" // REWE Ludwigshafen/Mundenheim (PICKUP enabled for zip 67065)
	testZipCode  = "67065"
)

func TestMain(m *testing.M) {
	// Try to load certificates from project root
	err := SetCertificate("../certificate.pem", "../private.key")
	if err == nil {
		certLoaded = true
	} else {
		println("NOTE: Certificates not found (certificate.pem, private.key in project root)")
		println("      API tests will be skipped. See docs/ for extraction instructions.")
	}
	os.Exit(m.Run())
}

func skipIfNoCert(t *testing.T) {
	if !certLoaded {
		t.Skip("certificate not loaded, skipping integration test")
	}
}

func TestMarketSearch(t *testing.T) {
	skipIfNoCert(t)

	markets, err := MarketSearch("Mannheim")
	if err != nil {
		t.Fatalf("MarketSearch failed: %v", err)
	}
	if len(markets) == 0 {
		t.Fatal("expected at least one market")
	}
	if markets[0].WWIdent == "" {
		t.Error("market WWIdent is empty")
	}
	if markets[0].Name == "" {
		t.Error("market Name is empty")
	}
}

func TestGetMarketDetails(t *testing.T) {
	skipIfNoCert(t)

	md, err := GetMarketDetails(testMarketID)
	if err != nil {
		t.Fatalf("GetMarketDetails failed: %v", err)
	}
	if md.Market.WWIdent != testMarketID {
		t.Errorf("expected WWIdent %s, got %s", testMarketID, md.Market.WWIdent)
	}
	if md.Market.Street == "" {
		t.Error("market Street is empty")
	}
	if md.Market.City == "" {
		t.Error("market City is empty")
	}
}

func TestGetProducts(t *testing.T) {
	skipIfNoCert(t)

	results, err := GetProducts(testMarketID, "Karotten", nil)
	if err != nil {
		t.Fatalf("GetProducts failed: %v", err)
	}
	if len(results.Products) == 0 {
		t.Fatal("expected at least one product")
	}
	p := results.Products[0]
	if p.ProductID == "" {
		t.Error("product ID is empty")
	}
	if p.Title == "" {
		t.Error("product Title is empty")
	}
}

func TestGetProductsWithPagination(t *testing.T) {
	skipIfNoCert(t)

	opts := &ProductOpts{Page: 1, ObjectsPerPage: 5}
	results, err := GetProducts(testMarketID, "Milch", opts)
	if err != nil {
		t.Fatalf("GetProducts with pagination failed: %v", err)
	}
	// API may return more products than requested due to bundling/ads
	if len(results.Products) == 0 {
		t.Fatal("expected at least one product")
	}
	t.Logf("requested 5, got %d products (API may include extras)", len(results.Products))
}

func TestGetProductsWithFilters(t *testing.T) {
	skipIfNoCert(t)

	opts := &ProductOpts{Filters: []ProductFilter{FilterVegan}}
	results, err := GetProducts(testMarketID, "", opts)
	if err != nil {
		t.Fatalf("GetProducts with filter failed: %v", err)
	}
	// Just verify we got a response, filter correctness depends on API
	if results.Pagination.ObjectCount == 0 {
		t.Log("no vegan products returned (may be valid)")
	}
}

func TestGetCategoryProducts(t *testing.T) {
	skipIfNoCert(t)

	results, err := GetCategoryProducts(testMarketID, "obst-gemuese", nil)
	if err != nil {
		t.Fatalf("GetCategoryProducts failed: %v", err)
	}
	if len(results.Products) == 0 {
		t.Fatal("expected at least one product in category")
	}
}

func TestGetProductByID(t *testing.T) {
	skipIfNoCert(t)

	// First get a product ID from search
	results, err := GetProducts(testMarketID, "Milch", nil)
	if err != nil {
		t.Fatalf("GetProducts failed: %v", err)
	}
	if len(results.Products) == 0 {
		t.Skip("no products found to test GetProductByID")
	}

	productID := results.Products[0].ProductID
	product, err := GetProductByID(testMarketID, productID)
	if err != nil {
		t.Fatalf("GetProductByID failed: %v", err)
	}
	if product.ProductID != productID {
		t.Errorf("expected product ID %s, got %s", productID, product.ProductID)
	}
	if product.Title == "" {
		t.Error("product Title is empty")
	}
}

func TestGetProductSuggestions(t *testing.T) {
	skipIfNoCert(t)

	suggestions, err := GetProductSuggestions("Milch", nil)
	if err != nil {
		t.Fatalf("GetProductSuggestions failed: %v", err)
	}
	if len(suggestions) == 0 {
		t.Fatal("expected at least one suggestion")
	}
	if suggestions[0].Title == "" {
		t.Error("suggestion Title is empty")
	}
}

func TestGetProductRecommendations(t *testing.T) {
	skipIfNoCert(t)

	// First get a listing ID from product search
	results, err := GetProducts(testMarketID, "Brot", nil)
	if err != nil {
		t.Fatalf("GetProducts failed: %v", err)
	}
	if len(results.Products) == 0 {
		t.Skip("no products found to test recommendations")
	}

	listingID := results.Products[0].Listing.ListingID
	if listingID == "" {
		t.Skip("no listing ID in product")
	}

	// Recommendations may return empty, that's valid
	recs, err := GetProductRecommendations(testMarketID, listingID)
	if err != nil {
		t.Fatalf("GetProductRecommendations failed: %v", err)
	}
	t.Logf("got %d recommendations", len(recs))
}

func TestRecipeSearch(t *testing.T) {
	skipIfNoCert(t)

	results, err := RecipeSearch(nil)
	if err != nil {
		t.Fatalf("RecipeSearch failed: %v", err)
	}
	if results.TotalCount == 0 {
		t.Fatal("expected at least one recipe")
	}
	if len(results.Recipes) == 0 {
		t.Fatal("expected at least one recipe in results")
	}
	if results.Recipes[0].Title == "" {
		t.Error("recipe Title is empty")
	}
}

func TestRecipeSearchWithFilters(t *testing.T) {
	skipIfNoCert(t)

	opts := &RecipeSearchOpts{
		Collection: CollectionVegetarisch,
		Difficulty: DifficultyEasy,
	}
	results, err := RecipeSearch(opts)
	if err != nil {
		t.Fatalf("RecipeSearch with filters failed: %v", err)
	}
	// May return fewer results with filters, that's ok
	t.Logf("got %d vegetarisch easy recipes", results.TotalCount)
}

func TestGetRecipeDetails(t *testing.T) {
	skipIfNoCert(t)

	// First get a recipe ID from search
	results, err := RecipeSearch(nil)
	if err != nil {
		t.Fatalf("RecipeSearch failed: %v", err)
	}
	if len(results.Recipes) == 0 {
		t.Skip("no recipes found to test details")
	}

	recipeID := results.Recipes[0].ID
	details, err := GetRecipeDetails(recipeID)
	if err != nil {
		t.Fatalf("GetRecipeDetails failed: %v", err)
	}
	if details.Recipe.ID != recipeID {
		t.Errorf("expected recipe ID %s, got %s", recipeID, details.Recipe.ID)
	}
	if details.Recipe.Title == "" {
		t.Error("recipe Title is empty")
	}
	if len(details.Recipe.Steps) == 0 {
		t.Error("recipe has no steps")
	}
}

func TestGetRecipePopularTerms(t *testing.T) {
	skipIfNoCert(t)

	terms, err := GetRecipePopularTerms()
	if err != nil {
		t.Fatalf("GetRecipePopularTerms failed: %v", err)
	}
	if len(terms) == 0 {
		t.Fatal("expected at least one popular term")
	}
	if terms[0].Title == "" {
		t.Error("term Title is empty")
	}
}

func TestGetRecipeHub(t *testing.T) {
	skipIfNoCert(t)

	hub, err := GetRecipeHub()
	if err != nil {
		t.Fatalf("GetRecipeHub failed: %v", err)
	}
	if hub.RecipeOfTheDay.Title == "" {
		t.Error("recipe of the day Title is empty")
	}
	if len(hub.PopularRecipes) == 0 {
		t.Error("no popular recipes")
	}
}

func TestGetRecalls(t *testing.T) {
	skipIfNoCert(t)

	recalls, err := GetRecalls()
	if err != nil {
		t.Fatalf("GetRecalls failed: %v", err)
	}
	// May be empty, that's valid
	t.Logf("got %d recalls", len(recalls))
}

func TestGetDiscountsRaw(t *testing.T) {
	skipIfNoCert(t)

	rd, err := GetDiscountsRaw(testMarketID)
	if err != nil {
		t.Fatalf("GetDiscountsRaw failed: %v", err)
	}
	// Check we got some data
	if rd.Data.Offers.Current.UntilDate == "" && rd.Data.Offers.Next.UntilDate == "" {
		t.Error("no until date in offers")
	}
}

func TestGetDiscounts(t *testing.T) {
	skipIfNoCert(t)

	ds, err := GetDiscounts(testMarketID)
	if err != nil {
		t.Fatalf("GetDiscounts failed: %v", err)
	}
	if len(ds.Categories) == 0 {
		t.Fatal("expected at least one discount category")
	}
	if ds.ValidUntil.IsZero() {
		t.Error("ValidUntil is zero")
	}
}

func TestGetShopOverview(t *testing.T) {
	skipIfNoCert(t)

	so, err := GetShopOverview(testMarketID)
	if err != nil {
		t.Fatalf("GetShopOverview failed: %v", err)
	}
	// Some markets may not have pickup service
	if len(so.ProductCategories) == 0 {
		t.Log("no product categories (market may not support PICKUP)")
	} else {
		if so.ProductCategories[0].Name == "" {
			t.Error("first category Name is empty")
		}
	}
}

func TestGetServicePortfolio(t *testing.T) {
	skipIfNoCert(t)

	sp, err := GetServicePortfolio(testZipCode)
	if err != nil {
		t.Fatalf("GetServicePortfolio failed: %v", err)
	}
	if sp.CustomerZipCode != testZipCode {
		t.Errorf("expected zip %s, got %s", testZipCode, sp.CustomerZipCode)
	}
}

func TestGetBulkyGoodsConfig(t *testing.T) {
	skipIfNoCert(t)

	// Get a delivery market first
	sp, err := GetServicePortfolio("50667")
	if err != nil {
		t.Fatalf("GetServicePortfolio failed: %v", err)
	}
	if sp.DeliveryMarket == nil {
		t.Skip("no delivery market for zip 50667")
	}

	config, err := GetBulkyGoodsConfig(sp.DeliveryMarket.WWIdent, ServiceDelivery)
	if err != nil {
		t.Fatalf("GetBulkyGoodsConfig failed: %v", err)
	}
	t.Logf("hasBeverageSurcharge: %v", config.HasBeverageSurcharge)
	if config.HasBeverageSurcharge && config.BeverageSurcharge != nil {
		t.Logf("softLimit: %d, hardLimit: %d, surcharge: %d",
			config.BeverageSurcharge.SoftLimit,
			config.BeverageSurcharge.HardLimit,
			config.BeverageSurcharge.Surcharge)
	}
}

func TestCreateBasket(t *testing.T) {
	skipIfNoCert(t)

	session, err := CreateBasket(testMarketID, testZipCode, ServicePickup)
	if err != nil {
		t.Fatalf("CreateBasket failed: %v", err)
	}
	if session.ID == "" {
		t.Error("basket ID is empty")
	}
	if session.DeviceID == "" {
		t.Error("device ID is empty")
	}
	t.Logf("created basket: id=%s, deviceId=%s, version=%d",
		session.ID, session.DeviceID, session.Version)
}

func TestBasketFullFlow(t *testing.T) {
	skipIfNoCert(t)

	// Create basket
	session, err := CreateBasket(testMarketID, testZipCode, ServicePickup)
	if err != nil {
		t.Fatalf("CreateBasket failed: %v", err)
	}
	t.Logf("created basket: %s (v%d)", session.ID, session.Version)

	// Get a product listing ID
	results, err := GetProducts(testMarketID, "Milch", nil)
	if err != nil {
		t.Fatalf("GetProducts failed: %v", err)
	}
	if len(results.Products) == 0 {
		t.Skip("no products found")
	}
	listingID := results.Products[0].Listing.ListingID
	if listingID == "" {
		t.Skip("product has no listing ID")
	}
	t.Logf("using listing: %s", listingID)

	// Add item to basket
	basket, err := session.SetItemQuantity(listingID, 2)
	if err != nil {
		t.Fatalf("SetItemQuantity failed: %v", err)
	}
	if len(basket.LineItems) == 0 {
		t.Fatal("expected item in basket")
	}
	if basket.LineItems[0].Quantity != 2 {
		t.Errorf("expected quantity 2, got %d", basket.LineItems[0].Quantity)
	}
	t.Logf("added item: qty=%d, price=%d, v%d",
		basket.LineItems[0].Quantity, basket.LineItems[0].TotalPrice, basket.Version)

	// Update quantity
	basket, err = session.SetItemQuantity(listingID, 1)
	if err != nil {
		t.Fatalf("SetItemQuantity (update) failed: %v", err)
	}
	if basket.LineItems[0].Quantity != 1 {
		t.Errorf("expected quantity 1, got %d", basket.LineItems[0].Quantity)
	}
	t.Logf("updated item: qty=%d, v%d", basket.LineItems[0].Quantity, basket.Version)

	// Remove item
	basket, err = session.RemoveItem(listingID)
	if err != nil {
		t.Fatalf("RemoveItem failed: %v", err)
	}
	if len(basket.LineItems) != 0 {
		t.Errorf("expected empty basket, got %d items", len(basket.LineItems))
	}
	t.Logf("removed item: basket now has %d items, v%d", len(basket.LineItems), basket.Version)
}

func TestGetShopOverviewDelivery(t *testing.T) {
	skipIfNoCert(t)

	// Get a delivery market
	sp, err := GetServicePortfolio("50667")
	if err != nil {
		t.Fatalf("GetServicePortfolio failed: %v", err)
	}
	if sp.DeliveryMarket == nil {
		t.Skip("no delivery market for zip 50667")
	}

	opts := &ShopOverviewOpts{
		ServiceType: ServiceDelivery,
		ZipCode:     "50667",
	}
	so, err := GetShopOverviewWithOpts(sp.DeliveryMarket.WWIdent, opts)
	if err != nil {
		t.Fatalf("GetShopOverviewWithOpts (DELIVERY) failed: %v", err)
	}
	if len(so.ProductCategories) == 0 {
		t.Log("no product categories for delivery")
	} else {
		t.Logf("got %d categories for delivery", len(so.ProductCategories))
	}
}
