package graph

import (
	"context"
	"fmt"
	"regexp"

	"github.com/graphql-go/graphql"
)

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: qureyType,
	},
)

// ExecuteQuery executes the query
func ExecuteQuery(ctx context.Context, query string) *graphql.Result {
	re := regexp.MustCompile(`\(\s*name:\s*"(.*?)"\)`)
	match := re.FindStringSubmatch(query)
	if len(match) > 0 {
		ctx = context.WithValue(ctx, "authorName", match[1])
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		VariableValues: map[string]interface{}{
			"Id": "1234",
		},
		Context: ctx,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, errors: %v", result.Errors)
	}
	return result
}
