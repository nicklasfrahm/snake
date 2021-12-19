package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const schema = `
CREATE TABLE IF NOT EXISTS 'queues' (
	'id' VARCHAR(36) PRIMARY KEY,
	'name' VARCHAR(128) NOT NULL,
	'owner' VARCHAR(128) NOT NULL,
	'title' VARCHAR(128) NOT NULL,
	'description' VARCHAR(2048) NOT NULL,
	'number' INTEGER NOT NULL,
	UNIQUE(name),
	UNIQUE(owner)
);
`

type Queue struct {
	ID          string `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Owner       string `db:"owner" json:"owner"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
	Number      int    `db:"number" json:"number"`
	Token       string `json:"token,omitempty"`
}

func MigrateDB() (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", "./database.sqlite3")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return db, nil
}
