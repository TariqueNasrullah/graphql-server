package graph

import (
	"context"
	"regexp"

	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
)

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    qureyType,
		Mutation: mutationType,
	},
)

// ExecuteQuery executes the query
func ExecuteQuery(ctx context.Context, query string) *graphql.Result {
	// Extract author's name from query string and pass it down with the context
	re := regexp.MustCompile(`\(\s*name:\s*"(.*?)"\)`)
	match := re.FindStringSubmatch(query)
	if len(match) > 0 {
		// set context
		ctx = context.WithValue(ctx, "authorName", match[1])
	}

	// GraphQL server call
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		Context:       ctx,
	})
	if len(result.Errors) > 0 {
		logrus.Errorf("wrong result, errors: %v\n", result.Errors)
	}
	return result
}
