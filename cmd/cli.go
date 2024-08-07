package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ByteSizedMarius/rewerse-engineering/pkg"
	"os"
	"strings"
)

func hdl(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	flag.Usage = func() {
		fmt.Println("Usage: rewerse-cli [flags] [subcommand] [subcommand-flags]")
		fmt.Println("\nFlags:")
		fmt.Println(align("-cert <path>") + "Path to the certificate file (default 'certificate.pem')")
		fmt.Println(align("-key <path>") + "Path to the key file (default 'private.key')")
		fmt.Println(align("-json") + "Output in JSON format (default false)")
		fmt.Println("\nSubcommands:")
		fmt.Println(align("marketsearch -query <query>") + "Search for markets")
		fmt.Println(align("marketdetails -id <market-id>") + "Get details for a market")
		fmt.Println(align("coupons") + "Get coupons")
		fmt.Println(align("recalls") + "Get recalls")
		fmt.Println(align("discounts -id <market-id> [-raw] [-catGroup]") + "Get discounts")
		fmt.Println(align("categories -id <market-id> [-printAll]") + "Get product categories")
		fmt.Println(align("products -id <market-id> [-category <category> -search <query>] [-page <page>] [-perPage <productsPerPage>]") + "\n" + align("") + "Get products from a category or by query")
		fmt.Println("\nSubcommand-Flags:")
		fmt.Println(align("-query <query>") + "Search query. Can be a city, postal code, street, etc.")
		fmt.Println(align("-id <market-id>") + "ID of the market or discount. Get it from marketsearch.")
		fmt.Println(align("-raw") + "Whether you want raw output format (directly from the API) or parsed.")
		fmt.Println(align("-catGroup") + "Group by product category instead of rewe-app category")
		fmt.Println(align("-printAll") + "Print all available product categories (very many)")
		fmt.Println(align("-category <category>") + "The slug of the category to fetch the products from")
		fmt.Println(align("-search <query>") + "Search query for products. Can be any term or EANs for example")
		fmt.Println(align("-page <page>") + "Page number for pagination. Starts at 1, default 1. The amount of available pages is included in the output")
		fmt.Println(align("-perPage <productsPerPage>") + "Number of products per page. Default 30")
		fmt.Println("\nExamples:")
		fmt.Println("   rewerse.exe -cert cert.pem -key p.key marketsearch -query Köln")
		fmt.Println("   rewerse.exe marketsearch -query Köln")
		fmt.Println("   rewerse.exe marketdetails -id 1763153")
		fmt.Println("   rewerse.exe discounts -id 1763153")
		fmt.Println("   rewerse.exe -json discounts -id 1763153 -raw")
		fmt.Println("   rewerse.exe discounts -id 1763153 -catGroup")
		fmt.Println("   rewerse.exe categories -id 831002 -printAll")
		fmt.Println("   rewerse.exe products -id 831002 -category kueche-haushalt -page 2 -perPage 10")
		fmt.Println("   rewerse.exe products -id 831002 -search Karotten")
	}

	certFile := flag.String("cert", "", "Path to the certificate file")
	keyFile := flag.String("key", "", "Path to the key file")
	jsonOutput := flag.Bool("json", false, "Output in JSON format (default false)")
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	var crt, key string
	if *certFile == "" {
		crt = "certificate.pem"
	} else {
		crt = *certFile
	}
	if *keyFile == "" {
		key = "private.key"
	} else {
		key = *keyFile
	}

	var err error
	err = rewerse.SetCertificate(crt, key)
	if err != nil && *certFile == "" && *keyFile == "" {
		fmt.Println("Please provide the paths to the certificate and key files.\n\nrewerse.exe -cert cert.pem -key p.key [...]")
		os.Exit(1)
	}
	hdl(err)

	marketSearchCmd := flag.NewFlagSet("marketsearch", flag.ExitOnError)
	marketSearchQuery := marketSearchCmd.String("query", "", "Search query")

	marketDetailsCmd := flag.NewFlagSet("marketdetails", flag.ExitOnError)
	marketDetailsID := marketDetailsCmd.String("id", "", "ID of the market details")

	discountsCmd := flag.NewFlagSet("discounts", flag.ExitOnError)
	discountsID := discountsCmd.String("id", "", "Market-ID")
	rawOutput := discountsCmd.Bool("raw", false, "Whether to return raw API Data")
	groupProduct := discountsCmd.Bool("catGroup", false, "Group by product category instead of app-category")

	pcategories := flag.NewFlagSet("categories", flag.ExitOnError)
	pcategoriesMarketID := pcategories.String("id", "", "Market-ID")
	all := pcategories.Bool("printAll", false, "Print all available product categories")

	products := flag.NewFlagSet("products", flag.ExitOnError)
	productsMarketID := products.String("id", "", "Market-ID")
	productsCategory := products.String("category", "", "Category slug")
	productsSearch := products.String("search", "", "Search query")
	productsPage := products.Int("page", 0, "Page number")
	productsPerPage := products.Int("perPage", 0, "Products per page")

	var data any
	switch flag.Arg(0) {
	case "marketsearch":
		hdl(marketSearchCmd.Parse(flag.Args()[1:]))
		if *marketSearchQuery == "" {
			fmt.Println("Expected search query")
			os.Exit(1)
		}
		data, err = rewerse.MarketSearch(*marketSearchQuery)
	case "marketdetails":
		hdl(marketDetailsCmd.Parse(flag.Args()[1:]))
		if *marketDetailsID == "" {
			fmt.Println("Expected market ID")
			os.Exit(1)
		}
		data, err = rewerse.GetMarketDetails(*marketDetailsID)
	case "coupons":
		var oc rewerse.OneScanCoupons
		var c rewerse.Coupons
		oc, c, err = rewerse.GetCoupons()
		data = []any{oc, c}
	case "recalls":
		data, err = rewerse.GetRecalls()
	case "discounts":
		hdl(discountsCmd.Parse(flag.Args()[1:]))
		if *rawOutput {
			data, err = rewerse.GetDiscountsRaw(*discountsID)
			*jsonOutput = true
		} else {
			data, err = rewerse.GetDiscounts(*discountsID)
		}
		if *groupProduct {
			data = data.(rewerse.Discounts).GroupByProductCategory()
		}
	case "categories":
		hdl(pcategories.Parse(flag.Args()[1:]))
		if *pcategoriesMarketID == "" {
			fmt.Println("Expected market ID")
			os.Exit(1)
		}

		var so rewerse.ShopOverview
		so, err = rewerse.GetShopOverview(*pcategoriesMarketID)
		hdl(err)

		if *all && *jsonOutput {
			fmt.Println("Cannot print all and output in JSON format (json is all by default)")
			return
		}

		if *all {
			fmt.Println(so.StringAll())
			return
		} else if !*jsonOutput {
			fmt.Println(so.String())
			return
		} else {
			data = so.ProductCategories
		}
	case "products":
		hdl(products.Parse(flag.Args()[1:]))
		if *productsMarketID == "" {
			fmt.Println("Expected market ID")
			os.Exit(1)
		}

		if *productsCategory == "" && *productsSearch == "" {
			fmt.Println("Expected category or search query")
			os.Exit(1)
		}

		if *productsCategory != "" && *productsSearch != "" {
			fmt.Println("Expected either category or search query, not both")
			os.Exit(1)
		}

		var opts *rewerse.ProductOpts
		if *productsPage != 0 {
			opts = &rewerse.ProductOpts{
				Page: *productsPage,
			}
		}

		if *productsPerPage != 0 {
			if opts == nil {
				opts = &rewerse.ProductOpts{
					ObjectsPerPage: *productsPerPage,
				}
			} else {
				opts.ObjectsPerPage = *productsPerPage
			}
		}

		if *productsCategory != "" {
			data, err = rewerse.GetCategoryProducts(*productsMarketID, *productsCategory, opts)
		} else {
			data, err = rewerse.GetProducts(*productsMarketID, *productsSearch, opts)
		}

	default:
		fmt.Println("Expected one of the following subcommands:\n- marketsearch\n- marketdetails\n- coupons\n- recalls\n- discounts\n- categories\n- products")
		os.Exit(1)
	}

	hdl(err)
	if *jsonOutput {
		var bt []byte
		bt, err = json.MarshalIndent(data, "", "  ")
		hdl(err)
		data = string(bt)
	}
	fmt.Println(data)
}

func align(s string) string {
	al := 50
	if len(s) > al {
		al = 0
	} else {
		al -= len(s)
	}

	return "   " + s + strings.Repeat(" ", al)
}
