package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TariqueNasrullah/graphql-server/database"
	"github.com/TariqueNasrullah/graphql-server/graph"
	"github.com/graph-gophers/dataloader"
)

func main() {
	// Initate database connection with migration
	database.InitDB()

	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")

		// Create DataLoader context and pass it down to the GraphQl Shcema
		var loaders = make(map[string]*dataloader.Loader, 1)

		loaders["GetAuthorByAuthorName"] = dataloader.NewBatchedLoader(graph.GetAuthorByAuthorNameBatchFn)
		loaders["GetAuthors"] = dataloader.NewBatchedLoader(graph.GetAuthorsBatchFn)
		loaders["GetBooks"] = dataloader.NewBatchedLoader(graph.GetBooksBatchFn)

		ctx := context.WithValue(context.Background(), "loaders", loaders)

		// Execute Query in GraphQL
		result := graph.ExecuteQuery(ctx, r.URL.Query().Get("query"))
		json.NewEncoder(w).Encode(result)
	})

	// Run  http Server
	fmt.Println("Server is running on port 8080")
	fmt.Println("Test with Get      : curl -g 'http://localhost:8080/graphql?query={Book(Id:\"1\"){Title}}'")
	http.ListenAndServe(":8080", nil)
}
