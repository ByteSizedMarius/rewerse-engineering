package rewerse

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

var ErrRecipeNotFound = errors.New("recipe not found")

var defaultRecipeOpts = RecipeSearchOpts{
	Page:           1,
	ObjectsPerPage: 20,
	TagConcat:      TagConcatAnd,
}

type RecipeSearchOpts struct {
	// SearchTerm matches everything when empty or "*"
	SearchTerm string
	// Page is 1-based
	Page           int
	ObjectsPerPage int
	Collections    []RecipeCollection
	Difficulties   []RecipeDifficulty
	// Tags are lowercase labels: "fisch", "gesund", "vegan"
	Tags                  []string
	TagConcat             TagConcat
	IncludeRecipeOfTheDay bool
}

// RecipeSearch searches recipes by term, collection, difficulty and tags.
// Endpoint: GET /api/recipes/search
// A nil opts searches all recipes with the defaults.
func RecipeSearch(opts *RecipeSearchOpts) (RecipeSearchResults, error) {
	if opts == nil {
		opts = &defaultRecipeOpts
	} else {
		if opts.Page <= 0 {
			opts.Page = defaultRecipeOpts.Page
		}
		if opts.ObjectsPerPage <= 0 {
			opts.ObjectsPerPage = defaultRecipeOpts.ObjectsPerPage
		}
		if opts.TagConcat == "" {
			opts.TagConcat = defaultRecipeOpts.TagConcat
		}
	}

	query := url.Values{}
	query.Add("searchTerm", opts.SearchTerm)
	query.Add("page", strconv.Itoa(opts.Page))
	query.Add("objectsPerPage", strconv.Itoa(opts.ObjectsPerPage))
	query.Add("tagConcat", string(opts.TagConcat))
	for _, c := range opts.Collections {
		query.Add("collections", string(c))
	}
	for _, d := range opts.Difficulties {
		query.Add("difficulties", strconv.Itoa(int(d)))
	}
	for _, t := range opts.Tags {
		query.Add("tags", t)
	}
	if opts.IncludeRecipeOfTheDay {
		query.Add("includeRecipeOfTheDay", "true")
	}

	req, err := BuildCustomRequest(clientHost, "recipes/search?"+query.Encode())
	if err != nil {
		return RecipeSearchResults{}, err
	}

	setCommonHeaders(req)

	var res recipeSearchResponse
	err = DoRequest(req, &res)
	if err != nil {
		return RecipeSearchResults{}, err
	}
	if len(res.Errors) > 0 {
		return RecipeSearchResults{}, fmt.Errorf("api error: %s", res.Errors[0].Message)
	}

	return RecipeSearchResults{
		Recipes:        res.Data.Recipes.Recipes,
		Metadata:       res.Data.Recipes.Metadata,
		RecipeOfTheDay: res.Data.RecipeOfTheDay,
	}, nil
}

// GetRecipeDetails returns a full recipe with ingredients, steps and nutrients.
// Endpoint: GET /api/recipes/{recipeId}
func GetRecipeDetails(recipeID string) (Recipe, error) {
	if recipeID == "" {
		return Recipe{}, fmt.Errorf("recipeID: cannot be empty")
	}

	req, err := BuildCustomRequest(clientHost, "recipes/"+recipeID)
	if err != nil {
		return Recipe{}, err
	}

	setCommonHeaders(req)

	var res recipeDetailsResponse
	err = DoRequest(req, &res)
	if err != nil {
		return Recipe{}, err
	}
	if len(res.Errors) > 0 {
		return Recipe{}, fmt.Errorf("api error: %s", res.Errors[0].Message)
	}
	if res.Data.RecipeByID == nil {
		return Recipe{}, fmt.Errorf("%w: %s", ErrRecipeNotFound, recipeID)
	}

	return *res.Data.RecipeByID, nil
}

// GetRecipePopularTerms returns the search terms the app shows on the recipe landing page.
// Endpoint: GET /api/recipes/search/landing-page
func GetRecipePopularTerms() ([]string, error) {
	req, err := BuildCustomRequest(clientHost, "recipes/search/landing-page")
	if err != nil {
		return nil, err
	}

	setCommonHeaders(req)

	var res recipePopularTermsResponse
	err = DoRequest(req, &res)
	if err != nil {
		return nil, err
	}
	if len(res.Errors) > 0 {
		return nil, fmt.Errorf("api error: %s", res.Errors[0].Message)
	}

	terms := make([]string, 0, len(res.Data.SearchTags.Tags))
	for _, t := range res.Data.SearchTags.Tags {
		terms = append(terms, t.Title)
	}
	return terms, nil
}
