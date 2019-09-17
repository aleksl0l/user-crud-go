package store

import (
	"database/sql"
	"fmt"
	config "github.com/spf13/viper"
)

type DB struct {
	*sql.DB
}

func NewDB() (*DB, error) {
	dataSourceName := fmt.Sprintf(
		"postgres://%s:%s@%s:%s",
		config.GetString(`database.user`),
		config.GetString(`database.password`),
		config.GetString(`database.host`),
		config.GetString(`database.port`),
	)
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
