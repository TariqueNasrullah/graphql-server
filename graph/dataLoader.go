package graph

import (
	"context"
	"fmt"

	"github.com/TariqueNasrullah/graphql-server/database"
	"github.com/graph-gophers/dataloader"
	"github.com/sirupsen/logrus"
)

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

// GetAuthorBatchFn batch function of lazy loading author data from database
func GetAuthorBatchFn(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {

	handleError := func(err error) []*dataloader.Result {
		var results []*dataloader.Result
		var result dataloader.Result
		result.Error = err
		results = append(results, &result)
		return results
	}

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
