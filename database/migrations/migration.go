package migration

import (
	"database/sql"
	"embed"
	"fmt"
	migrate "github.com/rubenv/sql-migrate"
)

//go:embed sql_migrations/*.sql
var dbMigrations embed.FS

func Initiator(dbParam *sql.DB) error {
	migrations := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: dbMigrations,
		Root:       "sql_migrations",
	}

	n, err := migrate.Exec(dbParam, "postgres", migrations, migrate.Up)
	if err != nil {
		return err
	}

	fmt.Println("Migration success, applied", n, "migrations!")
	return nil
}