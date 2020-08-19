package graph

// Book model
type Book struct {
	ID          string `json:"Id"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
	Author      Author `json:"Author"`
}

// Author Model
type Author struct {
	ID     string `json:"Id"`
	Name   string `json:"Name"`
	ISBNNo string `json:"isbn_no"`
	Books  []Book `json:"Books"`
}
