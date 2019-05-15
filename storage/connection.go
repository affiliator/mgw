package storage

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

var connection gorm.DB

func Connection() *gorm.DB {
	if &connection == nil {
		_ = initialize()
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
	return nil
}

func getConfigured() (dialect string, dsn string) {
	d, err := GetDialects().Configured()
	if err != nil {
		// Todo: Implement logging
		return "", ""
	}

	return d.Name(), d.DSN()
}
