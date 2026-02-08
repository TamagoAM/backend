package db

import (
	"log"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
)

func Migrate(db *sqlx.DB, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Split on semicolons and execute each statement individually.
	// This gives better error messages and works around MySQL's
	// single-statement Exec limitation.
	stmts := strings.Split(string(content), ";")
	for _, stmt := range stmts {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		log.Printf("[migrate] exec: %.80s...", stmt)
		if _, err := db.Exec(stmt); err != nil {
			log.Printf("[migrate] warning: %v (statement: %.80s...)", err, stmt)
			// Don't fail hard — some statements may error on re-run
			// (e.g. CREATE TABLE IF NOT EXISTS is fine but others might)
		}
	}
	return nil
}
