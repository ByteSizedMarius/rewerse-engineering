package rewerse

import (
	"net/url"
	"strconv"
)

var defaultRecipeOpts = RecipeSearchOpts{
	Sorting:        SortRelevance,
	Page:           1,
	ObjectsPerPage: 20,
}

// RecipeSearch searches for recipes with optional filters.
// Endpoint: GET /api/v3/recipe-search
func RecipeSearch(opts *RecipeSearchOpts) (results RecipeSearchResults, err error) {
	if opts == nil {
		opts = &defaultRecipeOpts
	} else {
		if opts.Sorting == "" {
			opts.Sorting = defaultRecipeOpts.Sorting
		}
		if opts.Page <= 0 {
			opts.Page = defaultRecipeOpts.Page
		}
		if opts.ObjectsPerPage <= 0 {
			opts.ObjectsPerPage = defaultRecipeOpts.ObjectsPerPage
		}
	}

	query := url.Values{}
	query.Add("searchTerm", opts.SearchTerm)
	query.Add("sorting", string(opts.Sorting))
	query.Add("page", strconv.Itoa(opts.Page))
	query.Add("objectsPerPage", strconv.Itoa(opts.ObjectsPerPage))

	if opts.Collection != "" {
		query.Add("collection", string(opts.Collection))
	}
	if opts.Difficulty != "" {
		query.Add("difficulty", string(opts.Difficulty))
	}

	req, err := BuildCustomRequest(apiHost, "v3/recipe-search?"+query.Encode())
	if err != nil {
		return
	}

	err = DoRequest(req, &results)
	return
}

// GetRecipeDetails returns the full recipe with ingredients and steps.
// Endpoint: GET /api/v3/recipe-hub-details?recipeId={id}
func GetRecipeDetails(recipeID string) (details RecipeDetails, err error) {
	query := url.Values{}
	query.Add("recipeId", recipeID)

	req, err := BuildCustomRequest(apiHost, "v3/recipe-hub-details?"+query.Encode())
	if err != nil {
		return
	}

	err = DoRequest(req, &details)
	return
}

// GetRecipePopularTerms returns popular recipe search terms.
// Endpoint: GET /api/v3/recipe-popular-search-terms
func GetRecipePopularTerms() (terms PopularSearchTerms, err error) {
	req, err := BuildCustomRequest(apiHost, "v3/recipe-popular-search-terms")
	if err != nil {
		return
	}

	err = DoRequest(req, &terms)
	return
}
