# rewerse engineering

> [!CAUTION]
> The certificates required for this repository are not included. You need to extract them from the APK.

Current supported APK version: 3.18.5 (as of 27.05.24)

## intro

In [March 2024](https://github.com/foo-git/rewe-discounts/issues/19), rewe started
using [Cloudflare MTLS](https://www.cloudflare.com/learning/access-management/what-is-mutual-tls/) to secure their api-endpoints, which
broke [existing solutions](https://github.com/foo-git/rewe-discounts) that allowed, for example, fetching discounts for a specific rewe market.
Github-user @torbenpfohl was ~~obsessed~~ persistent enough to figure this out, find the certificate and it's password. This repo is based on
his [work](https://github.com/torbenpfohl/rewe-discounts/blob/requests_based/how%20to%20get%20private.pem%20and%20private.key.txt) and aims to
document the required procedures and implement some of the endpoints.

This repo is not meant to cause any harm. It's sole purpose is to give access to data that is already freely
accessible via the app/website. It will never contain endpoints related to rewe account information (login/shoppinglist/etc.) or payback.

Not affiliated with Rewe in any way.

## contents

At some point this repo will contain:

- instructions/scripts for extracting the certificate
- an implementation of basic api endpoints in Go
- a cli for the Go-implementation
- (maybe) a basic python implementation for fetching discounts ([torbens fork](https://github.com/torbenpfohl/rewe-discounts) is already functional)

## misc

Feel free to open github issues for suggestions, questions, bugs. PRs welcome. Email: rewe at marius dot codes.

## attributions

- https://github.com/foo-git/rewe-discounts
- https://github.com/torbenpfohl/rewe-discounts
- https://github.com/charmbracelet/bubbletea
