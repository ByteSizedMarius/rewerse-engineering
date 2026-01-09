// Exit codes:
//   0 - Success or help displayed
//   1 - Error (invalid args, API failure, etc.)
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	rewerse "github.com/ByteSizedMarius/rewerse-engineering/pkg"
)

func main() {
	// Global flags
	certFile := flag.String("cert", "", "Path to the certificate file")
	keyFile := flag.String("key", "", "Path to the key file")
	jsonOutput := flag.Bool("json", false, "Output in JSON format")
	flag.Usage = mainHelp
	flag.Parse()

	if flag.NArg() == 0 {
		mainHelp()
		os.Exit(0)
	}

	// Load certificate
	crt := "certificate.pem"
	key := "private.key"
	if *certFile != "" {
		crt = *certFile
	}
	if *keyFile != "" {
		key = *keyFile
	}

	if err := rewerse.SetCertificate(crt, key); err != nil {
		if *certFile == "" && *keyFile == "" {
			fmt.Fprintln(os.Stderr, "Certificate not found. Use -cert and -key flags or place certificate.pem and private.key in current directory.")
		} else {
			fmt.Fprintf(os.Stderr, "Error loading certificate: %v\n", err)
		}
		os.Exit(1)
	}

	var data any
	var err error

	switch flag.Arg(0) {
	case "markets":
		data, err = handleMarkets(flag.Args()[1:])
	case "products":
		data, err = handleProducts(flag.Args()[1:])
	case "recipes":
		data, err = handleRecipes(flag.Args()[1:], *jsonOutput)
		if data == nil && err == nil {
			return // already printed
		}
	case "discounts":
		data, err = handleDiscounts(flag.Args()[1:])
	case "categories":
		data, err = handleCategories(flag.Args()[1:], *jsonOutput)
		if data == nil {
			return // already printed
		}
	case "recalls":
		data, err = rewerse.GetRecalls()
	case "services":
		data, err = handleServices(flag.Args()[1:])
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", flag.Arg(0))
		mainHelp()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if *jsonOutput {
		bt, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(bt))
	} else {
		fmt.Println(data)
	}
}

func mainHelp() {
	fmt.Printf(`Usage: %s [flags] <command> [subcommand] [flags]

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
  %s markets search -query Köln
  %s products search -market 831002 -query Milch
  %s products category -market 831002 -slug obst-gemuese
  %s recipes search -term Pasta
  %s discounts -market 840174
  %s categories -market 831002
  %s services -zip 50667

Run '%s <command>' for subcommand help.
`, binaryName, binaryName, binaryName, binaryName, binaryName, binaryName, binaryName, binaryName, binaryName)
}
