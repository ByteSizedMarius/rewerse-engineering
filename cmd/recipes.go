package main

import (
	"flag"
	"fmt"

	rewerse "github.com/ByteSizedMarius/rewerse-engineering/pkg"
)

func handleRecipes(args []string, jsonOutput bool) (any, error) {
	if wantsHelp(args) {
		recipesHelp()
		return nil, nil
	}

	switch args[0] {
	case "search":
		fs := flag.NewFlagSet("recipes search", flag.ContinueOnError)
		term := fs.String("term", "", "Search term")
		collection := fs.String("collection", "", "Collection filter (Vegetarisch, Vegan)")
		difficulty := fs.String("difficulty", "", "Difficulty (Gering, Mittel, Hoch)")
		page := fs.Int("page", 0, "Page number")
		perPage := fs.Int("perPage", 0, "Results per page")
		all := fs.Bool("all", false, "Fetch all results")
		if err := fs.Parse(args[1:]); err != nil {
			return nil, err
		}
		if err := checkUnexpectedArgs(fs); err != nil {
			return nil, err
		}

		opts := &rewerse.RecipeSearchOpts{
			SearchTerm:     *term,
			Collection:     rewerse.RecipeCollection(*collection),
			Difficulty:     rewerse.RecipeDifficulty(*difficulty),
			Page:           *page,
			ObjectsPerPage: *perPage,
		}

		results, err := rewerse.RecipeSearch(opts)
		if err != nil {
			return nil, err
		}

		// JSON output - return data for cli.go to handle
		if jsonOutput && !*all {
			return results, nil
		}

		// Effective values (library defaults: page=1, perPage=20)
		effectivePage := opts.Page
		if effectivePage == 0 {
			effectivePage = 1
		}
		effectivePerPage := opts.ObjectsPerPage
		if effectivePerPage == 0 {
			effectivePerPage = 20
		}

		// Fetch all pages if requested
		if *all {
			allRecipes := results.Recipes
			opts.Page = effectivePage
			for len(allRecipes) < results.TotalCount {
				opts.Page++
				more, err := rewerse.RecipeSearch(opts)
				if err != nil {
					return nil, err
				}
				if len(more.Recipes) == 0 {
					break
				}
				allRecipes = append(allRecipes, more.Recipes...)
			}
			results.Recipes = allRecipes

			if jsonOutput {
				return results, nil
			}

			fmt.Printf("Showing all %d recipes:\n", len(allRecipes))
			for _, recipe := range allRecipes {
				fmt.Println(recipe.String())
			}
			return nil, nil
		}

		// Paginated output with page info
		start := (effectivePage-1)*effectivePerPage + 1
		end := start + len(results.Recipes) - 1
		fmt.Printf("Showing %d-%d of %d recipes (page %d):\n", start, end, results.TotalCount, effectivePage)
		for _, recipe := range results.Recipes {
			fmt.Println(recipe.String())
		}
		return nil, nil

	case "details":
		fs := flag.NewFlagSet("recipes details", flag.ContinueOnError)
		id := fs.String("id", "", "Recipe UUID")
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
		return rewerse.GetRecipePopularTerms()

	case "hub":
		return rewerse.GetRecipeHub()

	default:
		recipesHelp()
		return nil, fmt.Errorf("unknown recipes subcommand: %s", args[0])
	}
}

func recipesHelp() {
	fmt.Printf(`Usage: %s recipes <subcommand> [flags]

Subcommands:
  search      Search for recipes
  details     Get recipe details
  popular     Get popular search terms
  hub         Get recipe hub (featured recipes)

recipes search:
  -term       Search term
  -collection Collection filter (e.g. Vegetarisch)
  -difficulty Difficulty (Gering, Mittel, Hoch)
  -page       Page number
  -perPage    Results per page
  -all        Fetch all results

recipes details:
  -id         Recipe UUID (required)

Examples:
  %s recipes search -term Pasta
  %s recipes search -collection Vegetarisch -difficulty Mittel
  %s recipes details -id 30ce3caf-4b3b-4c9e-8ea0-645fe75d1303
  %s recipes popular
  %s recipes hub
`, binaryName, binaryName, binaryName, binaryName, binaryName, binaryName)
}
