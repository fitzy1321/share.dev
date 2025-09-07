default:
    just --list

db_local_push:
    supabase db push --local

db_reset:
    supbase db reset

db_push:
    supbase db push

full_reset: db_reset
    supabase stop
    rm -rf supabae/.temp
    supabase start

start:
    supbase start

stop:
    supabase stop
