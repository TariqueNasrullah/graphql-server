package graph

import (
	"context"
	"fmt"

	"github.com/TariqueNasrullah/graphql-server/database"
	"github.com/graph-gophers/dataloader"
	"github.com/sirupsen/logrus"
)

var handleErrorFunc = func(err error) []*dataloader.Result {
	var results []*dataloader.Result
	var result dataloader.Result
	result.Error = err
	results = append(results, &result)
	return results
}

// GetBooksBatchFn data loader
func GetBooksBatchFn(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	handleError := handleErrorFunc

	var authorIds []string
	for _, key := range keys {
		authorIds = append(authorIds, key.String())
	}

	query := fmt.Sprintf(`for book in Book
		for authorIds in book.authors
		filter authorIds.id in %s
		return book`, prepareKeyListQueryFilterString(authorIds))
	bindVars := map[string]interface{}{}
	coursor, err := database.Db.Query(context.Background(), query, bindVars)
	if err != nil {
		logrus.Errorln(err)
		handleError(err)
	}
	defer coursor.Close()

	var booksOfAuthor = make(map[string][]Book)

	for _, ids := range authorIds {
		booksOfAuthor[ids] = []Book{}
	}

	for coursor.HasMore() {
		var book Book
		meta, _ := coursor.ReadDocument(ctx, &book)
		book.ID = meta.Key

		for _, author := range book.Authors {
			booksOfAuthor[author.ID] = append(booksOfAuthor[author.ID], book)
		}
		logrus.Infoln(book)
	}
	// for ids := range booksOfAuthor {
	// 	logrus.Infof("authorId-> %v   : %v", ids, booksOfAuthor[ids])
	// }

	var results []*dataloader.Result
	logrus.Infof("[GetBooksBachFn] batch size %d\n", len(results))
	return results
}

// GetAuthorsBatchFn batch function
func GetAuthorsBatchFn(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	handleError := handleErrorFunc

	var authorIds []string
	for _, key := range keys {
		authorIds = append(authorIds, key.String())
	}

	var authors = make([]Author, len(authorIds))

	metaSlice, errorSlice, err := database.AuthorCollecction.ReadDocuments(ctx, authorIds, authors)
	if err != nil {
		logrus.Errorln(err)
		return handleError(err)
	}

	var results []*dataloader.Result

	for idx, author := range authors {
		if errorSlice[idx] != nil {
			return handleError(errorSlice[idx])
		}
		author.ID = metaSlice[idx].Key
		result := dataloader.Result{
			Data:  author,
			Error: nil,
		}

		results = append(results, &result)
	}

	logrus.Infof("[GetAuthorBachFn] batch size %d\n", len(results))
	return results
}

// GetAuthorByAuthorNameBatchFn batch function of lazy loading author data from database
func GetAuthorByAuthorNameBatchFn(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	handleError := handleErrorFunc

	var authorName []string
	for _, key := range keys {
		authorName = append(authorName, key.String())
	}

	query := fmt.Sprintf("for author in Author filter author.name in %s return author", prepareKeyListQueryFilterString(authorName))
	bindVars := map[string]interface{}{}
	coursor, err := database.Db.Query(context.Background(), query, bindVars)
	if err != nil {
		handleError(err)
	}
	defer coursor.Close()

	var results []*dataloader.Result

	for coursor.HasMore() {
		var author Author
		meta, _ := coursor.ReadDocument(ctx, &author)
		author.ID = meta.Key

		result := dataloader.Result{
			Data:  author,
			Error: nil,
		}
		results = append(results, &result)
	}

	logrus.Infof("[GetAuthorBachFn] batch size %d\n", len(results))
	return results
}

func prepareKeyListQueryFilterString(keys []string) string {
	str := "["
	for idx, key := range keys {
		str += fmt.Sprintf("'%s'", key)
		if idx < len(keys)-1 {
			str += ","
		}
	}
	str += "]"
	return str
}
