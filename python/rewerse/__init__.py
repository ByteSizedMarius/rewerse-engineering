"""
Rewerse - Python bindings for the REWE mobile API library.

Example:
    from rewerse import Rewerse

    client = Rewerse(cert="client.pem", key="client.key")
    markets = client.market_search("Berlin")
    for market in markets:
        print(market["name"])
"""

from __future__ import annotations

import json

from ._ffi import call, RewerseError

__version__ = "0.1.0"
__all__ = ["Rewerse", "RewerseError", "__version__"]


class Rewerse:
    """
    Client for the REWE mobile API.

    Requires mTLS certificates for authentication.
    """

    def __init__(self, cert: str, key: str):
        """
        Initialize the client with mTLS certificates.

        Args:
            cert: Path to the client certificate file (.pem)
            key: Path to the client key file (.pem/.key)

        Raises:
            RewerseError: If certificates cannot be loaded
        """
        call("SetCertificate", cert, key)

    # --- Markets ---

    def market_search(self, query: str) -> list[dict]:
        """
        Search for REWE markets.

        Args:
            query: Search query (PLZ, city, street, market name, etc.)

        Returns:
            List of matching markets
        """
        return call("MarketSearch", query)

    def get_market_details(self, market_id: str) -> dict:
        """
        Get detailed information about a market.

        Args:
            market_id: The market ID (e.g., "831002")

        Returns:
            Market details including opening hours, address, etc.
        """
        return call("GetMarketDetails", market_id)

    # --- Products ---

    def get_products(
        self,
        market_id: str,
        search: str,
        *,
        page: int = 1,
        objects_per_page: int = 30,
        filters: list[str] | None = None,
    ) -> dict:
        """
        Search for products in a market.

        Args:
            market_id: The market ID
            search: Search query
            page: Page number (1-indexed)
            objects_per_page: Results per page
            filters: Product filters (e.g., ["attribute=vegan", "attribute=organic"])

        Returns:
            Product search results with pagination info
        """
        opts = {
            "Page": page,
            "ObjectsPerPage": objects_per_page,
            "Filters": filters or [],
        }
        return call("GetProducts", market_id, search, json.dumps(opts))

    def get_category_products(
        self,
        market_id: str,
        category_slug: str,
        *,
        page: int = 1,
        objects_per_page: int = 30,
        filters: list[str] | None = None,
    ) -> dict:
        """
        Get products in a specific category.

        Args:
            market_id: The market ID
            category_slug: Category slug from shop overview
            page: Page number (1-indexed)
            objects_per_page: Results per page
            filters: Product filters

        Returns:
            Product search results
        """
        opts = {
            "Page": page,
            "ObjectsPerPage": objects_per_page,
            "Filters": filters or [],
        }
        return call("GetCategoryProducts", market_id, category_slug, json.dumps(opts))

    def get_product_by_id(self, market_id: str, product_id: str) -> dict:
        """
        Get detailed product information.

        Args:
            market_id: The market ID
            product_id: The product ID

        Returns:
            Full product details
        """
        return call("GetProductByID", market_id, product_id)

    def get_product_suggestions(
        self,
        query: str,
        *,
        page: int = 1,
        objects_per_page: int = 25,
    ) -> list[dict]:
        """
        Get search autocomplete suggestions.

        Args:
            query: Partial search query
            page: Page number
            objects_per_page: Results per page

        Returns:
            List of suggestion objects
        """
        opts = {"Page": page, "ObjectsPerPage": objects_per_page}
        return call("GetProductSuggestions", query, json.dumps(opts))

    def get_product_recommendations(self, market_id: str, listing_id: str) -> list[dict]:
        """
        Get related product recommendations.

        Args:
            market_id: The market ID
            listing_id: The listing ID of the source product

        Returns:
            List of recommended products
        """
        return call("GetProductRecommendations", market_id, listing_id)

    # --- Discounts ---

    def get_discounts_raw(self, market_id: str) -> dict:
        """
        Get raw discount data from the API.

        Args:
            market_id: The market ID

        Returns:
            Raw discount response including handout links
        """
        return call("GetDiscountsRaw", market_id)

    def get_discounts(self, market_id: str) -> dict:
        """
        Get parsed discount data.

        Args:
            market_id: The market ID

        Returns:
            Discount categories with offers, valid_until date
        """
        return call("GetDiscounts", market_id)

    # --- Recipes ---

    def recipe_search(
        self,
        *,
        search_term: str = "",
        sorting: str = "RELEVANCE_DESC",
        difficulty: str = "",
        collection: str = "",
        page: int = 1,
        objects_per_page: int = 20,
    ) -> dict:
        """
        Search for recipes.

        Args:
            search_term: Search query
            sorting: Sort order (default: "RELEVANCE_DESC")
            difficulty: Filter by difficulty ("Gering", "Mittel", "Hoch")
            collection: Filter by collection ("Vegetarisch", "Vegan")
            page: Page number
            objects_per_page: Results per page

        Returns:
            Recipe search results with metadata
        """
        opts = {
            "SearchTerm": search_term,
            "Sorting": sorting,
            "Difficulty": difficulty,
            "Collection": collection,
            "Page": page,
            "ObjectsPerPage": objects_per_page,
        }
        return call("RecipeSearch", json.dumps(opts))

    def get_recipe_details(self, recipe_id: str) -> dict:
        """
        Get full recipe with ingredients and steps.

        Args:
            recipe_id: The recipe UUID

        Returns:
            Full recipe details
        """
        return call("GetRecipeDetails", recipe_id)

    def get_recipe_popular_terms(self) -> list[dict]:
        """
        Get popular recipe search terms.

        Returns:
            List of popular search terms
        """
        return call("GetRecipePopularTerms")

    # --- Misc ---

    def get_recalls(self) -> list[dict]:
        """
        Get current product recalls.

        Returns:
            List of active product recalls
        """
        return call("GetRecalls")

    def get_recipe_hub(self) -> dict:
        """
        Get the recipe hub homepage data.

        Returns:
            Recipe of the day, popular recipes, categories
        """
        return call("GetRecipeHub")

    def get_service_portfolio(self, zipcode: str) -> dict:
        """
        Get available REWE services for a zip code.

        Args:
            zipcode: German postal code

        Returns:
            Delivery and pickup availability
        """
        return call("GetServicePortfolio", zipcode)

    def get_shop_overview(
        self,
        market_id: str,
        *,
        service_type: str = "PICKUP",
        zip_code: str = "",
    ) -> dict:
        """
        Get product categories for a market's online shop.

        Args:
            market_id: The market ID
            service_type: "PICKUP" or "DELIVERY" (default: PICKUP)
            zip_code: Required for DELIVERY mode

        Returns:
            Shop category structure
        """
        if service_type == "PICKUP" and not zip_code:
            return call("GetShopOverview", market_id)
        opts = {"ServiceType": service_type, "ZipCode": zip_code}
        return call("GetShopOverviewWithOpts", market_id, json.dumps(opts))

    # --- Basket ---

    def create_basket(
        self,
        market_id: str,
        zip_code: str,
        service_type: str = "PICKUP",
    ) -> dict:
        """
        Create a new shopping basket.

        Args:
            market_id: The market ID
            zip_code: Customer's postal code
            service_type: "PICKUP" or "DELIVERY"

        Returns:
            Basket session info with basketId, deviceId, version
        """
        return call("CreateBasket", market_id, zip_code, service_type)

    def get_basket(
        self,
        basket_id: str,
        market_id: str,
        zip_code: str,
        service_type: str,
        version: int = 0,
    ) -> dict:
        """
        Get current basket state.

        Args:
            basket_id: The basket ID from create_basket
            market_id: The market ID
            zip_code: Customer's postal code
            service_type: "PICKUP" or "DELIVERY"
            version: Basket version for optimistic locking (default: 0)

        Returns:
            Full basket with lineItems, summary, violations
        """
        return call("GetBasket", basket_id, market_id, zip_code, service_type, version)

    def set_basket_item(
        self,
        basket_id: str,
        market_id: str,
        zip_code: str,
        service_type: str,
        listing_id: str,
        quantity: int,
        version: int,
    ) -> dict:
        """
        Set item quantity in basket (add/update).

        Args:
            basket_id: The basket ID
            market_id: The market ID
            zip_code: Customer's postal code
            service_type: "PICKUP" or "DELIVERY"
            listing_id: Product listing ID from search results
            quantity: Desired quantity (0 removes the item)
            version: Current basket version for optimistic locking

        Returns:
            Updated basket state
        """
        return call(
            "SetBasketItemQuantity",
            basket_id, market_id, zip_code, service_type,
            listing_id, quantity, version
        )

    def remove_basket_item(
        self,
        basket_id: str,
        market_id: str,
        zip_code: str,
        service_type: str,
        listing_id: str,
    ) -> dict:
        """
        Remove item from basket.

        Args:
            basket_id: The basket ID
            market_id: The market ID
            zip_code: Customer's postal code
            service_type: "PICKUP" or "DELIVERY"
            listing_id: Product listing ID to remove

        Returns:
            Updated basket state
        """
        return call(
            "RemoveBasketItem",
            basket_id, market_id, zip_code, service_type, listing_id
        )

    # --- Delivery ---

    def get_bulky_goods_config(
        self,
        market_id: str,
        service_type: str = "DELIVERY",
    ) -> dict:
        """
        Get beverage crate limits and surcharges for a market.

        Args:
            market_id: The market ID
            service_type: "PICKUP" or "DELIVERY"

        Returns:
            Beverage surcharge config (limits, prices)
        """
        return call("GetBulkyGoodsConfig", market_id, service_type)
