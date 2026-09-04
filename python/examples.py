"""
Usage examples for the rewerse Python library.

Prerequisites:
- Valid mTLS certificates (certificate.pem, private.key)
- Library built: python build_lib.py
"""

from rewerse import Rewerse, RewerseError


def main():
    # Initialize client with certificates
    client = Rewerse(cert="certificate.pem", key="private.key")

    # --- Market Operations ---

    # Search for markets by city, PLZ, or street
    markets = client.market_search("Berlin")
    print(f"Found {len(markets)} markets in Berlin")

    # Not every market supports pickup, so take one from the service portfolio
    portfolio = client.get_service_portfolio("67065")
    if not portfolio["pickupMarkets"]:
        print("No pickup markets available, exiting")
        return

    market_id = portfolio["pickupMarkets"][0]["wwIdent"]
    print(f"Using market: {portfolio['pickupMarkets'][0]['displayName']} ({market_id})")

    # Get detailed market info (hours, address, services)
    details = client.get_market_details(market_id)
    market = details["Market"]
    print(f"Address: {market['street']}, {market['zipCode']} {market['city']}")

    # --- Product Search ---

    # Basic product search
    results = client.get_products(market_id, "Milch", page=1, objects_per_page=5)
    print(f"\nSearch 'Milch': {results['pagination']['objectCount']} results")
    for p in results["products"][:3]:
        # Price is in cents
        price = p["listing"]["currentRetailPrice"] / 100
        print(f"  - {p['title']}: {price:.2f}€")

    # Search with filters
    results = client.get_products(market_id, "Joghurt", filters=["attribute=vegan"])
    print(f"\nVegan Joghurt: {results['pagination']['objectCount']} results")

    # Get product details by ID (uses productId, not listingId)
    if results["products"]:
        product_id = results["products"][0]["productId"]
        listing_id = results["products"][0]["listing"]["listingId"]
        product = client.get_product_by_id(market_id, product_id)
        print(f"\nProduct details: {product['title']}")

        # Get recommendations for a product (uses listingId)
        recs = client.get_product_recommendations(market_id, listing_id)
        print(f"Recommendations: {len(recs)} products")

    # --- Category Browsing ---

    # Get shop category structure
    overview = client.get_shop_overview(market_id)
    categories = overview["productCategories"]
    print(f"\nShop has {len(categories)} top-level categories")
    for cat in categories[:5]:
        print(f"  - {cat['name']} ({cat['slug']})")

    # Get products in a category
    if categories:
        slug = categories[0]["slug"]
        cat_results = client.get_category_products(market_id, slug, page=1, objects_per_page=5)
        print(f"\nProducts in '{categories[0]['name']}': {cat_results['pagination']['objectCount']}")

    # --- Discounts ---

    discounts = client.get_discounts(market_id)
    print(f"\nDiscounts valid until: {discounts['validUntil']}")
    for cat in discounts["categories"][:2]:
        print(f"  {cat['title']}: {len(cat['offers'])} offers")

    # --- Recipes ---

    # Search terms the app shows on the recipe landing page
    terms = client.get_recipe_popular_terms()
    print(f"\nPopular recipe terms: {', '.join(terms[:5])}")

    # Recipe search, filtered to easy fish recipes
    recipes = client.recipe_search(
        search_term="Lachs",
        objects_per_page=5,
        collections=["Fisch"],
        difficulties=[1],
    )
    print(f"Easy fish recipes for 'Lachs': {recipes['metadata']['totalRecipeCount']}")
    for r in recipes["recipes"][:3]:
        print(f"  - {r['title']}")

    # Full recipe with ingredients and steps
    if recipes["recipes"]:
        recipe = client.get_recipe_details(recipes["recipes"][0]["id"])
        print(f"\n{recipe['title']}: {len(recipe['ingredients'])} ingredients")
        for i in recipe["ingredients"][:3]:
            print(f"  - {i['displayName']}")

    # --- Misc ---

    # Check service availability for a postal code
    services = client.get_service_portfolio("10115")
    print(f"\nServices in 10115:")
    has_delivery = services.get("deliveryMarket") is not None
    has_pickup = len(services.get("pickupMarkets", [])) > 0
    print(f"  Delivery: {has_delivery}")
    print(f"  Pickup: {has_pickup} ({len(services.get('pickupMarkets', []))} markets)")

    # Product recalls
    recalls = client.get_recalls()
    print(f"\nActive recalls: {len(recalls)}")


def error_handling_example():
    """Demonstrates error handling."""
    client = Rewerse(cert="certificate.pem", key="private.key")

    try:
        # Invalid product ID raises RewerseError
        client.get_product_by_id("8534540", "nonexistent-product")
    except RewerseError as e:
        print(f"Expected error (invalid product): {e}")

    try:
        # Invalid market ID
        client.get_discounts("0000000")
    except RewerseError as e:
        print(f"Expected error (invalid market): {e}")


if __name__ == "__main__":
    main()
    print("\n--- Error Handling ---")
    error_handling_example()
