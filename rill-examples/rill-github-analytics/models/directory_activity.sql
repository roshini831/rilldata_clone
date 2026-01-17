-- Directory Activity Heatmap Model
-- Shows which parts of the codebase are most actively developed
-- @materialize: true

SELECT 
    first_directory as directory,
    COUNT(*) as changes,
    SUM(additions) as additions,
    SUM(deletions) as deletions,
    SUM(additions) - SUM(deletions) as net_changes,
    COUNT(DISTINCT username) as contributors,
    COUNT(DISTINCT commit_hash) as unique_commits,
    ROUND(SUM(additions) * 100.0 / (SELECT SUM(additions) FROM rill_commits_model), 2) as pct_of_project
FROM rill_commits_model
WHERE first_directory IS NOT NULL
GROUP BY first_directory
ORDER BY changes DESC
