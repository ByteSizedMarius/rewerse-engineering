package rewerse

import (
	"fmt"
	"strconv"
	"strings"
)

type RecipeDifficulty int

const (
	DifficultyEasy   RecipeDifficulty = 1
	DifficultyMedium RecipeDifficulty = 2
	DifficultyHard   RecipeDifficulty = 3
)

func (d RecipeDifficulty) String() string {
	switch d {
	case DifficultyEasy:
		return "Gering"
	case DifficultyMedium:
		return "Mittel"
	case DifficultyHard:
		return "Hoch"
	default:
		return strconv.Itoa(int(d))
	}
}

// RecipeCollection is a recipe category usable as a search filter.
type RecipeCollection string

// The complete set of collections the recipe search accepts.
const (
	CollectionVegetarian RecipeCollection = "Vegetarisch"
	CollectionMeat       RecipeCollection = "Fleisch"
	CollectionBaking     RecipeCollection = "Backen"
	CollectionDesserts   RecipeCollection = "Nachspeisen"
	CollectionFish       RecipeCollection = "Fisch"
	CollectionStarters   RecipeCollection = "Vorspeisen"
	CollectionCakes      RecipeCollection = "Kuchen"
)

// TagConcat controls how multiple tags combine in a recipe search.
type TagConcat string

const (
	TagConcatAnd TagConcat = "AND" // intersects the tag result sets
	TagConcatOr  TagConcat = "OR"  // unions the tag result sets
)

type recipeSearchResponse struct {
	Data struct {
		Recipes struct {
			Recipes  []RecipeSummary `json:"recipes"`
			Metadata RecipeMetadata  `json:"metadata"`
		} `json:"recipes"`
		RecipeOfTheDay *Recipe `json:"recipeOfTheDay"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type recipeDetailsResponse struct {
	Data struct {
		RecipeByID *Recipe `json:"recipeById"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type recipePopularTermsResponse struct {
	Data struct {
		SearchTags struct {
			Tags []struct {
				Title string `json:"title"`
			} `json:"tags"`
		} `json:"searchTags"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type RecipeSearchResults struct {
	Recipes  []RecipeSummary `json:"recipes"`
	Metadata RecipeMetadata  `json:"metadata"`
	// RecipeOfTheDay is set only when RecipeSearchOpts.IncludeRecipeOfTheDay is true
	RecipeOfTheDay *Recipe `json:"recipeOfTheDay"`
}

func (rs RecipeSearchResults) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d recipes (page shows %d)\n", rs.Metadata.TotalRecipeCount, len(rs.Recipes))
	if rs.RecipeOfTheDay != nil {
		fmt.Fprintf(&b, "Recipe of the day: %s\n", rs.RecipeOfTheDay.Title)
	}
	for _, r := range rs.Recipes {
		b.WriteString(r.String() + "\n")
	}
	return b.String()
}

// RecipeSummary is a recipe as returned by the search.
type RecipeSummary struct {
	// ID is what GetRecipeDetails takes: "blt1111111111111111"
	ID    string `json:"id"`
	Title string `json:"title"`
	// TimePreparation is in minutes, null for some recipes
	TimePreparation *int `json:"timePreparation"`
	// TimeCooking is in minutes, null for some recipes
	TimeCooking *int `json:"timeCooking"`
	// TimeTotal is in minutes, null for some recipes
	TimeTotal  *int             `json:"timeTotal"`
	Difficulty RecipeDifficulty `json:"difficulty"`
	// DetailURL is the full web URL: "https://www.rewe.de/rezepte/nudeln-mit-tomatensosse/"
	DetailURL string   `json:"detailUrl"`
	ImageURLs []string `json:"imageUrls"`
}

func (r RecipeSummary) String() string {
	return fmt.Sprintf("%s | %s/%s/%s min | %s",
		r.Title, minutes(r.TimePreparation), minutes(r.TimeCooking), minutes(r.TimeTotal), r.Difficulty)
}

// RecipeMetadata holds the facet counts of a recipe search.
type RecipeMetadata struct {
	// TotalRecipeCount is the number of matches across all pages: 275
	TotalRecipeCount int                     `json:"totalRecipeCount"`
	Collections      []RecipeCollectionCount `json:"collections"`
	Difficulties     []RecipeDifficultyCount `json:"difficulties"`
	// Tags was empty in every observed response
	Tags []string `json:"tags"`
}

// RecipeCollectionCount is the number of matches in one collection.
type RecipeCollectionCount struct {
	Count int              `json:"count"`
	Name  RecipeCollection `json:"name"`
}

// RecipeDifficultyCount is the number of matches at one difficulty.
type RecipeDifficultyCount struct {
	Count int              `json:"count"`
	Name  RecipeDifficulty `json:"name"`
}

// Recipe is a full recipe with ingredients, steps and nutrients.
type Recipe struct {
	ID         string           `json:"id"`
	Title      string           `json:"title"`
	Difficulty RecipeDifficulty `json:"difficulty"`
	Serving    RecipeServing    `json:"serving"`
	// TimeCooking is in minutes
	TimeCooking int `json:"timeCooking"`
	// TimePreparation is in minutes
	TimePreparation int `json:"timePreparation"`
	// TimeTotal is in minutes, null for some recipes
	TimeTotal   *int              `json:"timeTotal"`
	Preparation RecipePreparation `json:"preparation"`
	ImageURLs   []RecipeImage     `json:"imageUrls"`
	// DetailURL is only the slug: "nudeln-mit-tomatensosse"
	DetailURL string `json:"detailUrl"`
	// Tags are display labels: "Snacks", "Dessert", "Mittlerer Aufwand", "Hefeteig"
	Tags        []string           `json:"tags"`
	Ingredients []RecipeIngredient `json:"ingredients"`
	// Nutrients are the macro values per serving: "Energie", "Eiweiß"
	Nutrients []RecipeNutrient `json:"nutrients"`
	// VitalSubstances are the vitamin and mineral values per serving: "Vitamin B2", "Kalium"
	VitalSubstances []RecipeNutrient `json:"vitalSubstances"`
}

func (r Recipe) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s\n", r.Title)
	fmt.Fprintf(&b, "  Serving: %d %s\n", r.Serving.Quantity, r.Serving.Type)
	fmt.Fprintf(&b, "  Time: %d prep / %d cooking / %s total (min)\n",
		r.TimePreparation, r.TimeCooking, minutes(r.TimeTotal))
	fmt.Fprintf(&b, "  Difficulty: %s\n", r.Difficulty)

	b.WriteString("  Ingredients:\n")
	for _, i := range r.Ingredients {
		fmt.Fprintf(&b, "    - %s\n", i.DisplayName)
	}

	b.WriteString("  Steps:\n")
	for n, s := range r.Preparation.Steps {
		fmt.Fprintf(&b, "    %d. %s\n", n+1, s.Description)
	}

	return b.String()
}

// RecipeServing is the yield of a recipe.
type RecipeServing struct {
	// Type is the unit of the yield: "portions", "pieces"
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type RecipePreparation struct {
	Steps []RecipeStep `json:"steps"`
}

type RecipeStep struct {
	Description string `json:"description"`
}

type RecipeImage struct {
	URLs struct {
		Default string `json:"default"`
	} `json:"urls"`
}

type RecipeIngredient struct {
	Quantity float64 `json:"quantity"`
	// Unit is the measure: "g", "ml", "Würfel". Null for countable items like "3 Eier"
	Unit *string `json:"unit"`
	Name string  `json:"name"`
	// DisplayName is the formatted line: "200 ml Milch", "3 Eier"
	DisplayName string `json:"displayName"`
}

// RecipeNutrient is one nutrient value per serving.
type RecipeNutrient struct {
	Title string `json:"title"`
	Unit  string `json:"unit"`
	// RDA is the recommended daily amount
	RDA    float64 `json:"rda"`
	Amount float64 `json:"amount"`
}

// minutes renders a null value as "-".
func minutes(m *int) string {
	if m == nil {
		return "-"
	}
	return strconv.Itoa(*m)
}
