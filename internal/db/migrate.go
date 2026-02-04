package db

import (
	"os"

	"github.com/jmoiron/sqlx"
)

func Migrate(db *sqlx.DB, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	_, err = db.Exec(string(content))
	return err
}
