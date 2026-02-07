# rewerse

Python bindings for the REWE mobile API.

## Notes

- FFI wrapper around compiled Go library ([Github](https://github.com/ByteSizedMarius/rewerse-engineering/)), not native Python
- ~9 MB package size (includes .so and .dll)
- Linux x86_64 and Windows x86_64 only
- Response types are untyped dicts; see [Go structs](https://pkg.go.dev/github.com/ByteSizedMarius/rewerse-engineering/pkg) for field definitions

## Requirements

- Python 3.10+
- mTLS certificates from REWE APK ([extraction instructions](https://github.com/ByteSizedMarius/rewerse-engineering/tree/main/docs))

## Installation

```bash
pip install rewerse
```

## Building from source

Requires Go 1.21+ and a C compiler.

From repo root:
```bash
# Linux
python python/build_lib.py

# Windows (cross-compile from WSL)
python python/build_lib.py --platform windows
```

Then install locally:
```bash
cd python && pip install -e .
```

## Usage

```python
from rewerse import Rewerse

client = Rewerse(cert="certificate.pem", key="private.key")

markets = client.market_search("Berlin")
market_id = markets[0]["wwIdent"]

response = client.get_products(market_id, "Milch")
discounts = client.get_discounts(market_id)
recipes = client.recipe_search(search_term="Pasta")
```

## Methods

All methods raise `RewerseError` on failure. Responses are untyped dicts; see [Go structs](https://pkg.go.dev/github.com/ByteSizedMarius/rewerse-engineering/pkg) for field definitions.

### Markets

| Method | Description |
|--------|-------------|
| `market_search(query)` | Search markets by PLZ, city, street or name |
| `get_market_details(market_id)` | Get market details (hours, address, etc.) |

### Products

| Method | Description |
|--------|-------------|
| `get_products(market_id, search, *, page=1, objects_per_page=30, filters=None)` | Search products in a market |
| `get_category_products(market_id, category_slug, *, page=1, objects_per_page=30, filters=None)` | Get products by category slug |
| `get_product_by_id(market_id, product_id)` | Get product details |
| `get_product_suggestions(query, *, page=1, objects_per_page=25)` | Autocomplete suggestions |
| `get_product_recommendations(market_id, listing_id)` | Related product recommendations |

### Discounts

| Method | Description |
|--------|-------------|
| `get_discounts(market_id)` | Get parsed weekly discounts |
| `get_discounts_raw(market_id)` | Get raw discount data including handout links |

### Recipes

| Method | Description |
|--------|-------------|
| `recipe_search(*, search_term="", sorting="RELEVANCE_DESC", difficulty="", collection="", page=1, objects_per_page=20)` | Search recipes. Difficulty: `"Gering"`, `"Mittel"`, `"Hoch"`. Collection: `"Vegetarisch"`, `"Vegan"` |
| `get_recipe_details(recipe_id)` | Get full recipe with ingredients and steps |
| `get_recipe_popular_terms()` | Get popular recipe search terms |
| `get_recipe_hub()` | Get recipe homepage (recipe of the day, popular, categories) |

### Shop & Services

| Method | Description |
|--------|-------------|
| `get_shop_overview(market_id, *, service_type="PICKUP", zip_code="")` | Get product categories for a market's online shop |
| `get_service_portfolio(zipcode)` | Get available REWE services for a zip code |
| `get_recalls()` | Get current product recalls |

### Basket

| Method | Description |
|--------|-------------|
| `create_basket(market_id, zip_code, service_type="PICKUP")` | Create a new shopping basket |
| `get_basket(basket_id, market_id, zip_code, service_type, version=0)` | Get current basket state |
| `set_basket_item(basket_id, market_id, zip_code, service_type, listing_id, quantity, version)` | Add/update item quantity in basket |
| `remove_basket_item(basket_id, market_id, zip_code, service_type, listing_id)` | Remove item from basket |

### Delivery

| Method | Description |
|--------|-------------|
| `get_bulky_goods_config(market_id, service_type="DELIVERY")` | Get beverage crate limits and surcharges |

## License

MIT
