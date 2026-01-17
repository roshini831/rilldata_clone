-- Monthly Trends Model
-- Shows development activity changes over time by month
-- @materialize: true

SELECT 
    DATE_TRUNC('month', date) as month,
    YEAR(date) as year,
    MONTH(date) as month_num,
    COUNT(*) as file_changes,
    COUNT(DISTINCT commit_hash) as commits,
    COUNT(DISTINCT username) as active_contributors,
    SUM(additions) as additions,
    SUM(deletions) as deletions,
    SUM(additions) - SUM(deletions) as net_growth
FROM rill_commits_model
WHERE date >= DATE '2020-01-01'
GROUP BY DATE_TRUNC('month', date), YEAR(date), MONTH(date)
ORDER BY month DESC
