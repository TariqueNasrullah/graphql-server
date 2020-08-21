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

type body struct {
	Query string `json:"query"`
}

func main() {
	database.InitDB()
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		var p body
		json.NewDecoder(r.Body).Decode(&p)
		w.Header().Set("content-type", "application/json")

		// result := graph.ExecuteQuery(r.URL.Query().Get("query"))

		var loaders = make(map[string]*dataloader.Loader, 1)

		loaders["GetAuthorByAuthorName"] = dataloader.NewBatchedLoader(graph.GetAuthorByAuthorNameBatchFn)
		loaders["GetAuthors"] = dataloader.NewBatchedLoader(graph.GetAuthorsBatchFn)
		loaders["GetBooks"] = dataloader.NewBatchedLoader(graph.GetBooksBatchFn)

		ctx := context.WithValue(context.Background(), "loaders", loaders)

		result := graph.ExecuteQuery(ctx, p.Query)
		json.NewEncoder(w).Encode(result)
	})

	fmt.Println("Now server is running on port 8080")
	fmt.Println("Test with Get      : curl -g 'http://localhost:8080/graphql?query={Book(Id:\"1\"){Title}}'")
	http.ListenAndServe(":8080", nil)
}
