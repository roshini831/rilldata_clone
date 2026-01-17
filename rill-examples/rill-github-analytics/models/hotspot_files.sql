-- Hotspot Files Model
-- Shows the most frequently modified files in the codebase
-- @materialize: true

SELECT 
    file_path,
    filename,
    first_directory,
    file_extension,
    COUNT(*) as modification_count,
    COUNT(DISTINCT username) as unique_authors,
    COUNT(DISTINCT commit_hash) as unique_commits,
    SUM(additions) as total_additions,
    SUM(deletions) as total_deletions,
    SUM(additions) + SUM(deletions) as total_churn,
    MIN(date) as first_modified,
    MAX(date) as last_modified
FROM rill_commits_model
WHERE file_path IS NOT NULL
GROUP BY file_path, filename, first_directory, file_extension
HAVING COUNT(*) > 20
ORDER BY modification_count DESC
