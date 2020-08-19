package graph

import (
	"fmt"

	"github.com/graphql-go/graphql"
)

var authorType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Author",
		Fields: graphql.Fields{
			"Id": &graphql.Field{
				Type: graphql.String,
			},
			"Name": &graphql.Field{
				Type: graphql.String,
			},
			"isbn_no": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var bookType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Book",
		Fields: graphql.Fields{
			"Id": &graphql.Field{
				Type: graphql.String,
			},
			"Title": &graphql.Field{
				Type: graphql.String,
			},
			"Description": &graphql.Field{
				Type: graphql.String,
			},
			"Author": &graphql.Field{
				Type: authorType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					fmt.Println(p.Source)
					return author, nil
				},
			},
		},
	},
)

var qureyType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"Book": &graphql.Field{
				Type: bookType,
				Args: graphql.FieldConfigArgument{
					"Id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					idQuery, ok := p.Args["Id"].(string)
					if ok && book.ID == idQuery {
						return book, nil
					}
					return nil, nil
				},
			},
		},
	},
)

var book = Book{
	ID:          "1",
	Title:       "harry potter",
	Description: "A legendary trealer",
}

var author = Author{
	ID:     "1",
	Name:   "Tarique",
	ISBNNo: "12345",
}

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: qureyType,
	},
)

func init() {
	authorType.AddFieldConfig("Books", &graphql.Field{Type: graphql.NewList(authorType)})
}

// ExecuteQuery executes the query
func ExecuteQuery(query string) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}
