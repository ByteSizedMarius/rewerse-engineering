# rewerse engineering

> [!CAUTION]
> The certificates required for this repository are not included. You need to extract them from the APK. Documentation & an extraction-script for windows can be found in the [docs](./docs) directory.

Current supported APK version: 3.18.5 (as of 27.05.24)

## intro

In [March 2024](https://github.com/foo-git/rewe-discounts/issues/19), Rewe started
using [Cloudflare MTLS](https://www.cloudflare.com/learning/access-management/what-is-mutual-tls/) to secure their api-endpoints, which
broke [existing solutions](https://github.com/foo-git/rewe-discounts) that allowed, for example, fetching discounts for a specific Rewe market.
Github-user [@torbenpfohl](https://github.com/torbenpfohl) was ~~obsessed~~ persistent enough to figure this out, find the certificate and it's password. This repo is based on
his [work](https://github.com/torbenpfohl/rewe-discounts/blob/requests_based/how%20to%20get%20private.pem%20and%20private.key.txt) and aims to document the required procedures and implement some of the endpoints.

This repo is not meant to cause any harm. It's sole purpose is to give access to data that is already freely
accessible via the app/website. It will never contain endpoints related to rewe account information (login/shoppinglist/etc.) or payback.

Not affiliated with Rewe in any way.

## contents

A basic go implementation + documentation of the rewe api is available in the [pkg](pkg) directory. CLI is in [cmd](cmd).

At some point this repo will contain:
- a tui for the go-implementation
- a powershell script for extracting the password incase it changes (script will not work if anything else changes because of static linking)
- (maybe) a basic python implementation for fetching discounts if no one wants to maintain one ([torbens fork](https://github.com/torbenpfohl/rewe-discounts) is already functional)

## cli

```
> rewerse.exe
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

Subcommand-Flags:
   -query <query>                                    Search query. Can be a city, postal code, street, etc.
   -id <market-id>                                   ID of the market or discount. Get it from marketsearch.
   -raw                                              Raw output format (directly from the API). Otherwise, it's parsed and cleaned based on what I thought was useful.
   -catGroup                                         Group by product category instead of rewe-app-category

Examples:
   rewerse.exe marketsearch -query Köln
   rewerse.exe marketdetails -id 1763153
   rewerse.exe discounts -id 1763153
   rewerse.exe -json discounts -id 1763153 -raw
   rewerse.exe discounts -id 1763153 -catGroup
   rewerse.exe -cert cert.pem -key p.key marketsearch -query Köln
```

## misc

Feel free to open github issues for suggestions, questions, bugs. PRs welcome. Email: rewe at marius dot codes.

## attributions

- https://github.com/foo-git/rewe-discounts
- https://github.com/torbenpfohl/rewe-discounts
- https://github.com/charmbracelet/bubbletea
