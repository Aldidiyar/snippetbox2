package sql

import (
	"database/sql"
	"log"
)

type exampleModel struct {
	DB *sql.DB
}

func (m *exampleModel) exampleTransaction() error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO...")
	if err != nil {
		if rb := tx.Rollback(); rb != nil {
			log.Printf("query failed: %v, unable to abort: %v", err, rb)
			return rb
		}
		log.Printf("query failed: %v", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("transaction commit failed: %v", err)
		return err
	}

	return nil
}
