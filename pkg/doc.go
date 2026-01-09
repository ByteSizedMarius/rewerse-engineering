// Package rewerse provides a Go client for the REWE mobile API.
//
// Before making any API calls, you must initialize the client with a valid
// certificate using SetCertificate:
//
//	err := rewerse.SetCertificate("certificate.pem", "private.key")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// The package is safe for concurrent use after initialization. All API functions
// are thread-safe once SetCertificate has been called.
//
// # Markets
//
// Search for markets by location or get market details:
//
//	markets, _ := rewerse.MarketSearch("Mannheim")
//	details, _ := rewerse.GetMarketDetails("840174")
//
// # Products
//
// Search products, browse categories, or get product details:
//
//	results, _ := rewerse.GetProducts("840174", "Karotten", nil)
//	products, _ := rewerse.GetCategoryProducts("840174", "obst-gemuese", nil)
//	product, _ := rewerse.GetProductByID("840174", "9900011")
//
// # Discounts
//
// Get current weekly discounts for a market:
//
//	discounts, _ := rewerse.GetDiscounts("840174")
//
// # Recipes
//
// Search recipes or get recipe details:
//
//	results, _ := rewerse.RecipeSearch(&rewerse.RecipeSearchOpts{SearchTerm: "Pasta"})
//	details, _ := rewerse.GetRecipeDetails("recipe-uuid")
//
// # Basket
//
// Create and manage shopping baskets:
//
//	session, _ := rewerse.CreateBasket("840174", "67065", rewerse.ServicePickup)
//	session.SetItemQuantity("listing-id", 2)
//	basket, _ := session.GetBasket()
package rewerse
