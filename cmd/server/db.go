package main

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
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
	dbDir := "./data"
	dbPath := fmt.Sprintf("%s/v1.sqlite3", dbDir)

	// The leading zero is important as it forces the number to be parsed as octal.
	if err := os.MkdirAll(dbDir, os.FileMode(0700)); err != nil {
		log.Fatal().Err(err).Msg("Failed to create database directory")
	}

	db, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return db, nil
}
