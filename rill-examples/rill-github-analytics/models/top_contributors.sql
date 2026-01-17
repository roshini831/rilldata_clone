-- Top Contributors Model
-- Shows the most active developers ranked by total code changes
-- @materialize: true

SELECT 
    username,
    COUNT(*) as total_changes,
    COUNT(DISTINCT commit_hash) as unique_commits,
    SUM(additions) as lines_added,
    SUM(deletions) as lines_deleted,
    SUM(additions) - SUM(deletions) as net_lines,
    ROUND(SUM(additions) * 100.0 / NULLIF(SUM(additions) + SUM(deletions), 0), 1) as add_percentage,
    MIN(date) as first_commit,
    MAX(date) as last_commit
FROM rill_commits_model
GROUP BY username
ORDER BY total_changes DESC
