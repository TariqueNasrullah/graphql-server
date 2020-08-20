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

// InitDB initiates db
func InitDB() {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"http://localhost:8529"},
	})
	handle(err)
	c, err := driver.NewClient(
		driver.ClientConfig{
			Connection:     conn,
			Authentication: driver.BasicAuthentication("root", "admin"),
		},
	)
	handle(err)

	Db, err = c.Database(nil, "_system")
	handle(err)

	bookExists, err := Db.CollectionExists(nil, "Book")
	handle(err)
	if !bookExists {
		_, err = Db.CreateCollection(nil, "Book", nil)
		handle(err)
	}

	BookCollection, err = Db.Collection(nil, "Book")
	handle(err)
	AuthorCollecction, err = Db.Collection(nil, "Author")
	handle(err)

	logrus.Println("Db connection successfull")
}

func handle(err error) {
	if err != nil {
		logrus.Fatalln(err)
	}
}
