# pkg

[![Go Reference](https://pkg.go.dev/badge/github.com/ByteSizedMarius/rewerse-engineering.svg)](https://pkg.go.dev/github.com/ByteSizedMarius/rewerse-engineering)

This Go Package implements some of the endpoints used by the Rewe app. The following functions are currently available:

- `MarketSearch`: Search for Rewe markets using city, street, PLZ, market name, etc. Returns a list of markets with some basic information.
- `GetMarketDetails`: Get details about a specific market using the unique market-id.
- `GetDiscountsRaw`: Get all discounts for a specific market using the unique market-id. Returns the raw response parsed into a struct.
- `GetDiscounts`: Get all discounts for a specific market using the unique market-id. Returns data cleaned of information deemed unnecessary by me (maybe opinionated).
- `GetCoupons`: Get all global coupons. Returns raw data parsed into structs. I have no idea, what exactly the data means, but it's there. Open an issue if you have details.
- `GetRecalls`: Get product recalls.
- `GetRecipeHub`: Returns data from the recipe-page in the Rewe app.
- `GetShopOverview`: Returns Recalls (again) and the Product-Categories available in the given market. Only available for markets who are pickup- or delivery-enabled.
- `GetCategoryProducts`: Returns all products in a given category. Only available for markets who are pickup- or delivery-enabled.
- `GetProducts`: Returns all products for a given query. Only available for markets who are pickup- or delivery-enabled.

## Usage

To use any of the functions, the certificate of the Rewe app must be provided:
```go
err := rewerse.SetCertificate("certificate.pem", "private.key")
if err != nil {
    panic(err)
}
```

I have written very basic string-functions for all structs except RawDiscounts, which allows printing the data to console easily. Below are examples of the data:

**MarketSearch:**
```go
markets, _ := rewerse.MarketSearch("Köln")
fmt.Println(markets)
```

```
1763153: REWE Andrea Flammuth oHG, Kölner Str. 26, 50859 Köln / Lövenich
1763107: REWE Richrath SM GmbH & Co. oHG, Rhöndorfer Str. 17, 50939 Köln / Klettenberg
8100134: REWE Kelterbaum oHG, Rudolfplatz 9, 50674 Köln
1940357: REWE Holger Bertram oHG, Aachener Str. 1253, 50858 Köln / Weiden
1940281: REWE Silke Huerten oHG, Hohe Str. 30, 50667 Köln / Altstadt
1765268: REWE Markt Björn Rohe oHG, Bergisch Gladbacher Str. 667, 51067 Köln / Holweide
1765091: REWE Peter Ziegler oHG, Neusser Str. 100, 50670 Köln / Neustadt-Nord
1478368: S-Mart Lebensmittel Märkte GmbH & Co. KG, Sülzgürtel 47, 50937 Köln / Sülz
[...]
```

**GetMarketDetails:**
```go
market, _ := rewerse.GetMarketDetails("1763153")
fmt.Println(market)
```

```
---Allgemein---------------------------------------------------
   ID:                    1763153
   Name:                  REWE Andrea Flammuth oHG
   Type-ID:               MARKET
   Standort:              Kölner Str. 26, 50859 Köln / Lövenich (50.9463, 6.8336)
   Tel-Nr:                02234-9489945
   Rohdaten:              {"attributes":[],"postalCode":"50859","city":"Köln / Lövenich"}

---Öffnungszeiten---------------------------------------------
   Aktell:                Geschlossen
   Nächste Änderung:      Öffnet am Sa, 07:00 (Typ: OPENS_NEXT_DAY)
   Wochenübersicht:       [{"days":"Mo - Sa","hours":"07:00 - 22:00"}]
   Besonderheiten:        []

---Verschiedenes-----------------------------------------------
   Features:              [] []
   Services:              []
   Aktionen:              [{ROUTE Route berechnen} {CALL Markt anrufen} {RATE Markt bewerten}]
   Bewertungs-URL:        https://meinfeedback.rewe.de/app/?p1=1763153&p2=2024-06-21T21:13:35Z&app=android
   Lieferservice:         false
```

**GetDiscountsRaw:**
```go
discounts, _ := rewerse.GetDiscountsRaw("1763153")
fmt.Println(discounts)
```

Outputs data in the format of the struct [`RawDiscounts`](./discounts_struct.go).

**GetDiscounts:**
```go
discounts, _ := rewerse.GetDiscounts("1763153")
fmt.Println(discounts)
```

```
Top-Angebote in deinem Markt
    Philadelphia, 0.95€
    Gustavo Gusto Pizza Margherita, 3.49€
    [...]
Anpfiff zum Sparen!
    Salakis Schafskäse Natur, 1.79€
    Lorenz Saltletts Pausen Cracker, 1.49€
    [...]
Obst und Gemüse
    Dunkle Tafeltrauben, 1.49€
    Aubergine, 1.00€
    [...]
PAYBACK Punkte Highlights
    Meggle Produkten, 0.00€
    Jacobs Kaffeekapseln, 0.00€

[...]
```

Discounts can be reordered by product category:
```go
discounts, _ := rewerse.GetDiscounts("1763153")
discounts = discounts.GroupByProductCategory()
fmt.Println(discounts)
```

```
Garten & Outdoor
    Lavendel, 1.99€
    Schleierkraut-Duo, 1.69€
    [...]
Tiefkühl
    Gustavo Gusto Pizza Margherita, 3.49€
    Mövenpick Bourbon Vanille, 1.99€
    [...]
Süßes & Salziges
    Ritter Sport Schokolade, 0.88€
    Chio Tortillas, 1.11€
    [...]
Küche & Haushalt
    Somat Excellence Geschirrreiniger, 7.77€
    Somat Excellence Geschirrreiniger, 7.77€
    [...]
 
[...]
```

**GetCoupons:**
```go
oneScanCoupons, coupons, _ := rewerse.GetCoupons()
fmt.Println(oneScanCoupons)
fmt.Println(coupons)
```

They are probably the same but in different formats. Not sure, I never used them. Let me know :)

```
Got 16 OneScan-Coupons

auf alle REWE Beste Wahl Produkte
Dieser Coupon kann pro Einkauf nur einmal eingelöst werden und ist nicht mit anderen Aktionen oder Gutscheinen kombinierbar.  *Dieser Coupon gilt nur in der REWE App und beim Kauf von mindestens zwei REWE Beste Wahl Produkten. Ausgenommen sind Aktionsartikel. Gültig vom 10.06. - 23.06.2024. Der Coupon ist nicht bei REWE nahkauf einlösbar. Eine Barauszahlung ist nicht möglich. Dieser Coupon ist auch einlösbar im REWE Lieferservice (ggf. abweichende Angebote) und REWE Abholservice. Bei REWE Lieferservice Bestellungen ist der Tag der Lieferung und bei REWE Abholservice Bestellungen der Tag der Abholung (nicht der Bestellung) maßgeblich.Gültig bis 23.06.2024
default_provider

Ritter Sport Schokolade versch. Sorten
Dieser Coupon kann pro Einkauf nur einmal eingelöst werden und ist nicht mit anderen Aktionen oder Gutscheinen kombinierbar.  *Dieser Coupon ist nicht bei REWE nahkauf einlösbar. Der Coupon ist nur beim Kauf von Ritter Sport Schokolade versch. Sorten je 100 g Tafel (1 kg = 7.70) gültig. Eine Barauszahlung ist nicht möglich. Keine Mehrfachrabattierung des gleichen Artikels durch gleichen oder weiteren Coupon möglich. Nur solange der Vorrat reicht. Änderungen im Sortiment vorbehalten. Produkte können ggf. nicht in allen Märkten erhältlich sein. Dieser Coupon ist auch einlösbar im REWE Lieferservice (ggf. abweichende Angebote) und REWE Abholservice. Bei REWE Lieferservice Bestellungen ist der Tag der Lieferung und bei REWE Abholservice Bestellungen der Tag der Abholung (nicht der Bestellung) maßgeblich.Gültig bis 23.06.2024
default_provider

[...]


Got 16 Coupons

auf alle REWE Beste Wahl Produkte // 5% Rabatt*
Beim Kauf von mindestens 2 Produkten

Ritter Sport Schokolade versch. Sorten // Knaller 0.77 € statt Aktionspreis 0.88 €!*
je 100 g Tafel (1 kg = 7.70)

[...]
```

**GetRecalls:**
```go
recalls, _ := rewerse.GetRecalls()
fmt.Println(recalls)
```

```
Recalls:
Vorsorglicher Produktrückruf von „ja! Basilikum gerebelt“
Mögliche mikrobielle Verunreinigung durch Salmonellen
https://mediacenter.rewe.de/produktrueckrufe/ja-basilikum-gerebelt

Vorsorglicher Produktrückruf von ELVIS Minze Zartbitter 400 ml
Das Produkt enthält eine fehlerhafte Allergenkennzeichnung auf der Vorderseite und auf der Rückseite.
https://mediacenter.rewe.de/produktrueckrufe/elvis-minze-zartbitter-400-ml
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
Kräuter-Fisch auf Bandnudeln
40 min
Einfach
https://www.rewe.de/rezepte/kraeuter-fisch-bandnudeln/

Popular Recipes
--------------------
Kräuter-Fisch auf Bandnudeln
40 min
Einfach
https://www.rewe.de/rezepte/kraeuter-fisch-bandnudeln/

Hähnchen Marinade Honig-Senf
15 min
Einfach
https://www.rewe.de/rezepte/haehnchen-marinade-honig-senf/

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

**GetShopOverview:**

```go
so, _ := rewerse.GetShopOverview("831002")
fmt.Println(so.StringAll())
```

```
---Grillsaison-------------------------------------------------
   ID:                    3752
   Name:                  Grillsaison
   Slug:                  grillsaison
   ProductCount:          475
   ImageURL:              https://shop.rewe-static.de/mobile/categories/images/v2/3752.png
   Child-Amnt:            8
   Children:
      Grillfleisch & -fisch
         ID:                    3753
         Name:                  Grillfleisch & -fisch
         Slug:                  grillfleisch-fisch
         ProductCount:          87
         ImageURL:              https://shop.rewe-static.de/mobile/categories/images/v2/3753.png
         Child-Amnt:            5
         Children:
            Bratwürstchen
               ID:                    3754
               Name:                  Bratwürstchen
               Slug:                  bratwuerstchen
               ProductCount:          14
               ImageURL:              https://shop.rewe-static.de/mobile/categories/images/v2/3754.png
               Child-Amnt:            0
               [...]
```

**GetCategoryProducts:**

```go
cs, _ := rewerse.GetCategoryProducts("831002", "bratwuerstchen", nil)
fmt.Println(cs)
```

```
Products for query *
Page 1 of 1
---------------------------------------------------------------
Wiesenhof Bruzzzler 400g
Wilhelm Brandenburg Rostbratwürstchen 200g
REWE Bio Original Nürnberger Rostbratwürstchen 160g
REWE Feine Welt Salsiccia Peperoncino 300g
Wolf Berner Würstchen 250g, 3 Stück
Wilhelm Brandenburg Metzgerbratwurst mittelgrob 400g
Steinhaus Krakauer mit Emmentaler geräuchert 500g
Die Thüringer Rostbratwurst 500g
[...]
```

**GetProducts:**

```go
pr, _ := rewerse.GetProducts("831002", "paprika", nil)
fmt.Println(pr)
```

```
Products for query paprika
Page 1 of 9
---------------------------------------------------------------
Paprika rot 500g
Spitzpaprika gelb 500g
Spitzpaprika rot 500g
Paprika Mix 500g
Paprika rot ca. 250g
Bio Paprika 400g
[...]
```

