# Share.dev

A nonsense Social Media clone.

Built using Golang + HTMX!

## Full Tech Stack

- Golang v1.24.6
- Echo API Framework
- [a-h/Templ](https://github.com/a-h/templ) package
- [HTMX](https://github.com/bigskysoftware/htmx)
- [Bulma CSS](https://github.com/jgthms/bulma)
- Supabase - [supabase-go](https://github.com/supabase-community/supabase-go) package
- [Air (hot reloading)](https://github.com/air-verse/air)
- Fly.io hosting

## Current Scope

- [x] Landing Page
- [x] Login / Signup Auth Actions
  - [ ] Supabase Auth
  - [ ] Email Confirmation
- [ ] Scalable User Sessions
- [ ] "Feed" Page
  - [ ] Post Text
  - [ ] See Subscribed User posts
- [ ] Subscribe to user
  - [ ] Unsubscribe
- [ ] Might need individual User page, Follow users
- [ ] Supabase DB for posts?
- [ ] Supabase storage for images, videos, etc?

## Dev Notes

Using Air for live reloading.  
Install air: `go install github.com/air-verse/air@latest`.  
_Note_: GOPATH (i.e. where go install will place your package, usually `$Home/go/bin`) must be in your $PATH  

Supabase cli: `brew install supabase/tap/supabase`.  
Fly.io cli: `brew install flyctl`

Run app locally with supabase

```sh
supabase start
air # live reload go files, templates, and html
```

## References

- Primeagon talk on golang and htmx: <https://www.youtube.com/watch?v=x7v6SNIgJpE>
- Supabase-go package: <https://github.com/supabase-community/supabase-go>
- Local Development with Supabase <https://supabase.com/docs/guides/local-development>
- Fly.io for hosting and deployments <https://fly.io/docs/>
