package models

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // postges dialect for database connection
)

// Db is our database struct used for interacting with the database
type Db struct {
	*gorm.DB
}

// New makes a new database using the connection string and
// returns it, otherwise returns the error
func New(dbURIString string) (*Db, error) {
	// connect to the database
	db, err := gorm.Open("postgres", dbURIString)
	if err != nil {
		fmt.Printf("Cannot connect to %s database:\n ", "postgres")
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %v database\n", dbURIString)
	}
	return &Db{db}, nil
}

// ConnString returns a connection string based on the parameters it's given
// This would normally also contain the password, however we're not using one
func ConnString(host string, port string, user string, dbName string, dbPassword string) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		host,
		port,
		user,
		dbName,
		dbPassword,
	)
}

// User shape
type User struct {
	ID         int
	Name       string
	Age        int
	Profession string
	Friendly   bool
}

// GetUsersByName is called within our user query for graphql
func (d *Db) GetUsersByName(name string) []User {
	// Create slice of Users for our response
	users := []User{}
	d.Table("users").Where("Name = ?", name).Find(&users)
	if len(users) <= 0 {
		return []User{}
	}
	return users
}
