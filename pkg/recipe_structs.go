package rewerse

import "fmt"

// RecipeSorting defines sorting options for recipe search
type RecipeSorting string

// Probably not exhaustive
const (
	SortRelevance RecipeSorting = "RELEVANCE_DESC"
)

// RecipeCollection defines recipe collection filters
type RecipeCollection string

// Not exhaustive - other collections may exist
const (
	CollectionVegetarisch RecipeCollection = "Vegetarisch"
)

// RecipeDifficulty defines difficulty level filters
type RecipeDifficulty string

const (
	DifficultyEasy   RecipeDifficulty = "Gering"
	DifficultyMedium RecipeDifficulty = "Mittel"
	DifficultyHard   RecipeDifficulty = "Hoch"
)

// RecipeSearchOpts contains options for recipe search
type RecipeSearchOpts struct {
	// SearchTerm is the text query to search for
	SearchTerm string
	// Collection filters by recipe collection
	Collection RecipeCollection
	// Difficulty filters by difficulty level
	Difficulty RecipeDifficulty
	// Sorting determines result order (default: SortRelevance)
	Sorting RecipeSorting
	// Page is the page number (1-indexed)
	Page int
	// ObjectsPerPage is the number of results per page (default 20)
	ObjectsPerPage int
}

// RecipeSearchResults contains paginated recipe search results
// Endpoint: GET /api/v3/recipe-search
type RecipeSearchResults struct {
	// TotalCount is the total number of matching recipes
	TotalCount int `json:"totalCount"`
	// Recipes contains the recipe summaries for this page
	Recipes []Recipe `json:"recipes"`
	// Metadata contains available filter options
	Metadata RecipeMetadata `json:"metadata"`
}

func (r RecipeSearchResults) String() string {
	s := fmt.Sprintf("Found %d recipes:\n", r.TotalCount)
	for _, recipe := range r.Recipes {
		s += recipe.String() + "\n"
	}
	return s
}

// RecipeMetadata contains available filter options for recipe search
type RecipeMetadata struct {
	// Collections are recipe categories: ["Vegetarisch"]
	Collections []string `json:"collections"`
	// Tags are recipe tags: ["Geringer Aufwand", "Abendessen", "schnell", "Vegan", ...]
	Tags []string `json:"tags"`
	// Difficulties are available difficulty levels: ["Gering", "Mittel", "Hoch"]
	Difficulties []string `json:"difficulties"`
}

// RecipeDetails contains full recipe information including ingredients and steps
// Endpoint: GET /api/v3/recipe-hub-details?recipeId={id}
type RecipeDetails struct {
	Recipe RecipeDetail `json:"recipe"`
}

func (rd RecipeDetails) String() string {
	return rd.Recipe.StringFull()
}

// RecipeDetail is the full recipe with ingredients and cooking steps
type RecipeDetail struct {
	// ID is the recipe UUID: "30ce3caf-4b3b-4c9e-8ea0-645fe75d1303"
	ID string `json:"id"`
	// Title is the recipe name: "Gnocchi-Pfanne mit Rosenkohl und getrockneten Tomaten"
	Title string `json:"title"`
	// DetailURL is the web page URL: "https://www.rewe.de/rezepte/..."
	DetailURL string `json:"detailUrl"`
	// ImageURL is the recipe image
	ImageURL string `json:"imageUrl"`
	// Duration is the cooking time: "50 min"
	Duration string `json:"duration"`
	// DifficultyLevel is numeric difficulty: 1=easy, 2=medium, 3=hard
	DifficultyLevel int `json:"difficultyLevel"`
	// DifficultyDescription is human-readable: "Einfach", "Mittel", "Schwer"
	DifficultyDescription string `json:"difficultyDescription"`
	// Ingredients contains the ingredient list with portions
	Ingredients RecipeIngredients `json:"ingredients"`
	// Steps contains the cooking instructions
	Steps []string `json:"steps"`
}

func (r RecipeDetail) StringFull() string {
	s := fmt.Sprintf("%s\n", r.Title)
	s += fmt.Sprintf("  Dauer: %s, Schwierigkeit: %s\n", r.Duration, r.DifficultyDescription)
	s += fmt.Sprintf("  URL: %s\n\n", r.DetailURL)

	s += fmt.Sprintf("Zutaten (für %d Portionen):\n", r.Ingredients.Portions)
	for _, ing := range r.Ingredients.Items {
		if ing.Quantity > 0 {
			s += fmt.Sprintf("  - %.0f %s %s\n", ing.Quantity, ing.Unit, ing.Name)
		} else {
			s += fmt.Sprintf("  - %s\n", ing.Name)
		}
	}

	s += "\nZubereitung:\n"
	for i, step := range r.Steps {
		s += fmt.Sprintf("  %d. %s\n", i+1, step)
	}

	return s
}

// RecipeIngredients contains the ingredient list with portion info
type RecipeIngredients struct {
	// Portions is the number of servings the recipe makes
	Portions int `json:"portions"`
	// Items contains the individual ingredients
	Items []RecipeIngredient `json:"items"`
}

// RecipeIngredient is a single ingredient with quantity and unit
type RecipeIngredient struct {
	// Name is the ingredient name: "frischer Rosenkohl (ersatzweise TK)"
	Name string `json:"name"`
	// Quantity is the amount needed: 250, 1, 0 (for "to taste")
	Quantity float64 `json:"quantity"`
	// Unit is the measurement unit: "g", "EL", "Zehe(n)", "" (for count)
	Unit string `json:"unit"`
}

// PopularSearchTerm is a suggested search term
// Endpoint: GET /api/v3/recipe-popular-search-terms
type PopularSearchTerm struct {
	// ID is the search term identifier
	ID string `json:"id"`
	// Title is the display text: "Lachs", "Low Carb", "Kürbis", ...
	Title string `json:"title"`
}

type PopularSearchTerms []PopularSearchTerm

func (p PopularSearchTerms) String() string {
	s := "Popular search terms:\n"
	for _, term := range p {
		s += fmt.Sprintf("  - %s\n", term.Title)
	}
	return s
}
