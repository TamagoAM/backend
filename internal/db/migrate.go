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

	// Split SQL into statements respecting quoted strings.
	// A semicolon only counts as a statement terminator when it is
	// outside single-quoted ('…') literals — this avoids breaking
	// on semicolons that appear inside VARCHAR/TEXT values.
	stmts := splitStatements(string(content))
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

// splitStatements splits raw SQL on semicolons that are outside of
// single-quoted string literals.  Escaped quotes (”) inside literals
// are handled correctly.
func splitStatements(sql string) []string {
	var stmts []string
	var buf strings.Builder
	inQuote := false
	for i := 0; i < len(sql); i++ {
		ch := sql[i]
		if ch == '\'' {
			if inQuote {
				// Check for escaped quote ('')
				if i+1 < len(sql) && sql[i+1] == '\'' {
					buf.WriteByte(ch)
					buf.WriteByte(ch)
					i++ // skip the second quote
					continue
				}
				inQuote = false
			} else {
				inQuote = true
			}
			buf.WriteByte(ch)
		} else if ch == ';' && !inQuote {
			stmts = append(stmts, buf.String())
			buf.Reset()
		} else {
			buf.WriteByte(ch)
		}
	}
	// Remaining text after the last semicolon (if any)
	if trailing := buf.String(); strings.TrimSpace(trailing) != "" {
		stmts = append(stmts, trailing)
	}
	return stmts
}
