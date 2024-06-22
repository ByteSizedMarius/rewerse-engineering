package rewerse

import "fmt"

// Recalls is the struct for Rewe Product-Recalls
type Recalls []Recall

func (rs Recalls) String() string {
	if len(rs) == 0 {
		return "Aktuell keine Produktrückrufe"
	}

	recalls := "Produktrückrufe:\n"
	for _, r := range rs {
		recalls += r.String() + "\n"
	}

	return recalls
}

// Recall is the struct for a single recall
type Recall struct {
	RecallURL      string `json:"recallUrl"`
	SubjectProduct string `json:"subjectProduct"`
	SubjectReason  string `json:"subjectReason"`
}

func (r Recall) String() string {
	return fmt.Sprintf("%s\n%s\n%s\n", r.SubjectProduct, r.SubjectReason, r.RecallURL)
}

// GetRecalls returns all currently ongoing recalls from Rewe
func GetRecalls() (r Recalls, err error) {
	req, err := BuildCustomRequest(apiHost, "v3/productrecalls")
	if err != nil {
		return
	}

	err = DoRequest(req, &r)
	if err != nil {
		return
	}

	return
}

// RecipeHub is the struct for the Data returned by the Rewe Recipe-Page
type RecipeHub struct {
	RecipeOfTheDay Recipe   `json:"recipeOfTheDay"`
	PopularRecipes []Recipe `json:"popularRecipes"`
	Categories     []struct {
		Type        string `json:"type"`
		Title       string `json:"title"`
		SearchQuery string `json:"searchQuery"`
	} `json:"categories"`
}

func (rh RecipeHub) String() string {
	recipeHub := "Recipe Hub\n\n"
	recipeHub += "Rezept des Tages\n--------------------\n" + rh.RecipeOfTheDay.String() + "\n"

	recipeHub += "Beliebte Rezepte\n--------------------\n"
	for _, r := range rh.PopularRecipes {
		recipeHub += r.String() + "\n"
	}

	recipeHub += "Verfügbare Rezept-Kategorien\n--------------------\n"
	for _, c := range rh.Categories {
		recipeHub += c.Title + "\n"
	}

	return recipeHub
}

// Recipe is the struct for a single recipe
type Recipe struct {
	ID                    string `json:"id"`
	Title                 string `json:"title"`
	DetailURL             string `json:"detailUrl"`
	ImageURL              string `json:"imageUrl"`
	Duration              string `json:"duration"`
	DifficultyLevel       int    `json:"difficultyLevel"`
	DifficultyDescription string `json:"difficultyDescription"`
}

func (r Recipe) String() string {
	return fmt.Sprintf("%s\n%s\n%s\n%s\n", r.Title, r.Duration, r.DifficultyDescription, r.DetailURL)
}

// GetRecipeHub returns the Data from the RecipeHub
func GetRecipeHub() (r RecipeHub, err error) {
	req, err := BuildCustomRequest(apiHost, "v3/recipe-hub")
	if err != nil {
		return
	}

	err = DoRequest(req, &r)
	if err != nil {
		return
	}

	return
}
