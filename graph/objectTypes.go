package graph

import (
	"github.com/graph-gophers/dataloader"
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
			"author": &graphql.Field{
				Type: authorType,
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var (
						v       = p.Context.Value
						loaders = v("loaders").(map[string]*dataloader.Loader)
					)

					authorName := p.Args["name"].(string)

					thunk := loaders["GetAuthor"].Load(p.Context, dataloader.StringKey(authorName))

					return func() (interface{}, error) {
						author, err := thunk()
						if err != nil {
							return nil, err
						}
						return author, nil
					}, nil
				},
			},
		},
	},
)
