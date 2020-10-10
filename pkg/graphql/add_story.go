package main

import (
	"encoding/json"
	"fmt"
	"github.com/graphql-go/graphql"
	"log"
)

type Error struct {
	description string
}

func addStory(query string) *graphql.Result {
	errCfg := graphql.ObjectConfig{
		Name:   "Error",
		Fields: "",
	}

	fields := graphql.Fields{
		"success": &graphql.Field{
			Type: graphql.Boolean,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return true, nil
			},
		},
		"data": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "data", nil
			},
		},
		"error": &graphql.Field{
			Type: graphql.NewObject(errCfg),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return &Error{}, nil
			},
		},
	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	}

	return r
}

func main() {

	// Query
	query := `
		{
			success
			data
			error
		}
	`

	r := addStory(query)

	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJSON)
}
