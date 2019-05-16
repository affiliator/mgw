package storage

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
)

var (
	connection  gorm.DB
	initialized = false
)

func Connection() *gorm.DB {
	if initialized == false {
		err := initialize()

		if err != nil {
			panic(err)
		}
	}

	return &connection
}

func Close() {
	if &connection == nil {
		return
	}

	err := Connection().Close()
	if err != nil {
		fmt.Printf("failed to close storage connection: %s", err)
	}
}

func initialize() error {
	d, dsn := getConfigured()
	if d == "" || dsn == "" {
		return errors.New("unable to create new storage instance")
	}

	db, err := gorm.Open(d, dsn)
	if err != nil {
		return err
	}

	connection = *db
	initialized = true

	return nil
}

func getConfigured() (dialect string, dsn string) {
	d, err := GetDialects().Configured()
	if err != nil {
		// Todo: Log & maybe SQLite Fallback?
		panic(err)
	}

	return d.Name, d.DSN()
}
