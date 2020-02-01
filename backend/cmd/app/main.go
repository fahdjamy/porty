package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/friendsofgo/graphiql"
	"github.com/graphql-go/graphql"
)

//Dummy struct
type Dummy struct {
	ID             int      `json:"id"`
	Position       string   `json:"position"`
	Company        string   `json:"company"`
	Description    string   `json:"description"`
	SkillsRequired []string `json:"skillsRequired"`
	Location       string   `json:"location"`
	EmploymentType string   `json:"employmentType"`
}

var dummyType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Dummy",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"position": &graphql.Field{
				Type: graphql.String,
			},
			"company": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"location": &graphql.Field{
				Type: graphql.String,
			},
			"employmentType": &graphql.Field{
				Type: graphql.String,
			},
			"skillsRequired": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
		},
	},
)

type reqBody struct {
	Query string `json:"query"`
}

func main() {
	graphiqlHandler, err := graphiql.NewGraphiqlHandler("/graphql")
	if err != nil {
		panic(err)
	}

	http.Handle("/graphql", gqlHandler())
	http.Handle("/graphiql", graphiqlHandler)
	log.Fatal(http.ListenAndServe(":7000", nil))
}

func gqlHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			http.Error(w, "No query data", 400)
			return
		}

		var rBody reqBody
		err := json.NewDecoder(r.Body).Decode(&rBody)
		if err != nil {
			http.Error(w, "Error parsing JSON request body", 400)
		}
		fmt.Fprintf(w, "%s", processQuery(rBody.Query))
	})
}

func processQuery(query string) (result string) {
	retrieveDummyData := retrieveDummyDataFromFile()
	params := graphql.Params{Schema: gqlSchema(retrieveDummyData), RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		fmt.Printf("failed to execute graphql operation, errors: %+v", r.Errors)
	}
	rJSON, _ := json.Marshal(r)

	return fmt.Sprintf("%s", rJSON)
}

//Open the file data.json and retrieve json data
func retrieveDummyDataFromFile() func() []Dummy {
	return func() []Dummy {
		jsonf, err := os.Open("data.json")
		if err != nil {
			fmt.Printf("failed to open json file, error: %v", err)
		}
		jsonDataFromFile, _ := ioutil.ReadAll(jsonf)
		defer jsonf.Close()

		var dummyData []Dummy

		err = json.Unmarshal(jsonDataFromFile, dummyData)
		if err != nil {
			fmt.Printf("failed to parse json, error: %v", err)
		}

		return dummyData
	}
}

// Define the GraphQL Schema
func gqlSchema(queryData func() []Dummy) graphql.Schema {
	fields := graphql.Fields{
		"allData": &graphql.Field{
			Type:        graphql.NewList(dummyType),
			Description: "All data",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return queryData(), nil
			},
		},
		"data": &graphql.Field{
			Type:        dummyType,
			Description: "Get Jobs by ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, success := params.Args["id"].(int)
				if success {
					for _, job := range queryData() {
						if int(job.ID) == id {
							return job, nil
						}
					}
				}
				return nil, nil
			},
		},
	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		fmt.Printf("failed to create new schema, error: %v", err)
	}
	return schema
}
