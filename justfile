default:
    just --list

local_db_migrations:
    supbase db reset

start:
    supbase start

stop:
    supabase stop

templ:
    go tool templ generate
