package graph

import (
	"context"
	"fmt"

	"github.com/TariqueNasrullah/graphql-server/database"
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
)

// GraphQL query type
var qureyType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "GetBooksOfAuthor",
		Fields: graphql.Fields{
			"book": &graphql.Field{
				Type: graphql.NewList(listBookType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Get AuthorName from the context
					authorName := p.Context.Value("authorName").(string)

					// Construct query to get bookList of the author
					query := fmt.Sprintf(`for author in Author
						filter author.name=='%s'
						for book in Book
						for authid in book.authors
						filter authid.id==author._key
						return book`, authorName)

					// Execute query
					bindVars := map[string]interface{}{}
					coursor, err := database.Db.Query(context.Background(), query, bindVars)
					if err != nil {
						logrus.Errorln(err)
						return nil, err
					}
					defer coursor.Close()

					// Extract books from query result
					var books []Book
					for coursor.HasMore() {
						var b Book
						meta, _ := coursor.ReadDocument(context.Background(), &b)
						b.ID = meta.Key
						books = append(books, b)
					}
					return books, nil
				},
			},
			// books - returns all Books currently in the database
			"books": &graphql.Field{
				Type: graphql.NewList(bookType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Construct query
					query := "for book in Book return book"

					// Execute query
					bindVars := map[string]interface{}{}
					coursor, err := database.Db.Query(context.Background(), query, bindVars)
					if err != nil {
						logrus.Errorln(err)
						return nil, err
					}
					defer coursor.Close()

					// Extract books from the query result
					var books []Book
					for coursor.HasMore() {
						var b Book
						meta, _ := coursor.ReadDocument(context.Background(), &b)
						b.ID = meta.Key
						books = append(books, b)
					}
					return books, nil
				},
			},
			// authors - returns all Authors in the database
			"authors": &graphql.Field{
				Type: graphql.NewList(authorType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Construct query
					query := "for author in Author return author"

					// Execute qury
					bindVars := map[string]interface{}{}
					coursor, err := database.Db.Query(context.Background(), query, bindVars)
					if err != nil {
						logrus.Errorln(err)
						return nil, err
					}
					defer coursor.Close()

					// Extract Authors from the query result
					var authors []Author
					for coursor.HasMore() {
						var author Author
						meta, _ := coursor.ReadDocument(context.Background(), &author)
						author.ID = meta.Key
						authors = append(authors, author)
					}
					return authors, nil
				},
			},
		},
	},
)
