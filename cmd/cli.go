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
		fmt.Println("\nSubcommand-Flags:")
		fmt.Println(align("-query <query>") + "Search query. Can be a city, postal code, street, etc.")
		fmt.Println(align("-id <market-id>") + "ID of the market or discount. Get it from marketsearch.")
		fmt.Println(align("-raw") + "Raw output format (directly from the API). Otherwise, it's parsed and cleaned based on what I thought was useful.")
		fmt.Println(align("-catGroup") + "Group by product category instead of rewe-app-category")
		fmt.Println("\nExamples:")
		fmt.Println("   rewerse.exe marketsearch -query Köln")
		fmt.Println("   rewerse.exe marketdetails -id 1763153")
		fmt.Println("   rewerse.exe discounts -id 1763153")
		fmt.Println("   rewerse.exe -json discounts -id 1763153 -raw")
		fmt.Println("   rewerse.exe discounts -id 1763153 -catGroup")
		fmt.Println("   rewerse.exe -cert cert.pem -key p.key marketsearch -query Köln")
	}
	certFile := flag.String("cert", "certificate.pem", "Path to the certificate file")
	keyFile := flag.String("key", "private.key", "Path to the key file")
	jsonOutput := flag.Bool("json", false, "Output in JSON format (default false)")
	flag.Parse()

	var err error
	err = rewerse.SetCertificate(*certFile, *keyFile)
	hdl(err)

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	marketSearchCmd := flag.NewFlagSet("marketsearch", flag.ExitOnError)
	marketSearchQuery := marketSearchCmd.String("query", "", "Search query")

	marketDetailsCmd := flag.NewFlagSet("marketdetails", flag.ExitOnError)
	marketDetailsID := marketDetailsCmd.String("id", "", "ID of the market details")

	discountsCmd := flag.NewFlagSet("discounts", flag.ExitOnError)
	discountsID := discountsCmd.String("id", "", "Market-ID")
	rawOutput := discountsCmd.Bool("raw", false, "Whether to return raw API Data")
	groupProduct := discountsCmd.Bool("catGroup", false, "Group by product category instead of app-category")

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

	default:
		fmt.Println("Expected 'marketsearch', 'marketdetails', 'coupons', 'recalls' or 'discounts' subcommands")
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
	return "   " + s + strings.Repeat(" ", 40-len(s))
}
