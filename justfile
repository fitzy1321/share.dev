default:
    just --list

db_rest:
    supbase db reset

full_reset: db_reset
    supabase stop
    rm -rf supabae/.temp
    supabase start

start:
    supbase start

stop:
    supabase stop

templ:
    go tool templ generate
