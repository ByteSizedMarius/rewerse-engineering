# rewerse engineering

<div>
  <img src="gopher.png" alt="Project Logo" width="170" align="right">
  <p>This repository aims to implement all publicly accessible (unauthenticated) API endpoints used by the Rewe app for querying current discounts, products, recipes and recalls.</p>
  <p>Current supported APK version: 5.7.3 (as of 09.01.26)</p> 
</div>

> [!CAUTION]
> The certificates required for talking to the rewe api are not included in this repository. You need to extract them from the APK. Documentation & an extraction-script for windows can be found in the [docs](./docs) directory.

## quick start

1. **Extract certificates** from the rewe apk; see [docs](./docs) for instructions

2. **Download** the [latest release](https://github.com/ByteSizedMarius/rewerse-engineering/releases/latest) or clone and build:
   ```
   go build -o rewerse ./cmd
   ```

3. **Run** (certificates must be in working directory or specified via flags):
   ```
   ./rewerse.exe -json discounts -market 840174
   ```

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

A basic go implementation + documentation of the rewe api is available in the [pkg](pkg) directory. See the [readme](pkg/readme.md) for API documentation with usage examples. The CLI-implementation is in [cmd](cmd). Releases are in [releases](https://github.com/ByteSizedMarius/rewerse-engineering/releases). 

Please note that since this is an unsigned go binary that does some encryption/decryption of certificates and sends webrequests to the rewe api, it will likely get flagged by your antivirus. There are no dependencies, so you can easily compile it yourself.

## cli

```
Usage: ./rewerse.exe [flags] <command> [subcommand] [flags]

Flags:
  -cert <path>    Certificate file (default: certificate.pem)
  -key <path>     Key file (default: private.key)
  -json           Output as JSON

Commands:
  markets         Search and get market details
  products        Search, browse, and get product info
  recipes         Search and browse recipes
  discounts       Get market discounts
  categories      Get product categories
  recalls         Get product recalls
  services        Get service portfolio by zip

Examples:
  ./rewerse.exe markets search -query Köln
  ./rewerse.exe products search -market 831002 -query Milch
  ./rewerse.exe products category -market 831002 -slug obst-gemuese
  ./rewerse.exe recipes search -term Pasta
  ./rewerse.exe discounts -market 840174
  ./rewerse.exe categories -market 831002
  ./rewerse.exe services -zip 50667

Run './rewerse.exe <command>' for subcommand help.
```

## misc

Feel free to open github issues for suggestions, questions, bugs. PRs welcome. Email: rewe at marius dot codes.

## attribution

- https://github.com/foo-git/rewe-discounts
- https://github.com/torbenpfohl/rewe-discounts
- https://github.com/egonelbre/gophers
