package graph

import (
	"github.com/graphql-go/graphql"
)

var authorType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "author",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
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
		Name: "book",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var listBookType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "book",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"author": &graphql.Field{
				Type: authorType,
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: getBooksOfAuthorResolver,
			},
		},
	},
)

func init() {
	bookType.AddFieldConfig("authors", &graphql.Field{
		Type:    graphql.NewList(authorType),
		Resolve: getAuthorsResolver,
	})
	authorType.AddFieldConfig("books", &graphql.Field{
		Type:    graphql.NewList(bookType),
		Resolve: getBooksResolver,
	})
}
