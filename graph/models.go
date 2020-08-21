package graph

// Book model
type Book struct {
	ID          string   `json:"id,omitempty"`
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Authors     []Author `json:"authors,omitempty"`
}

// Author Model
type Author struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	ISBNNo string `json:"isbn_no,omitempty"`
	Books  []Book `json:"books,omitempty"`
}
