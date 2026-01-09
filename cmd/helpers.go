package main

import (
	"flag"
	"fmt"
	"regexp"
	"strings"

	rewerse "github.com/ByteSizedMarius/rewerse-engineering/pkg"
)

var (
	numericRegex = regexp.MustCompile(`^\d+$`)
	zipRegex     = regexp.MustCompile(`^\d{5}$`)
)

// wantsHelp checks if args request help (empty args, "help", "-h", "--help")
func wantsHelp(args []string) bool {
	if len(args) == 0 {
		return true
	}
	switch args[0] {
	case "help", "-h", "--help":
		return true
	}
	return false
}

// validateFlag returns an error if the flag value is empty or whitespace-only
func validateFlag(name, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("-%s is required", name)
	}
	return nil
}

// validateNumeric returns an error if the value is empty or non-numeric
func validateNumeric(name, value string) error {
	if err := validateFlag(name, value); err != nil {
		return err
	}
	if !numericRegex.MatchString(value) {
		return fmt.Errorf("-%s must be numeric (got %q)", name, value)
	}
	return nil
}

// validateZipCode returns an error if the zip code is not 5 digits
func validateZipCode(value string) error {
	if err := validateFlag("zip", value); err != nil {
		return err
	}
	if !zipRegex.MatchString(value) {
		return fmt.Errorf("-zip must be 5 digits (got %q)", value)
	}
	return nil
}

// buildOpts creates ProductOpts from pagination and service params
func buildOpts(page, perPage int, service string) *rewerse.ProductOpts {
	if page == 0 && perPage == 0 && service == "" {
		return nil
	}
	return &rewerse.ProductOpts{
		Page:           page,
		ObjectsPerPage: perPage,
		ServiceType:    rewerse.ServiceType(service),
	}
}

// checkUnexpectedArgs returns an error if there are leftover positional arguments
func checkUnexpectedArgs(fs *flag.FlagSet) error {
	if fs.NArg() > 0 {
		return fmt.Errorf("unexpected argument: %s (use quotes for multi-word values)", fs.Arg(0))
	}
	return nil
}
