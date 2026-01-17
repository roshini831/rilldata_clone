package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/marcboeker/go-duckdb"
	wire "github.com/jeroenrinzema/psql-wire"
)

var (
	listenAddr string
	dbPath     string
)

func init() {
	flag.StringVar(&listenAddr, "listen", "127.0.0.1:5432", "Address to listen on")
	flag.StringVar(&dbPath, "db", "", "Path to DuckDB database file (required)")
}

var db *sql.DB
var handler *QueryHandler

func main() {
	flag.Parse()

	if dbPath == "" {
		fmt.Println("Usage: psql-rill-bridge -db /path/to/database.db")
		fmt.Println("")
		fmt.Println("Example with Rill GitHub Analytics:")
		fmt.Println("  psql-rill-bridge -db ./rill-github-analytics/stage.db")
		fmt.Println("")
		fmt.Println("Then connect with:")
		fmt.Println("  psql -h localhost -p 5432 -U rill")
		os.Exit(1)
	}

	// Open DuckDB database
	var err error
	db, err = sql.Open("duckdb", dbPath)
	if err != nil {
		log.Fatalf("Failed to open DuckDB database: %v", err)
	}
	defer db.Close()

	// Verify connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping DuckDB database: %v", err)
	}

	log.Printf("Connected to DuckDB database: %s", dbPath)

	// Create the handler
	handler = NewQueryHandler(db)

	log.Printf("PostgreSQL wire protocol server listening on %s", listenAddr)
	log.Println("Connect with: psql -h localhost -p 5432 -U rill")
	log.Println("")
	log.Println("Example queries:")
	log.Println("  SHOW TABLES;")
	log.Println("  SELECT * FROM <table_name> LIMIT 10;")

	// Start the server using the simple API
	if err := wire.ListenAndServe(listenAddr, queryHandler); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// queryHandler is the top-level handler function for the wire protocol
func queryHandler(ctx context.Context, query string) (wire.PreparedStatements, error) {
	return handler.HandleQuery(ctx, query)
}
