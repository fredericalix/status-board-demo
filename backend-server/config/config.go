package config

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

// ConnectToDB établit une connexion à la base de données PostgreSQL.
func ConnectToDB() (*sql.DB, error) {
	dbURI := os.Getenv("POSTGRESQL_ADDON_URI")
	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		return nil, err
	}
	return db, nil
}
