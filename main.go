package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	gq "github.com/graphql-go/graphql"
)

type Todo struct {
	ID   int    `json:"id"`
	Task string `json:"task"`
}

var todos = []Todo{
	{ID: 1, Task: "Learn Go"},
	{ID: 2, Task: "Learn Chi"},
}

var todoType = gq.NewObject(
	gq.ObjectConfig{
		Name: "Todo",
		Fields: gq.Fields{
			"id": &gq.Field{
				Type: gq.Int,
			},
			"task": &gq.Field{
				Type: gq.String,
			},
		},
	},
)

var queryType = gq.NewObject(
	gq.ObjectConfig{
		Name: "Query",
		Fields: gq.Fields{
			"todos": &gq.Field{
				Type:        gq.NewList(todoType),
				Description: "get all todos",
				Resolve: func(p gq.ResolveParams) (interface{}, error) {
					return todos, nil
				},
			},
		},
	},
)

var schema, _ = gq.NewSchema(
	gq.SchemaConfig{
		Query: queryType,
	},
)

func execQuery(query string, schema gq.Schema) *gq.Result {
	result := gq.Do(gq.Params{
		Schema:        schema,
		RequestString: query,
	})

	if len(result.Errors) > 0 {
		fmt.Printf("errors: %v", result.Errors)
	}

	return result
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/gq", func(w http.ResponseWriter, r *http.Request) {
		result := execQuery(r.URL.Query().Get("query"), schema)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	http.ListenAndServe(":3000", r)
}
