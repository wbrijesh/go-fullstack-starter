package database

import "database/sql"

func InitialiseDatabase(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS users (
      id TEXT PRIMARY KEY,
      email TEXT NOT NULL,
      hashed_password TEXT NOT NULL
    );
  `
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
