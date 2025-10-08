default:
    just --list

db-local-push:
    supabase db push --local

db-reset:
    supbase db reset

db-push:
    supbase db push

full-reset: db-reset
    supabase stop
    rm -rf supabae/.temp
    supabase start

start:
    supabase start

stop:
    supabase stop

fly-deploy:
    flyctl deploy

templ:
    go tool templ generate

update-deps:
    go get -u all
    go mod tidy
