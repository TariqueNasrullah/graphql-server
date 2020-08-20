package graph

// Book model
type Book struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Authors     []Author `json:"authors"`
}

// Author Model
type Author struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	ISBNNo string `json:"isbn_no"`
	Books  []Book `jsoin:"books"`
}
