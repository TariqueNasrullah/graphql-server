package database

import (
	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/sirupsen/logrus"
)

var (
	//Db Database connection
	Db driver.Database

	// BookCollection database connection to book collection
	BookCollection driver.Collection

	// AuthorCollecction database connection to author collection
	AuthorCollecction driver.Collection
)

// InitDB initiates Database connection and Runs
// Necessary Migrations such as Creating Author and
// Book Collection if not exists in the database
func InitDB() {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"http://localhost:8529"},
	})
	handle(err)

	// Create client with db configs
	c, err := driver.NewClient(
		driver.ClientConfig{
			Connection:     conn,
			Authentication: driver.BasicAuthentication("root", "admin"),
		},
	)
	handle(err)

	// Connect to db with pre-created client
	// ArangoDB default database is _system
	Db, err = c.Database(nil, "_system")
	handle(err)

	// Book List migrations - create if not exists
	bookExists, err := Db.CollectionExists(nil, "Book")
	handle(err)
	if !bookExists {
		_, err = Db.CreateCollection(nil, "Book", nil)
		handle(err)
		logrus.Infoln("Book Collection Migrations Successfull.")
	}

	// Author List migrations - create if not exists
	authorExists, err := Db.CollectionExists(nil, "Author")
	handle(err)
	if !authorExists {
		_, err = Db.CreateCollection(nil, "Author", nil)
		handle(err)
		logrus.Infoln("")
	}

	// Database connecton to Book Collection
	BookCollection, err = Db.Collection(nil, "Book")
	handle(err)

	// Databse connection to Author Collection
	AuthorCollecction, err = Db.Collection(nil, "Author")
	handle(err)

	logrus.Println("Db connection successfull")
}

func handle(err error) {
	if err != nil {
		logrus.Fatalln(err)
	}
}
