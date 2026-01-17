package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	// Connect to the psql-rill-bridge
	connStr := "host=localhost port=15432 user=rill password=any sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer db.Close()

	fmt.Println("✓ Connected to psql-rill-bridge!")
	fmt.Println()

	// Determine query
	query := "SHOW TABLES"
	if len(os.Args) > 1 {
		query = os.Args[1]
	}

	fmt.Printf("Query: %s\n", query)
	fmt.Println("─────────────────────────────────────────")

	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		log.Fatalf("Failed to get columns: %v", err)
	}

	// Print header
	for i, col := range columns {
		if i > 0 {
			fmt.Print(" | ")
		}
		fmt.Print(col)
	}
	fmt.Println()
	fmt.Println("─────────────────────────────────────────")

	// Prepare scan destination
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// Print rows
	rowCount := 0
	for rows.Next() {
		err := rows.Scan(valuePtrs...)
		if err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}

		for i, v := range values {
			if i > 0 {
				fmt.Print(" | ")
			}
			if v == nil {
				fmt.Print("NULL")
			} else {
				fmt.Printf("%v", v)
			}
		}
		fmt.Println()
		rowCount++
	}

	fmt.Println("─────────────────────────────────────────")
	fmt.Printf("(%d rows)\n", rowCount)
}
