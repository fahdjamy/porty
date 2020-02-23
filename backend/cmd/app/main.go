package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fahdJamy/porty/src/gql"
	"github.com/fahdJamy/porty/src/models"
	"github.com/fahdJamy/porty/src/server"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/graphql-go/graphql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

func main() {
	router, db := initializer()
	defer db.Close()

	fmt.Println("running on port 7000")
	log.Fatal(http.ListenAndServe(":7000", router))
}

func initializer() (*chi.Mux, *models.Db) {
	// Create a new connection to our pg database
	// Create a new router
	router := chi.NewRouter()

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Print("Failed to load environment variables and trying again:\n", err)
		err = godotenv.Load("../../.env")
	}
	if err != nil {
		log.Print("Failed to load environment variables completely ", err)
	}

	var dbUser = os.Getenv("DB_USER")
	var dbPassword = os.Getenv("DB_PASSWORD")
	var dbName = os.Getenv("DB_NAME")
	var dbPort = os.Getenv("DB_PORT")
	var dbHost = os.Getenv("DB_HOST")

	// Build connection string
	dbURI := models.ConnString(dbHost, dbPort, dbUser, dbName, dbPassword)
	postgresDb, err := models.New(dbURI)

	// connect to the database
	if err != nil {
		fmt.Printf("Cannot connect to %s database:\n ", "postgres")
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %v database\n", dbName)
	}

	// Create our root query for graphql
	rootQuery := gql.NewRoot(postgresDb)
	// Create a new graphql schema, passing in the the root query
	sc, err := graphql.NewSchema(
		graphql.SchemaConfig{Query: rootQuery.Query},
	)
	if err != nil {
		fmt.Println("Error creating schema: ", err)
	}

	// Create a server struct that holds a pointer to our database as well
	// as the address of our graphql schema
	s := server.Server{
		GqlSchema: &sc,
	}

	// Add some middleware to our router
	router.Use(
		render.SetContentType(render.ContentTypeJSON), // set content-type headers as application/json
		middleware.Logger,          // log api request calls
		middleware.DefaultCompress, // compress results, mostly gzipping assets and json
		middleware.StripSlashes,    // match paths with a trailing slash, strip it, and continue routing through the mux
		middleware.Recoverer,       // recover from panics without crashing server
	)

	// Create the graphql route with a Server method to handle it
	router.Post("/graphql", s.GraphQL())

	return router, postgresDb
}
