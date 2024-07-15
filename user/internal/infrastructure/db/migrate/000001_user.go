package migrate

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTables, downCreateTables)
}

func upCreateTables(tx *sql.Tx) error {
	// Create users table
	_, err := tx.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            uuid VARCHAR(255) NOT NULL UNIQUE,
            first_name VARCHAR(255) NOT NULL,
            last_name VARCHAR(255) NOT NULL,
            username VARCHAR(255) NOT NULL UNIQUE,
            email VARCHAR(255) NOT NULL UNIQUE,
            password VARCHAR(255) NOT NULL
        );
    `)
	if err != nil {
		return fmt.Errorf("could not create users table: %v", err)
	}

	return nil
}

func downCreateTables(tx *sql.Tx) error {
	// Drop tables in reverse order of creation to avoid foreign key constraints
	_, err := tx.Exec(`DROP TABLE IF EXISTS users;`)
	if err != nil {
		return fmt.Errorf("could not drop users table: %v", err)
	}
	return nil
}
