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
	// Create payment_session table
	_, err := tx.Exec(`
        CREATE TABLE IF NOT EXISTS session_in_api (
            bank_id VARCHAR(255) NOT NULL,
            user_id VARCHAR(255) NOT NULL,
            user_uuid VARCHAR(255) NOT NULL,
            session_id VARCHAR(255) NOT NULL,
            canceled BOOLEAN NOT NULL
        );
    `)
	if err != nil {
		return fmt.Errorf("could not create payment_session table: %v", err)
	}

	// Create user_balance table
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS session_in_bank (
		    user_uuid VARCHAR(255) NOT NULL,
		    bank_id VARCHAR(255) NOT NULL,
		    amount INT NOT NULL,
		    canceled BOOLEAN NOT NULL,
		    returned BOOLEAN,
		);
	`)
	if err != nil {
		return fmt.Errorf("could not create user_balance table: %v", err)
	}
	return nil
}

func downCreateTables(tx *sql.Tx) error {

	// Drop tables
	_, err := tx.Exec(`DROP TABLE IF EXISTS session_in_api;`)
	if err != nil {
		return fmt.Errorf("could not drop user_balance table: %v", err)
	}

	_, err = tx.Exec(`DROP TABLE IF EXISTS session_in_bank;`)
	if err != nil {
		return fmt.Errorf("could not drop payment_session table: %v", err)
	}

	return nil
}
