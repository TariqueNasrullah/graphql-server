package graph

import (
	"context"

	"github.com/TariqueNasrullah/graphql-server/database"
	"github.com/graphql-go/graphql"
)

// GraphQL Mutation Type
var mutationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		// book - Creates a Book
		"book": &graphql.Field{
			Type:        bookType,
			Description: "Create new book",
			Args: graphql.FieldConfigArgument{
				"title": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"description": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"authors": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.NewList(graphql.String)),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// Get Argument values
				title := p.Args["title"].(string)
				description := p.Args["description"].(string)
				authorIds := p.Args["authors"].([]interface{})

				// Create List of Author from argument - 'authors'
				var authors []Author
				for _, authorID := range authorIds {
					author := Author{
						ID: authorID.(string),
					}
					authors = append(authors, author)
				}

				// Create Book instance
				book := Book{
					Title:       title,
					Description: description,
					Authors:     authors,
				}

				// Write document to database
				meta, err := database.BookCollection.CreateDocument(context.Background(), book)
				if err != nil {
					return nil, err
				}
				book.ID = meta.Key

				return book, nil
			},
		},
		// author - Creates an Author
		"author": &graphql.Field{
			Type:        authorType,
			Description: "Create new author",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"isbn_no": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// Get Argument values
				name := p.Args["name"].(string)
				isbnNo := p.Args["isbn_no"].(string)

				// Construct new Author instance
				author := Author{
					Name:   name,
					ISBNNo: isbnNo,
				}

				// Write document to database
				meta, err := database.AuthorCollecction.CreateDocument(context.Background(), author)
				if err != nil {
					return nil, err
				}
				author.ID = meta.Key

				return author, nil
			},
		},
	},
})
