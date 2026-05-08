package connection

import (
	"database/sql"
	"fmt"
	"os"
	_ "github.com/lib/pq"
)

var (
	DBConnections *sql.DB 
)

func Initiator() {
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	panic(err)
	// }

	psqlInfo := fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`,
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGDATABASE"),
	)
	var err error
	DBConnections, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(fmt.Sprintf("Error opening database: %v", err))
	}


	// Check connection
	if err = DBConnections.Ping(); err != nil {
		panic(fmt.Sprintf("Database not reachable: %v", err))
	}

	fmt.Println("Successfully connected to database")
}