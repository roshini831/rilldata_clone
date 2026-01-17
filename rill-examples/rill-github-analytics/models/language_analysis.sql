-- Language Analysis Model
-- Shows file extension distribution across the codebase
-- @materialize: true

SELECT 
    CASE 
        WHEN file_extension = '.go' THEN 'Go'
        WHEN file_extension = '.ts' THEN 'TypeScript'
        WHEN file_extension = '.svelte' THEN 'Svelte'
        WHEN file_extension = '.sql' THEN 'SQL'
        WHEN file_extension = '.yaml' OR file_extension = '.yml' THEN 'YAML'
        WHEN file_extension = '.json' THEN 'JSON'
        WHEN file_extension = '.md' THEN 'Markdown'
        WHEN file_extension = '.proto' THEN 'Protobuf'
        WHEN file_extension = '.css' THEN 'CSS'
        WHEN file_extension = '.html' THEN 'HTML'
        ELSE 'Other'
    END as language,
    file_extension,
    COUNT(*) as file_changes,
    SUM(additions) as lines_added,
    SUM(deletions) as lines_deleted,
    COUNT(DISTINCT username) as contributors
FROM rill_commits_model
WHERE file_extension IS NOT NULL AND file_extension != ''
GROUP BY file_extension
ORDER BY file_changes DESC
