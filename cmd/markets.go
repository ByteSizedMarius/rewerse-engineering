package main

import (
	"flag"
	"fmt"

	rewerse "github.com/ByteSizedMarius/rewerse-engineering/pkg"
)

func handleMarkets(args []string) (any, error) {
	if wantsHelp(args) {
		marketsHelp()
		return nil, nil
	}

	switch args[0] {
	case "search":
		fs := flag.NewFlagSet("markets search", flag.ContinueOnError)
		query := fs.String("query", "", "Search query (city, zip, street)")
		if err := fs.Parse(args[1:]); err != nil {
			return nil, err
		}
		if err := checkUnexpectedArgs(fs); err != nil {
			return nil, err
		}
		if err := validateFlag("query", *query); err != nil {
			return nil, err
		}
		return rewerse.MarketSearch(*query)

	case "details":
		fs := flag.NewFlagSet("markets details", flag.ContinueOnError)
		id := fs.String("id", "", "Market ID")
		if err := fs.Parse(args[1:]); err != nil {
			return nil, err
		}
		if err := checkUnexpectedArgs(fs); err != nil {
			return nil, err
		}
		if err := validateNumeric("id", *id); err != nil {
			return nil, err
		}
		return rewerse.GetMarketDetails(*id)

	default:
		marketsHelp()
		return nil, fmt.Errorf("unknown markets subcommand: %s", args[0])
	}
}

func marketsHelp() {
	fmt.Printf(`Usage: %s markets <subcommand> [flags]

Subcommands:
  search      Search for markets
  details     Get market details

markets search:
  -query      Search query (city, zip code, street)

markets details:
  -id         Market ID

Examples:
  %s markets search -query Mannheim
  %s markets search -query 68199
  %s markets details -id 840174
`, binaryName, binaryName, binaryName, binaryName)
}
