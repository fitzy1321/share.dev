-- Return bool, enum exists
create or replace function enum_exists(enum_name text, enum_schema text default 'public')
returns boolean as $$
  select exists (
    select 1
    from pg_type t
    join pg_namespace n on t.typnamespace = n.oid
        where t.typname = enum_name
          and n.nspname = enum_schema
    );
$$ language sql set search_path = '';

-- return bool, table exists
create or replace function table_exists(table_name text, schema_name text default 'public')
returns boolean as $$
  select exists (
    select 1
    from pg_tables
    where schemaname = schema_name
      and tablename = table_name
  );
$$ language sql set search_path = '';
