package main

import (
	"flag"
	"fmt"
	"os"

	rewerse "github.com/ByteSizedMarius/rewerse-engineering/pkg"
)

func handleDiscounts(args []string) (any, error) {
	if wantsHelp(args) {
		discountsHelp()
		return nil, nil
	}

	fs := flag.NewFlagSet("discounts", flag.ContinueOnError)
	market := fs.String("market", "", "Market ID")
	raw := fs.Bool("raw", false, "Raw API response")
	group := fs.Bool("group", false, "Group by product category")
	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if err := checkUnexpectedArgs(fs); err != nil {
		return nil, err
	}

	if err := validateNumeric("market", *market); err != nil {
		return nil, err
	}

	if *raw {
		return rewerse.GetDiscountsRaw(*market)
	}

	data, err := rewerse.GetDiscounts(*market)
	if err != nil {
		return nil, err
	}
	if *group {
		return data.GroupByProductCategory(), nil
	}
	return data, nil
}

func handleCategories(args []string, jsonOutput bool) (any, error) {
	if wantsHelp(args) {
		categoriesHelp()
		return nil, nil
	}

	fs := flag.NewFlagSet("categories", flag.ContinueOnError)
	market := fs.String("market", "", "Market ID")
	service := fs.String("service", "", "Service type: PICKUP or DELIVERY")
	all := fs.Bool("all", false, "Print all categories")
	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if err := checkUnexpectedArgs(fs); err != nil {
		return nil, err
	}

	if err := validateNumeric("market", *market); err != nil {
		return nil, err
	}

	var opts *rewerse.ShopOverviewOpts
	if *service != "" {
		opts = &rewerse.ShopOverviewOpts{ServiceType: rewerse.ServiceType(*service)}
	}
	so, err := rewerse.GetShopOverviewWithOpts(*market, opts)
	if err != nil {
		return nil, err
	}

	if *all && jsonOutput {
		fmt.Fprintln(os.Stderr, "Cannot use -all with -json (json includes all by default)")
		return nil, nil
	}

	if *all {
		fmt.Println(so.StringAll())
		return nil, nil
	} else if !jsonOutput {
		fmt.Println(so.String())
		return nil, nil
	}
	return so.ProductCategories, nil
}

func handleServices(args []string) (any, error) {
	if wantsHelp(args) {
		servicesHelp()
		return nil, nil
	}

	fs := flag.NewFlagSet("services", flag.ContinueOnError)
	zip := fs.String("zip", "", "Zip code")
	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if err := checkUnexpectedArgs(fs); err != nil {
		return nil, err
	}

	if err := validateZipCode(*zip); err != nil {
		return nil, err
	}
	return rewerse.GetServicePortfolio(*zip)
}

func discountsHelp() {
	fmt.Printf(`Usage: %s discounts [flags]

Flags:
  -market     Market ID (required)
  -raw        Return raw API response
  -group      Group results by product category

Examples:
  %s discounts -market 840174
  %s discounts -market 840174 -raw
  %s discounts -market 840174 -group
`, binaryName, binaryName, binaryName, binaryName)
}

func categoriesHelp() {
	fmt.Printf(`Usage: %s categories [flags]

Flags:
  -market     Market ID (required)
  -service    PICKUP or DELIVERY (default: PICKUP, must match market capabilities)
  -all        Print all categories with subcategories

Note: Use 'markets details -id <market>' to check hasPickup field for service support.

Examples:
  %s categories -market 831002
  %s categories -market 320509 -service DELIVERY
  %s categories -market 831002 -all
`, binaryName, binaryName, binaryName, binaryName)
}

func servicesHelp() {
	fmt.Printf(`Usage: %s services [flags]

Flags:
  -zip        Zip code (required, 5 digits)

Examples:
  %s services -zip 50667
`, binaryName, binaryName)
}
