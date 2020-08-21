package graph

import (
	"fmt"
	"strings"

	"github.com/graph-gophers/dataloader"
	"github.com/graphql-go/graphql"
)

var getBooksOfAuthorResolver = func(p graphql.ResolveParams) (interface{}, error) {
	var (
		v       = p.Context.Value
		loaders = v("loaders").(map[string]*dataloader.Loader)
	)

	authorName := p.Args["name"].(string)

	thunk := loaders["GetAuthorByAuthorName"].Load(p.Context, dataloader.StringKey(authorName))

	return func() (interface{}, error) {
		author, err := thunk()
		if err != nil {
			return nil, err
		}
		return author, nil
	}, nil
}

var getBooksResolver = func(p graphql.ResolveParams) (interface{}, error) {
	var (
		sourceAuthor = p.Source.(Author)
		v            = p.Context.Value
		loaders      = v("loaders").(map[string]*dataloader.Loader)
		// handleErrors = func(errors []error) error {
		// 	var errs []string
		// 	for _, e := range errors {
		// 		errs = append(errs, e.Error())
		// 	}
		// 	return fmt.Errorf(strings.Join(errs, "\n"))
		// }
	)

	thunk := loaders["GetBooks"].Load(p.Context, dataloader.StringKey(sourceAuthor.ID))

	return func() (interface{}, error) {
		books, err := thunk()
		if err != nil {
			return nil, err
		}
		return books, nil
	}, nil
}

var getAuthorsResolver = func(p graphql.ResolveParams) (interface{}, error) {
	var (
		sourceBook   = p.Source.(Book)
		v            = p.Context.Value
		loaders      = v("loaders").(map[string]*dataloader.Loader)
		handleErrors = func(errors []error) error {
			var errs []string
			for _, e := range errors {
				errs = append(errs, e.Error())
			}
			return fmt.Errorf(strings.Join(errs, "\n"))
		}
	)

	var authorIds []string

	for _, author := range sourceBook.Authors {
		authorIds = append(authorIds, author.ID)
	}

	var keys = dataloader.NewKeysFromStrings(authorIds)

	thunk := loaders["GetAuthors"].LoadMany(p.Context, keys)

	return func() (interface{}, error) {
		author, err := thunk()
		if err != nil {
			return nil, handleErrors(err)
		}
		return author, nil
	}, nil
}

var author = Author{
	ID:     "1",
	Name:   "tarique",
	ISBNNo: "1234",
}

var book = Book{
	ID:          "1",
	Title:       "A book of happiness",
	Description: "An incredible book",
}
