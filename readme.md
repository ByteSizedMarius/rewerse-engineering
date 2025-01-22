# rewerse engineering

<div>
  <img src="gopher.png" alt="Project Logo" width="170" align="right">
  <p>This repository implements many API endpoints used by the Rewe app for querying current discounts, products, coupons and recalls.</p>
  <p>Current supported APK version: 4.0.3 (as of 22.01.25)</p> 
</div>

> [!CAUTION]
> The certificates required for talking to the rewe api are not included in this repository. You need to extract them from the APK. Documentation & an extraction-script for windows can be found in the [docs](./docs) directory.

## intro

In [March 2024](https://github.com/foo-git/rewe-discounts/issues/19), Rewe started
using [Cloudflare MTLS](https://www.cloudflare.com/learning/access-management/what-is-mutual-tls/) to secure their api-endpoints, which
broke [existing solutions](https://github.com/foo-git/rewe-discounts) that allowed, for example, fetching discounts for a specific Rewe market.
Github-user [@torbenpfohl](https://github.com/torbenpfohl) was ~~obsessed~~ persistent enough to figure this out, find the certificate and it's password. This repo is based on
his [work](https://github.com/torbenpfohl/rewe-discounts/blob/main/how%20to%20get%20private.pem%20and%20private.key.txt) and aims to document the required procedures and implement some of the endpoints.

This repo is not meant to cause any harm. It's sole purpose is to give access to data that is already freely
accessible via the app/website. It will never contain endpoints related to rewe account information (login/shoppinglist/etc.) or payback.

Not affiliated with Rewe in any way.

## contents

A basic go implementation + documentation of the rewe api is available in the [pkg](pkg) directory. The CLI-implementation (quick & dirty currently) is in [cmd](cmd). Releases are in [releases](https://github.com/ByteSizedMarius/rewerse-engineering/releases). 

Please note that since this is an unsigned go binary that does some encryption/decryption of certificates and sends webrequests to the rewe api, it will likely get flagged by your antivirus. The only dependency is [google/uuid](https://github.com/google/uuid), so you can easily compile it yourself.

At some point this repo will contain:
- a tui for the go-implementation
- a powershell script for extracting the password incase it changes (script will not work if anything else changes because of static linking)
- (maybe) a basic python implementation for fetching discounts if no one wants to maintain one ([torbens fork](https://github.com/torbenpfohl/rewe-discounts) is already functional)

## cli

```
Usage: rewerse-cli [flags] [subcommand] [subcommand-flags]

Flags:
   -cert <path>                                      Path to the certificate file (default 'certificate.pem')
   -key <path>                                       Path to the key file (default 'private.key')
   -json                                             Output in JSON format (default false)

Subcommands:
   marketsearch -query <query>                       Search for markets
   marketdetails -id <market-id>                     Get details for a market
   coupons                                           Get coupons
   recalls                                           Get recalls
   discounts -id <market-id> [-raw] [-catGroup]      Get discounts
   categories -id <market-id> [-printAll]            Get product categories
   products -id <market-id> [-category <category> -search <query>] [-page <page>] [-perPage <productsPerPage>]
                                                     Get products from a category or by query

Subcommand-Flags:
   -query <query>                                    Search query. Can be a city, postal code, street, etc.
   -id <market-id>                                   ID of the market or discount. Get it from marketsearch.
   -raw                                              Whether you want raw output format (directly from the API) or parsed.
   -catGroup                                         Group by product category instead of rewe-app category
   -printAll                                         Print all available product categories (very many)
   -category <category>                              The slug of the category to fetch the products from
   -search <query>                                   Search query for products. Can be any term or EANs for example
   -page <page>                                      Page number for pagination. Starts at 1, default 1. The amount of available pages is included in the output
   -perPage <productsPerPage>                        Number of products per page. Default 30

Examples:
   rewerse.exe -cert cert.pem -key p.key marketsearch -query Köln
   rewerse.exe marketsearch -query Köln
   rewerse.exe marketdetails -id 1763153
   rewerse.exe discounts -id 1763153
   rewerse.exe -json discounts -id 1763153 -raw
   rewerse.exe discounts -id 1763153 -catGroup
   rewerse.exe categories -id 831002 -printAll
   rewerse.exe products -id 831002 -category kueche-haushalt -page 2 -perPage 10
   rewerse.exe products -id 831002 -search Karotten
```

## misc

Feel free to open github issues for suggestions, questions, bugs. PRs welcome. Email: rewe at marius dot codes.

## attribution

- https://github.com/foo-git/rewe-discounts
- https://github.com/torbenpfohl/rewe-discounts
- https://github.com/charmbracelet/bubbletea
- https://github.com/egonelbre/gophers