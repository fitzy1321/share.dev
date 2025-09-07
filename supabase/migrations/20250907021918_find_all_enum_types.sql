-- Find all enum types
SELECT n.nspname AS enum_schema,
       t.typname AS enum_name,
       string_agg(e.enumlabel, ', ' ORDER BY e.enumsortorder) AS enum_values
FROM pg_type t
JOIN pg_enum e ON t.oid = e.enumtypid
JOIN pg_catalog.pg_namespace n ON n.oid = t.typnamespace
GROUP BY n.nspname, t.typname
ORDER BY n.nspname, t.typname;  
