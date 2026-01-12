# pkg

[![Go Reference](https://pkg.go.dev/badge/github.com/ByteSizedMarius/rewerse-engineering.svg)](https://pkg.go.dev/github.com/ByteSizedMarius/rewerse-engineering)

- [Usage](#usage)
- [Service Types](#service-types)
- [Output Examples](#output-examples)

This Go package implements publicly accessible (unauthenticated) API endpoints used by the Rewe app for querying current discounts, products, recipes and recalls. The following functions are currently available:

- `MarketSearch`: Search for Rewe markets using city, street, PLZ, market name, etc. Returns a list of markets with some basic information.
- `GetMarketDetails`: Get details about a specific market using the unique market-id.
- `GetDiscountsRaw`: Get all discounts for a specific market using the unique market-id. Returns the raw response parsed into a struct.
- `GetDiscounts`: Get all discounts for a specific market using the unique market-id. Returns data cleaned of information deemed unnecessary by me (opinionated).
- `GetRecalls`: Get product recalls.
- `GetRecipeHub`: Returns data from the recipe-page in the Rewe app.
- `GetShopOverview`: Returns the Product-Categories available in the given market. Only available for markets who are pickup- or delivery-enabled.
- `GetCategoryProducts`: Returns all products in a given category. Only available for markets who are pickup- or delivery-enabled.
- `GetProducts`: Returns all products for a given query. Only available for markets who are pickup- or delivery-enabled.
- `GetServicePortfolio`: Returns delivery/pickup availability for a zip code.

## Usage

To use any of the functions, the certificate of the Rewe app must be provided:
```go
err := rewerse.SetCertificate("certificate.pem", "private.key")
if err != nil {
    panic(err)
}
```

This also initializes a simulated device session:
- A random user agent is selected from 16 phone models (Samsung, Pixel, OnePlus, Xiaomi, LG)
- UUIDs are generated for `x-instana-android` (APM telemetry) and `rdfa` (device fingerprint) headers

These values are set once and reused for all subsequent requests. The config is stored using an atomic pointer, so concurrent API calls are safe after initialization.

## Service Types

Product and category queries (`GetProducts`, `GetCategoryProducts`, `GetShopOverview`) require a service type that must match market capabilities. Markets support either PICKUP, DELIVERY, or both.

Use `GetMarketDetails` to check the `hasPickup` field in `serviceFlags`. If a market doesn't support PICKUP, specify DELIVERY via `ProductOpts.ServiceType` or `ShopOverviewOpts.ServiceType`. Using the wrong service type returns HTTP 400 "Invalid market selection".

```go
// For markets that support PICKUP (default)
pr, _ := rewerse.GetProducts("831002", "milch", nil)

// For markets that only support DELIVERY
opts := &rewerse.ProductOpts{ServiceType: rewerse.ServiceDelivery}
pr, _ := rewerse.GetProducts("840174", "milch", opts)
```

## Output Examples

I have written very basic string-functions for all structs except RawDiscounts, which allows printing the data to console easily. Below are examples of the data:

**MarketSearch:**
```go
markets, _ := rewerse.MarketSearch("Mannheim")
fmt.Println(markets)
```

```
ID      Location
862012: REWE Markt, N 1, 68161 Mannheim
830843: REWE Markt, Q6/Q7, 68161 Mannheim
840099: REWE Markt, S 6 29 - 32, 68161 Mannheim
840906: REWE Markt, Lange Rötter Str. 11 - 17, 68167 Mannheim
831297: REWE Markt, Ulmenweg 9, 68167 Mannheim
840832: REWE Markt, Steubenstr. 76, 68199 Mannheim
565236: REWE Center, Amselstr. 10, 68307 Mannheim
840174: REWE Markt, Rheingoldstr. 18-20, 68199 Mannheim / Neckarau
[...]
```

**GetMarketDetails:**
```go
market, _ := rewerse.GetMarketDetails("840174")
fmt.Println(market)
```

```
---Allgemein---------------------------------------------------
   ID:                    840174
   Name:                  REWE Markt (REWE Markt Vuthaj oHG)
   Type:                  REWE Markt
   Standort:              Rheingoldstr. 18-20, 68199 Mannheim / Neckarau (49.4539, 8.4904)
   Tel-Nr:                0621-8414498

---Öffnungszeiten---------------------------------------------
   Aktuell:               Schließt bald (um 22:00 Uhr)
   Wochenübersicht:       [{"days":"Mo - Sa","hours":"07:00 - 22:00"}]

---Services----------------------------------------------------
   Services:              {"fixed":[{"text":"Fleisch- und Wursttheke",...}],"editable":[...]}
   Pickup:                false
```

**GetDiscountsRaw:**
```go
discounts, _ := rewerse.GetDiscountsRaw("840174")
fmt.Println(discounts)
```

Outputs data in the format of the struct [`RawDiscounts`](./discounts_struct.go).

**GetDiscounts:**
```go
discounts, _ := rewerse.GetDiscounts("840174")
fmt.Println(discounts)
```

```
Top-Angebote in deinem Markt
    Haribo Goldbären oder Color-Rado, 0.77€
    Coca-Cola, Fanta, Sprite oder Mezzo Mix, 11.99€
    YFood Trinkmahlzeit, 3.29€
    Ritter Sport Schokolade, 1.11€
    [...]
Bonus-Aktionen
    YFood Trinkmahlzeit, 3.29€
    Rügenwalder Vegane Mühlen Cordon bleu, 2.49€
    Kinder Country, 2.22€
    Philadelphia, 1.11€
    [...]

[...]
```

Discounts can be reordered by product category:
```go
discounts, _ := rewerse.GetDiscounts("840174")
discounts = discounts.GroupByProductCategory()
fmt.Println(discounts)
```

```
wein-und-spirituosen
    Grand Sud Vin de France, 3.49€
    MM Extra Sekt, 2.99€
    Henkell Sekt, 3.99€
    [...]

suesses-und-salziges
    Haribo Goldbären oder Color-Rado, 0.77€
    Ritter Sport Schokolade, 1.11€
    Kinder Country, 2.22€
    [...]

obst-und-gemuese
    Bio Speisemöhren, 0.99€
    Bio Zitronen, 1.11€
    Dunkle Tafeltrauben, 1.99€
    [...]

[...]
```

**GetRecalls:**
```go
recalls, _ := rewerse.GetRecalls()
fmt.Println(recalls)
```

```
Recalls:
Vorsorglicher Produktrückruf von verschiedenen Beba Produkten
Mögliches Vorhandensein von Cereulid
https://mediacenter.rewe.de/produktrueckrufe/beba-produkte
```

**GetRecipeHub:**
```go
recipes, _ := rewerse.GetRecipeHub()
fmt.Println(recipes)
```

```
Recipe Hub

Recipe of the Day
--------------------
Goi Xoai Mangosalat
30 min
Mittel
https://www.rewe.de/rezepte/goi-xoai-mangosalat/

Popular Recipes
--------------------
Goi Xoai Mangosalat
30 min
Mittel
https://www.rewe.de/rezepte/goi-xoai-mangosalat/

Hack-Reis-Pfanne
30 min
Einfach
https://www.rewe.de/rezepte/hack-reis-pfanne/

[...]

Available Categories
--------------------
Vegetarisch
Fleisch
Backen
Nachspeisen
Fisch
Vorspeisen
Kuchen
```

**GetServicePortfolio:**

```go
sp, _ := rewerse.GetServicePortfolio("50667")
fmt.Println(sp)
```

```
Service Portfolio for 50667:
  Delivery available from market 320509
  20 pickup markets available:
    - 3200008 (REWE Markt): Neumarkt 35-37, 50667 Köln
    - 1765506 (REWE Markt): Salierring 47-53, 50677 Köln
    - 1940295 (REWE Markt): Vogelsangerstr. 16, 50823 Köln / Ehrenfeld
    [...]
```

**GetShopOverview:**

```go
so, _ := rewerse.GetShopOverview("831002")
fmt.Println(so.StringAll())
```

```
---Monats-Highlights-------------------------------------------
   ID:                    3844
   Name:                  Monats-Highlights
   Slug:                  monatshighlights
   ProductCount:          538
   ImageURL:              https://shop.rewe-static.de/mobile/categories/images/v2/3844.png
   Child-Amnt:            6
   Children:
      Neu im Sortiment
         ID:                    3845
         Name:                  Neu im Sortiment
         Slug:                  neu-im-sortiment
         ProductCount:          97
         ImageURL:              https://shop.rewe-static.de/...
         Child-Amnt:            0
      Fit ins neue Jahr
         ID:                    3908
         Name:                  Fit ins neue Jahr
         Slug:                  fit-ins-neue-jahr
         ProductCount:          228
         [...]
```

**GetCategoryProducts:**

```go
cs, _ := rewerse.GetCategoryProducts("831002", "obst-gemuese", nil)
fmt.Println(cs)
```

```
Products for query *
Page 1 of 8
---------------------------------------------------------------
Tafeltrauben hell kernlos 500g
Cherry Romatomaten 250g
Zespri Kiwi Green 1 Stück
Salatgurke 1 Stück
REWE Beste Wahl Delizioso Mini Cherry Rispentomaten 200g
Chicoree 500g
Zwiebel rot 500g im Netz
Zucchini grün ca. 300g
[...]
```

**GetProducts:**

```go
pr, _ := rewerse.GetProducts("831002", "milch", nil)
fmt.Println(pr)
```

```
Products for query milch
Page 1 of 24
---------------------------------------------------------------
Weihenstephan H-Milch 3,5% 1l
vly Erbsen-Drink High Protein vegan 1l
yfood Trinkmahlzeit Smooth Vanilla 500ml
REWE Frei von Fettarme H-Milch laktosefrei 1,5% 1l
Arla LactoFree Laktosefreie Milch 1,5% 1l
MinusL H-Milch laktosefrei 1,5% 1l
Weihenstephan Haltbare Milch 1,5% 1l
[...]
```
