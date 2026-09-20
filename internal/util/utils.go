package util

import (
	"context"
	"fmt"
	"os"
	db "url_shortener/internal/database/sqlc"

	"github.com/jackc/pgx/v5"
)

var conn *pgx.Conn
var queries *db.Queries

func ConnectToDB() {
	var err error
	conn, err = pgx.Connect(context.Background(), os.Getenv("DB_SOURCE"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unalbe to connect to database: %v\n", err)
		os.Exit(1)
	}
	queries = db.New(conn)
}

func GetQueries() *db.Queries {
	return queries
}