CREATE TABLE commits (
    id INTEGER PRIMARY KEY,
    author_name VARCHAR,
    commit_message VARCHAR,
    commit_date DATE,
    additions INTEGER,
    deletions INTEGER,
    file_path VARCHAR
);

INSERT INTO commits VALUES
    (1, 'Alice', 'Initial commit', '2024-01-01', 100, 0, 'main.go'),
    (2, 'Bob', 'Add README', '2024-01-02', 50, 5, 'README.md'),
    (3, 'Alice', 'Fix bug in parser', '2024-01-03', 25, 10, 'parser/parser.go'),
    (4, 'Charlie', 'Add tests', '2024-01-04', 200, 20, 'tests/test_main.go'),
    (5, 'Bob', 'Update dependencies', '2024-01-05', 30, 45, 'go.mod'),
    (6, 'Alice', 'Refactor handlers', '2024-01-06', 80, 60, 'handlers/api.go'),
    (7, 'Diana', 'Add documentation', '2024-01-07', 150, 10, 'docs/README.md'),
    (8, 'Charlie', 'Performance optimization', '2024-01-08', 40, 25, 'core/engine.go'),
    (9, 'Bob', 'Fix security issue', '2024-01-09', 15, 8, 'auth/auth.go'),
    (10, 'Alice', 'Add new feature', '2024-01-10', 120, 30, 'features/new.go');

CREATE TABLE projects (
    id INTEGER PRIMARY KEY,
    name VARCHAR,
    description VARCHAR,
    stars INTEGER,
    forks INTEGER
);

INSERT INTO projects VALUES
    (1, 'rill', 'BI-as-code platform', 5000, 350),
    (2, 'duckdb', 'In-process SQL OLAP database', 15000, 1200),
    (3, 'psql-wire', 'PostgreSQL wire protocol', 500, 50);
