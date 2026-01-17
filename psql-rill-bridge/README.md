# PSQL-Rill Bridge 🔌

A PostgreSQL wire protocol bridge for Rill data. This allows you to query Rill's DuckDB-based data using any PostgreSQL-compatible client (psql, DBeaver, Python's psycopg2, etc.).

## How It Works

```
┌─────────────────┐     PostgreSQL Wire     ┌────────────────────┐     DuckDB     ┌─────────────┐
│  PostgreSQL     │ ◀───────────────────▶   │  psql-rill-bridge  │ ◀──────────▶   │  Rill Data  │
│  Client (psql)  │      Protocol           │  (Go Server)       │    Driver      │  (DuckDB)   │
└─────────────────┘                          └────────────────────┘                └─────────────┘
```

## Prerequisites

- Go 1.21+
- A Rill project with data (run `rill start` first to create the DuckDB database)

## Building

```bash
cd psql-rill-bridge
go mod tidy
go build -o psql-rill-bridge .
```

## Usage

### Step 1: Start a Rill project to generate data

First, you need to run a Rill project to generate the DuckDB database:

```bash
cd ../rill-examples/rill-github-analytics
rill start
```

Wait for the data to load, then stop the server (Ctrl+C). The database will be at `stage.db`.

### Step 2: Start the PostgreSQL bridge

```bash
./psql-rill-bridge -db ../rill-examples/rill-github-analytics/stage.db
```

Options:
- `-db` - Path to DuckDB database file (required)
- `-listen` - Address to listen on (default: `127.0.0.1:5432`)

### Step 3: Connect with a PostgreSQL client

```bash
# Using psql
psql -h localhost -p 5432 -U rill

# Using any password (authentication is open)
```

## Example Queries

Once connected, you can run SQL queries:

```sql
-- Show all tables
SHOW TABLES;

-- Query commits data
SELECT * FROM rill_commits_model LIMIT 10;

-- Get top contributors
SELECT username, COUNT(*) as commits
FROM rill_commits_model
GROUP BY username
ORDER BY commits DESC
LIMIT 20;

-- Analyze file changes
SELECT first_directory, SUM(additions) as additions, SUM(deletions) as deletions
FROM rill_commits_model
GROUP BY first_directory
ORDER BY additions DESC;

-- Time-based analysis
SELECT DATE_TRUNC('month', date) as month, COUNT(*) as commits
FROM rill_commits_model
GROUP BY month
ORDER BY month DESC;
```

## Supported Features

- ✅ Simple queries (SELECT, SHOW TABLES, DESCRIBE)
- ✅ Aggregations (GROUP BY, COUNT, SUM, AVG)
- ✅ Filtering (WHERE, HAVING)
- ✅ Ordering (ORDER BY)
- ✅ Limits (LIMIT, OFFSET)
- ✅ Joins
- ✅ Clear text authentication (accepts any credentials)

## Architecture

This bridge uses:

- **[psql-wire](https://github.com/jeroenrinzema/psql-wire)** - PostgreSQL wire protocol implementation
- **[duckdb-go](https://github.com/duckdb/duckdb-go)** - DuckDB Go driver
- **Rill's data model** - DuckDB databases created by Rill projects

## Connecting from Python

```python
import psycopg2

conn = psycopg2.connect(
    host="localhost",
    port=5432,
    user="rill",
    password="any"
)

cursor = conn.cursor()
cursor.execute("SELECT * FROM rill_commits_model LIMIT 5")
rows = cursor.fetchall()

for row in rows:
    print(row)

conn.close()
```

## Connecting from Other Tools

### DBeaver
1. Create new PostgreSQL connection
2. Host: `localhost`, Port: `5432`
3. Username: `rill`, Password: anything
4. Connect and browse tables

### DataGrip
1. Add new PostgreSQL data source
2. Configure connection to `localhost:5432`
3. Use any credentials

## Limitations

- Read-only access (INSERT/UPDATE/DELETE are passed through but DuckDB may reject them depending on the database mode)
- No prepared statement caching
- Basic authentication (accepts all credentials)
- Single connection mode

## License

Apache 2.0 (same as Rill)
