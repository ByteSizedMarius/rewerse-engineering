# rewerse

Python bindings for the REWE mobile API.

## Notes

- FFI wrapper around compiled Go library ([Github](https://github.com/ByteSizedMarius/rewerse-engineering/)), not native Python
- ~9 MB package size (includes .so and .dll)
- Linux x86_64 and Windows x86_64 only (macOS: [build from source](#macos))
- Response types are untyped dicts; see [Go structs](https://pkg.go.dev/github.com/ByteSizedMarius/rewerse-engineering/pkg) for field definitions

## Requirements

- Python 3.10+
- mTLS certificates from REWE APK ([extraction instructions](https://github.com/ByteSizedMarius/rewerse-engineering/tree/main/docs))

## Installation

```bash
pip install rewerse
```

Alternatively, build from source on Windows (requires Go 1.21+ and gcc):

```bash
python python/build_lib.py
cd python && pip install -e .
```

This builds both the Linux .so and Windows .dll. For macOS, see [building from source](#macos).

## Usage

```python
from rewerse import Rewerse

client = Rewerse(cert="certificate.pem", key="private.key")

markets = client.market_search("Berlin")
# [
#     {
#         "wwIdent": "8547534",
#         "name": "REWE Markt",
#         "companyName": "REWE Karsten Schmidt oHG",
#         "street": "Schönhauser Allee 80",
#         "zipCode": "10439",
#         "city": "Berlin",
#         "location": {"latitude": 52.54984, "longitude": 13.41488},
#         "openingStatus": {"openState": "CLOSED", "statusText": "Geschlossen", ...},
#         "serviceFlags": {"hasPickup": true},
#         ...
#     },
#     ...
# ]

market_id = markets[0]["wwIdent"]

products = client.get_products(market_id, "Milch")
# {
#     "products": [
#         {
#             "productId": "677067",
#             "title": "Weihenstephan H-Milch 3,5% 1l",
#             "listing": {
#                 "currentRetailPrice": 99,
#                 "grammage": "1l",
#                 "discount": {"__typename": "RegularProductDiscount", "validTo": "04.04."},
#                 "loyaltyBonus": {"bonusType": "CENT", "bonusValue": 10},
#                 ...
#             },
#             "attributes": {"isOrganic": false, "isVegan": false, "isBulkyGood": true, ...},
#             ...
#         },
#         ...
#     ],
#     "pagination": {"objectsPerPage": 30, "currentPage": 1, "pageCount": 24, "objectCount": 706},
#     "searchTerm": {"original": "Milch", "corrected": null}
# }

discounts = client.get_discounts(market_id)
# {
#     "categories": [
#         {
#             "id": "markt-topangebote",
#             "title": "Top-Angebote in deinem Markt",
#             "index": 3,
#             "offers": [
#                 {
#                     "title": "Kinder Riegel",
#                     "subtitle": "je 18 x 21-g-Pckg. (1 kg = 11.56)",
#                     "price": 3.99,
#                     "priceRaw": "3,99 €",
#                     "loyaltyBonus": 0.1,
#                     "manufacturer": "FERRERO",
#                     "articleNo": "249183",
#                     "productCategory": "suesses-und-salziges",
#                     ...
#                 },
#                 ...
#             ]
#         },
#         ...
#     ],
#     "validUntil": "2026-04-05T00:00:00Z"
# }

```

## Methods

All methods raise `RewerseError` on failure. Responses are untyped dicts; field definitions, filter constants, and enum values are documented in the [Go reference](https://pkg.go.dev/github.com/ByteSizedMarius/rewerse-engineering/pkg#pkg-types).

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
| `get_product_recommendations(market_id, listing_id)` | Related product recommendations |

See [`ProductFilter`](https://pkg.go.dev/github.com/ByteSizedMarius/rewerse-engineering/pkg#ProductFilter). Prices in `currentRetailPrice` are in **cents** (divide by 100 for euros).

### Discounts

| Method | Description |
|--------|-------------|
| `get_discounts(market_id)` | Get parsed weekly discounts |
| `get_discounts_raw(market_id)` | Get raw discount data including handout links |

### Shop & Services

| Method | Description |
|--------|-------------|
| `get_shop_overview(market_id, *, service_type="PICKUP", zip_code="")` | Get product categories for a market's online shop |
| `get_service_portfolio(zipcode)` | Get available REWE services for a zip code |
| `get_recalls()` | Get current product recalls |

See [`ServiceType`](https://pkg.go.dev/github.com/ByteSizedMarius/rewerse-engineering/pkg#ServiceType). Must match market capabilities (check `serviceFlags.hasPickup` from `get_market_details()`).

### Basket

| Method | Description |
|--------|-------------|
| `create_basket(market_id, zip_code, service_type="PICKUP")` | Create a new shopping basket |
| `get_basket(basket_id, market_id, zip_code, service_type, version=0)` | Get current basket state |
| `set_basket_item(basket_id, market_id, zip_code, service_type, listing_id, quantity, version)` | Add/update item quantity in basket |
| `remove_basket_item(basket_id, market_id, zip_code, service_type, listing_id)` | Remove item from basket |

See [`ServiceType`](https://pkg.go.dev/github.com/ByteSizedMarius/rewerse-engineering/pkg#ServiceType). `version`: basket state version for optimistic locking (from previous basket response). Prices in cents.

### Delivery

| Method | Description |
|--------|-------------|
| `get_bulky_goods_config(market_id, service_type="DELIVERY")` | Get beverage crate limits and surcharges |

## macOS

The PyPI wheel doesn't include a macOS binary. Build from source on a Mac:

```bash
git clone https://github.com/ByteSizedMarius/rewerse-engineering.git
cd rewerse-engineering
python python/build_lib.py --platform darwin
cd python && pip install -e .
```

Requires Go 1.21+ and Xcode command line tools (`xcode-select --install`).

## License

MIT
