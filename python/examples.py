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

    if not markets:
        print("No markets found, exiting")
        return

    # Use first market for subsequent examples
    market_id = markets[0]["wwIdent"]
    print(f"Using market: {markets[0]['name']} ({market_id})")

    # Get detailed market info (hours, address, services)
    details = client.get_market_details(market_id)
    market = details["Market"]
    print(f"Address: {market['street']}, {market['zipCode']} {market['city']}")

    # --- Product Search ---

    # Basic product search (response wrapped in data.products)
    response = client.get_products(market_id, "Milch", page=1, objects_per_page=5)
    products_data = response["data"]["products"]
    print(f"\nSearch 'Milch': {products_data['pagination']['objectCount']} results")
    for p in products_data["products"][:3]:
        # Price is in cents
        price = p["listing"]["currentRetailPrice"] / 100
        print(f"  - {p['title']}: {price:.2f}€")

    # Search with filters
    response = client.get_products(market_id, "Joghurt", filters=["attribute=vegan"])
    products_data = response["data"]["products"]
    print(f"\nVegan Joghurt: {products_data['pagination']['objectCount']} results")

    # Get product details by ID (uses productId, not listingId)
    if products_data["products"]:
        product_id = products_data["products"][0]["productId"]
        listing_id = products_data["products"][0]["listing"]["listingId"]
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
        response = client.get_category_products(market_id, slug, page=1, objects_per_page=5)
        cat_data = response["data"]["products"]
        print(f"\nProducts in '{categories[0]['name']}': {cat_data['pagination']['objectCount']}")

    # --- Discounts ---

    discounts = client.get_discounts(market_id)
    print(f"\nDiscounts valid until: {discounts['ValidUntil']}")
    for cat in discounts["Categories"][:2]:
        print(f"  {cat['Title']}: {len(cat['Offers'])} offers")

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
