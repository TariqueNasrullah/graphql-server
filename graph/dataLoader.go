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

// GetBooksBatchFn data loader - expect list of AuthorIDs
// Query database with the AuthorIDs and returns List of Books
// belogs to the AuthorIDs
func GetBooksBatchFn(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	handleError := handleErrorFunc

	// Extract authorIds from keys into list of string
	var authorIds []string
	for _, key := range keys {
		authorIds = append(authorIds, key.String())
	}

	// Construct query
	query := fmt.Sprintf(`for book in Book
		let cnt = (for author in book.authors
		filter author.id in %s
		return author)
		filter count(cnt) != 0
		return book`, prepareKeyListQueryFilterString(authorIds))

	// Execute Query
	bindVars := map[string]interface{}{}
	coursor, err := database.Db.Query(context.Background(), query, bindVars)
	if err != nil {
		logrus.Errorln(err)
		handleError(err)
	}
	defer coursor.Close()

	// Each author might have multiple books, booksOfAuthor
	// is a map of key='authorId' value=[]Books
	var booksOfAuthor = make(map[string][]Book)
	for _, ids := range authorIds {
		booksOfAuthor[ids] = []Book{}
	}

	for coursor.HasMore() {
		var book Book
		meta, _ := coursor.ReadDocument(ctx, &book)
		book.ID = meta.Key

		// Insert books to the booksOfAuthor map
		for _, author := range book.Authors {
			booksOfAuthor[author.ID] = append(booksOfAuthor[author.ID], book)
		}
	}

	// Create result in exact order of keys
	var results []*dataloader.Result
	for _, ids := range authorIds {
		result := dataloader.Result{
			Data:  booksOfAuthor[ids],
			Error: nil,
		}
		results = append(results, &result)
	}

	logrus.Infof("[GetBooksBachFn] batch size %d\n", len(results))
	return results
}

// GetAuthorsBatchFn batch function
func GetAuthorsBatchFn(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	handleError := handleErrorFunc

	// Extract authorIDs from keys
	var authorIds []string
	for _, key := range keys {
		authorIds = append(authorIds, key.String())
	}

	// Author slice with length as same of the number of keys/authorIDs
	var authors = make([]Author, len(authorIds))

	// Execute query - Read list of author documents
	metaSlice, errorSlice, err := database.AuthorCollecction.ReadDocuments(ctx, authorIds, authors)
	if err != nil {
		logrus.Errorln(err)
		return handleError(err)
	}

	// Populate result slice
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

	logrus.Infof("[GetAuthorsBachFn] batch size %d\n", len(results))
	return results
}

// GetAuthorByAuthorNameBatchFn batch function of lazy loading author data from database
func GetAuthorByAuthorNameBatchFn(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	handleError := handleErrorFunc

	var authorName []string
	for _, key := range keys {
		authorName = append(authorName, key.String())
	}

	// Construct query
	query := fmt.Sprintf("for author in Author filter author.name in %s return author", prepareKeyListQueryFilterString(authorName))

	// Execute query
	bindVars := map[string]interface{}{}
	coursor, err := database.Db.Query(context.Background(), query, bindVars)
	if err != nil {
		handleError(err)
	}
	defer coursor.Close()

	var results []*dataloader.Result

	// Extract Authors from the query result -
	// Populate result slice
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

	logrus.Infof("[GetAuthorByAuthorNameBatchFn] batch size %d\n", len(results))
	return results
}

// prepareKeyListQueryFilterString converts a list of array into follwing format
// ['id_1', 'id_2', ......]
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
