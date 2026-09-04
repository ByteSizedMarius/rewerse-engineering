package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	rewerse "github.com/ByteSizedMarius/rewerse-engineering/pkg"
)

// popularTerms prints one search term per line.
type popularTerms []string

func (p popularTerms) String() string {
	return strings.Join(p, "\n")
}

func handleRecipes(args []string) (any, error) {
	if wantsHelp(args) {
		recipesHelp()
		return nil, nil
	}

	switch args[0] {
	case "search":
		fs := flag.NewFlagSet("recipes search", flag.ContinueOnError)
		term := fs.String("term", "", "Search term")
		collections := fs.String("collections", "", "Comma-separated collections")
		difficulties := fs.String("difficulties", "", "Comma-separated difficulties: 1, 2, 3")
		tags := fs.String("tags", "", "Comma-separated lowercase tags")
		tagConcat := fs.String("tagConcat", "AND", "Tag combination: AND or OR")
		page := fs.Int("page", 0, "Page number")
		perPage := fs.Int("perPage", 0, "Results per page")
		rotd := fs.Bool("rotd", false, "Include the recipe of the day")
		if err := fs.Parse(args[1:]); err != nil {
			return nil, err
		}
		if err := checkUnexpectedArgs(fs); err != nil {
			return nil, err
		}

		opts := rewerse.RecipeSearchOpts{
			SearchTerm:            *term,
			Page:                  *page,
			ObjectsPerPage:        *perPage,
			Tags:                  splitList(*tags),
			IncludeRecipeOfTheDay: *rotd,
		}
		for _, c := range splitList(*collections) {
			opts.Collections = append(opts.Collections, rewerse.RecipeCollection(c))
		}
		for _, d := range splitList(*difficulties) {
			n, err := strconv.Atoi(d)
			if err != nil || n < 1 || n > 3 {
				return nil, fmt.Errorf("-difficulties must be 1, 2 or 3 (got %q)", d)
			}
			opts.Difficulties = append(opts.Difficulties, rewerse.RecipeDifficulty(n))
		}
		switch rewerse.TagConcat(*tagConcat) {
		case rewerse.TagConcatAnd, rewerse.TagConcatOr:
			opts.TagConcat = rewerse.TagConcat(*tagConcat)
		default:
			return nil, fmt.Errorf("-tagConcat must be AND or OR (got %q)", *tagConcat)
		}

		return rewerse.RecipeSearch(&opts)

	case "details":
		fs := flag.NewFlagSet("recipes details", flag.ContinueOnError)
		id := fs.String("id", "", "Recipe ID")
		if err := fs.Parse(args[1:]); err != nil {
			return nil, err
		}
		if err := checkUnexpectedArgs(fs); err != nil {
			return nil, err
		}
		if err := validateFlag("id", *id); err != nil {
			return nil, err
		}
		return rewerse.GetRecipeDetails(*id)

	case "popular":
		fs := flag.NewFlagSet("recipes popular", flag.ContinueOnError)
		if err := fs.Parse(args[1:]); err != nil {
			return nil, err
		}
		if err := checkUnexpectedArgs(fs); err != nil {
			return nil, err
		}
		terms, err := rewerse.GetRecipePopularTerms()
		if err != nil {
			return nil, err
		}
		return popularTerms(terms), nil

	default:
		recipesHelp()
		return nil, fmt.Errorf("unknown recipes subcommand: %s", args[0])
	}
}

// splitList drops empty and whitespace-only entries.
func splitList(value string) []string {
	var out []string
	for _, part := range strings.Split(value, ",") {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func recipesHelp() {
	fmt.Printf(`Usage: %s recipes <subcommand> [flags]

Subcommands:
  search      Search recipes
  details     Get recipe details
  popular     Get the popular search terms of the recipe landing page

recipes search:
  -term          Search term (empty matches all recipes)
  -collections   Comma-separated: Vegetarisch, Fleisch, Backen, Nachspeisen, Fisch, Vorspeisen, Kuchen
  -difficulties  Comma-separated: 1 (Gering), 2 (Mittel), 3 (Hoch)
  -tags          Comma-separated lowercase tags: fisch, gesund, vegan
  -tagConcat     AND or OR (default: AND)
  -page          Page number
  -perPage       Results per page
  -rotd          Include the recipe of the day

recipes details:
  -id            Recipe ID (required)

Examples:
  %s recipes search -term Lachs -difficulties 1
  %s recipes search -collections Fisch,Fleisch -tags gesund
  %s recipes details -id blt4aaa7361ba69f8c8
  %s recipes popular
`, binaryName, binaryName, binaryName, binaryName, binaryName)
}
