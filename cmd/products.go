package main

import (
	"flag"
	"fmt"

	rewerse "github.com/ByteSizedMarius/rewerse-engineering/pkg"
)

func handleProducts(args []string) (any, error) {
	if wantsHelp(args) {
		productsHelp()
		return nil, nil
	}

	switch args[0] {
	case "search":
		fs := flag.NewFlagSet("products search", flag.ContinueOnError)
		market := fs.String("market", "", "Market ID")
		query := fs.String("query", "", "Search query")
		service := fs.String("service", "", "Service type: PICKUP or DELIVERY")
		page := fs.Int("page", 0, "Page number")
		perPage := fs.Int("perPage", 0, "Results per page")
		if err := fs.Parse(args[1:]); err != nil {
			return nil, err
		}
		if err := checkUnexpectedArgs(fs); err != nil {
			return nil, err
		}
		if err := validateNumeric("market", *market); err != nil {
			return nil, err
		}
		if err := validateFlag("query", *query); err != nil {
			return nil, err
		}
		return rewerse.GetProducts(*market, *query, buildOpts(*page, *perPage, *service))

	case "category":
		fs := flag.NewFlagSet("products category", flag.ContinueOnError)
		market := fs.String("market", "", "Market ID")
		slug := fs.String("slug", "", "Category slug")
		service := fs.String("service", "", "Service type: PICKUP or DELIVERY")
		page := fs.Int("page", 0, "Page number")
		perPage := fs.Int("perPage", 0, "Results per page")
		if err := fs.Parse(args[1:]); err != nil {
			return nil, err
		}
		if err := checkUnexpectedArgs(fs); err != nil {
			return nil, err
		}
		if err := validateNumeric("market", *market); err != nil {
			return nil, err
		}
		if err := validateFlag("slug", *slug); err != nil {
			return nil, err
		}
		return rewerse.GetCategoryProducts(*market, *slug, buildOpts(*page, *perPage, *service))

	case "details":
		fs := flag.NewFlagSet("products details", flag.ContinueOnError)
		market := fs.String("market", "", "Market ID")
		product := fs.String("product", "", "Product ID")
		if err := fs.Parse(args[1:]); err != nil {
			return nil, err
		}
		if err := checkUnexpectedArgs(fs); err != nil {
			return nil, err
		}
		if err := validateNumeric("market", *market); err != nil {
			return nil, err
		}
		if err := validateFlag("product", *product); err != nil {
			return nil, err
		}
		return rewerse.GetProductByID(*market, *product)

	case "suggest":
		fs := flag.NewFlagSet("products suggest", flag.ContinueOnError)
		query := fs.String("query", "", "Search query")
		page := fs.Int("page", 0, "Page number")
		perPage := fs.Int("perPage", 0, "Results per page")
		if err := fs.Parse(args[1:]); err != nil {
			return nil, err
		}
		if err := checkUnexpectedArgs(fs); err != nil {
			return nil, err
		}
		if err := validateFlag("query", *query); err != nil {
			return nil, err
		}
		return rewerse.GetProductSuggestions(*query, buildOpts(*page, *perPage, ""))

	case "recommend":
		fs := flag.NewFlagSet("products recommend", flag.ContinueOnError)
		market := fs.String("market", "", "Market ID")
		listing := fs.String("listing", "", "Listing ID")
		if err := fs.Parse(args[1:]); err != nil {
			return nil, err
		}
		if err := checkUnexpectedArgs(fs); err != nil {
			return nil, err
		}
		if err := validateNumeric("market", *market); err != nil {
			return nil, err
		}
		if err := validateFlag("listing", *listing); err != nil {
			return nil, err
		}
		return rewerse.GetProductRecommendations(*market, *listing)

	default:
		productsHelp()
		return nil, fmt.Errorf("unknown products subcommand: %s", args[0])
	}
}

func productsHelp() {
	fmt.Printf(`Usage: %s products <subcommand> [flags]

Subcommands:
  search      Search for products
  category    Get products from a category
  details     Get product details
  suggest     Get search suggestions
  recommend   Get product recommendations

products search:
  -market     Market ID (required)
  -query      Search query (required)
  -service    PICKUP or DELIVERY (default: PICKUP, must match market capabilities)
  -page       Page number
  -perPage    Results per page

products category:
  -market     Market ID (required)
  -slug       Category slug (required, use 'categories' command to list)
  -service    PICKUP or DELIVERY (default: PICKUP, must match market capabilities)
  -page       Page number
  -perPage    Results per page

products details:
  -market     Market ID (required)
  -product    Product ID (required)

products suggest:
  -query      Search query (required)
  -page       Page number
  -perPage    Results per page

products recommend:
  -market     Market ID (required)
  -listing    Listing ID (required)

Note: Use 'markets details -id <market>' to check hasPickup field for service support.

Examples:
  %s products search -market 831002 -query Karotten
  %s products search -market 840174 -query Tomaten -service DELIVERY
  %s products category -market 831002 -slug obst-gemuese
  %s products details -market 831002 -product 9900011
  %s products suggest -query Milch
`, binaryName, binaryName, binaryName, binaryName, binaryName, binaryName)
}
