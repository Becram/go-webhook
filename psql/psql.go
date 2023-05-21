package psql

import (
	"log"
	"os"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type Notification struct {
	Id        int64
	Arthur    string
	Title     string
	Timestamp time.Time
	Version   string
	Body      string
}

// The function connects to a PostgreSQL database and creates a schema.
func Connect() *pg.DB {
	psqlConn := &pg.Options{
		User:     os.Getenv("POSTGRESQL_USERNAME"),
		Password: os.Getenv("POSTGRESQL_PASSWORD"),
		Addr:     os.Getenv("POSTGRESQL_ADDRESS"),
		Database: os.Getenv("POSTGRESQL_DATABASE"),
	}
	var db *pg.DB = pg.Connect(psqlConn)

	if db == nil {
		log.Fatal("failed to connect to database")
	}
	log.Println("successfully connected to database")

	err := createSchema(db)
	CheckNilErr(err)
	return db
}

// The function creates database tables for specified models in Go using the pg library.
func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*Notification)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			log.Fatal(err)
		}

	}
	return nil
}

// The function checks if an error is nil and logs it as fatal if it is not.
func CheckNilErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// This function creates a new notification in a database using the provided notification object.
func CreateNotification(notification Notification) {
	db := Connect()
	defer db.Close()

	_, err := db.Model(&notification).Insert()

	CheckNilErr(err)
}
