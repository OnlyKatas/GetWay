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
        CREATE TABLE IF NOT EXISTS payment_session (
            id SERIAL primary key,
            bank_id VARCHAR(255) UNIQUE, 
            user_id INT NOT NULL,
            user_uuid VARCHAR(255) NOT NULL,
            amount INT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            status INT NOT NULL,
            returnded BOOLEAN NOT NULL DEFAULT FALSE,
        );
    `)
	if err != nil {
		return fmt.Errorf("could not create payment_session table: %v", err)
	}

	// Create user_balance table
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS user_balance (
		    user_id INT NOT NULL,
		    user_uuid VARCHAR(255) NOT NULL,
		    balance INT CHECK (balance >= 0),
		    hold_balance INT,
		    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return fmt.Errorf("could not create user_balance table: %v", err)
	}

	// Create function to update updated_at column
	_, err = tx.Exec(`
		CREATE OR REPLACE FUNCTION update_updated_at_column()
		RETURNS TRIGGER AS $$
		BEGIN
		   NEW.updated_at = NOW();
		   RETURN NEW;
		END;
		$$ language 'plpgsql';
	`)
	if err != nil {
		return fmt.Errorf("could not create function update_updated_at_column: %v", err)
	}

	// Create triggers for updating updated_at column
	_, err = tx.Exec(`
		CREATE TRIGGER update_user_balance_updated_at
		BEFORE UPDATE ON user_balance
		FOR EACH ROW
		EXECUTE PROCEDURE update_updated_at_column();

		CREATE TRIGGER update_payment_session_updated_at
		BEFORE UPDATE ON payment_session
		FOR EACH ROW
		EXECUTE PROCEDURE update_updated_at_column();
	`)
	if err != nil {
		return fmt.Errorf("could not create triggers: %v", err)
	}

	return nil
}

func downCreateTables(tx *sql.Tx) error {
	// Drop triggers first to avoid dependency issues
	_, err := tx.Exec(`
		DROP TRIGGER IF EXISTS update_user_balance_updated_at ON user_balance;
		DROP TRIGGER IF EXISTS update_payment_session_updated_at ON payment_session;
	`)
	if err != nil {
		return fmt.Errorf("could not drop triggers: %v", err)
	}

	// Drop function
	_, err = tx.Exec(`DROP FUNCTION IF EXISTS update_updated_at_column;`)
	if err != nil {
		return fmt.Errorf("could not drop function update_updated_at_column: %v", err)
	}

	// Drop tables
	_, err = tx.Exec(`DROP TABLE IF EXISTS user_balance;`)
	if err != nil {
		return fmt.Errorf("could not drop user_balance table: %v", err)
	}

	_, err = tx.Exec(`DROP TABLE IF EXISTS payment_session;`)
	if err != nil {
		return fmt.Errorf("could not drop payment_session table: %v", err)
	}

	return nil
}
