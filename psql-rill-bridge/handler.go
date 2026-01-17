package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	wire "github.com/jeroenrinzema/psql-wire"
)

// QueryHandler handles PostgreSQL queries and routes them to DuckDB
type QueryHandler struct {
	db *sql.DB
}

// NewQueryHandler creates a new query handler
func NewQueryHandler(db *sql.DB) *QueryHandler {
	return &QueryHandler{db: db}
}

// HandleQuery processes incoming SQL queries from PostgreSQL clients
func (h *QueryHandler) HandleQuery(ctx context.Context, query string) (wire.PreparedStatements, error) {
	query = strings.TrimSpace(query)
	log.Printf("Query: %s", query)

	// Handle empty query
	if query == "" {
		return wire.Prepared(wire.NewStatement(func(ctx context.Context, writer wire.DataWriter, parameters []wire.Parameter) error {
			return writer.Complete("OK")
		})), nil
	}

	// Handle special commands
	upperQuery := strings.ToUpper(query)

	// Handle SHOW TABLES (DuckDB compatible)
	if upperQuery == "SHOW TABLES" || upperQuery == "SHOW TABLES;" {
		return h.executeQuery(ctx, "SHOW TABLES")
	}

	// Handle \dt style table listing
	if strings.HasPrefix(upperQuery, "\\DT") || strings.HasPrefix(upperQuery, "\\D") {
		return h.executeQuery(ctx, "SHOW TABLES")
	}

	// Handle DESCRIBE/DESC
	if strings.HasPrefix(upperQuery, "DESCRIBE ") || strings.HasPrefix(upperQuery, "DESC ") {
		parts := strings.Fields(query)
		if len(parts) >= 2 {
			tableName := strings.TrimSuffix(parts[1], ";")
			return h.executeQuery(ctx, fmt.Sprintf("DESCRIBE %s", tableName))
		}
	}

	// Handle regular queries
	return h.executeQuery(ctx, query)
}

// executeQuery runs a query against DuckDB and returns results in PostgreSQL format
func (h *QueryHandler) executeQuery(ctx context.Context, query string) (wire.PreparedStatements, error) {
	// Execute the query to get schema info first
	rows, err := h.db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("Query error: %v", err)
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	// Get column information
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, fmt.Errorf("failed to get column types: %w", err)
	}

	// Build column definitions for PostgreSQL wire protocol
	wireColumns := make(wire.Columns, len(columns))
	for i, col := range columns {
		wireColumns[i] = wire.Column{
			Name:  col,
			Oid:   mapDuckDBTypeToOid(columnTypes[i].DatabaseTypeName()),
			Table: 0,
			Width: 256,
		}
	}

	// Read all rows into memory
	var allRows [][]any
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	for rows.Next() {
		err := rows.Scan(valuePtrs...)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		rowCopy := make([]any, len(values))
		for i, v := range values {
			rowCopy[i] = convertValue(v)
		}
		allRows = append(allRows, rowCopy)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	rowCount := len(allRows)

	// Create statement handler
	handle := func(ctx context.Context, writer wire.DataWriter, parameters []wire.Parameter) error {
		for _, row := range allRows {
			err := writer.Row(row)
			if err != nil {
				return fmt.Errorf("failed to write row: %w", err)
			}
		}
		return writer.Complete(fmt.Sprintf("SELECT %d", rowCount))
	}

	return wire.Prepared(wire.NewStatement(handle, wire.WithColumns(wireColumns))), nil
}

// mapDuckDBTypeToOid maps DuckDB types to PostgreSQL OIDs
func mapDuckDBTypeToOid(duckdbType string) uint32 {
	switch strings.ToUpper(duckdbType) {
	case "INTEGER", "INT", "INT4", "INT32":
		return pgtype.Int4OID
	case "BIGINT", "INT8", "INT64":
		return pgtype.Int8OID
	case "SMALLINT", "INT2", "INT16":
		return pgtype.Int2OID
	case "FLOAT", "FLOAT4", "REAL":
		return pgtype.Float4OID
	case "DOUBLE", "FLOAT8":
		return pgtype.Float8OID
	case "BOOLEAN", "BOOL":
		return pgtype.BoolOID
	case "DATE":
		return pgtype.DateOID
	case "TIME":
		return pgtype.TimeOID
	case "TIMESTAMP", "TIMESTAMP WITH TIME ZONE", "TIMESTAMPTZ":
		return pgtype.TimestampOID
	case "INTERVAL":
		return pgtype.IntervalOID
	case "BLOB", "BYTEA":
		return pgtype.ByteaOID
	case "UUID":
		return pgtype.UUIDOID
	case "JSON":
		return pgtype.JSONOID
	default:
		// Default to text for VARCHAR, TEXT, and unknown types
		return pgtype.TextOID
	}
}

// convertValue converts DuckDB values to PostgreSQL wire format
func convertValue(v interface{}) interface{} {
	if v == nil {
		return nil
	}

	switch val := v.(type) {
	case []byte:
		return string(val)
	case int64, int32, int16, int8, int:
		return val
	case float64, float32:
		return val
	case bool:
		return val
	case string:
		return val
	default:
		// Convert to string representation for unknown types
		return fmt.Sprintf("%v", val)
	}
}
