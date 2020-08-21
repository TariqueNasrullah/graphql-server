package graph

import (
	"context"
	"fmt"

	"github.com/TariqueNasrullah/graphql-server/database"
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
)

var qureyType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "GetBooksOfAuthor",
		Fields: graphql.Fields{
			"book": &graphql.Field{
				Type: graphql.NewList(listBookType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					authorName := p.Context.Value("authorName").(string)
					query := fmt.Sprintf(`for author in Author
						filter author.name=='%s'
						for book in Book
						for authid in book.authors
						filter authid.id==author._key
						return book`, authorName)

					bindVars := map[string]interface{}{}
					coursor, err := database.Db.Query(context.Background(), query, bindVars)
					if err != nil {
						logrus.Errorln(err)
						return nil, err
					}
					defer coursor.Close()

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
			"books": &graphql.Field{
				Type: graphql.NewList(bookType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					query := "for book in Book return book"
					bindVars := map[string]interface{}{}
					coursor, err := database.Db.Query(context.Background(), query, bindVars)
					if err != nil {
						logrus.Errorln(err)
						return nil, err
					}
					defer coursor.Close()

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
			"authors": &graphql.Field{
				Type: graphql.NewList(authorType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					query := "for author in Author return author"
					bindVars := map[string]interface{}{}
					coursor, err := database.Db.Query(context.Background(), query, bindVars)
					if err != nil {
						logrus.Errorln(err)
						return nil, err
					}
					defer coursor.Close()

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
