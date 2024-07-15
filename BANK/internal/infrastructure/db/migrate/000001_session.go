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
	_, err := tx.Exec(`
        CREATE TABLE IF NOT EXISTS bank (
            bank_uuid varchar(255) NOT NULL PRIMARY KEY,
            user_uuid varchar(255) NOT NULL,
            amount int NOT NULL,
            status int,
        );
    `)
	if err != nil {
		return fmt.Errorf("could not create payment_session table: %v", err)
	}

	return nil
}

func downCreateTables(tx *sql.Tx) error {
	// Drop tables
	_, err := tx.Exec(`DROP TABLE IF EXISTS bank;`)
	if err != nil {
		return fmt.Errorf("could not drop user_balance table: %v", err)
	}

	return nil
}
