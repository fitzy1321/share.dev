# Share.dev

A nonsense Social Media clone.

Built using Golang + HTMX!

## Full Tech Stack

- Golang
  - Builtin http server
- Templ
- [Air (hot reloading)](https://github.com/air-verse/air)
- HTMX
- Bulma CSS
- Supabase

## Current Scope

- [ ] Decent Homepage
- [ ] Login / Signup pages
- [ ] Scalable User Sessions
- [ ]" Feed" Page
- [ ] Post Text
- [ ] Follow User thingy
- [ ] Might need individual User page, Follow users

Supabase for Auth and DB

## Dev Notes

Using Air for live reloading.  
Install air: `go install github.com/air-verse/air@latest`.  
_Note_: GOPATH (i.e. where go install will place your package, usually `$Home/go/bin`) must be in your $PATH  

Supabase cli: `brew install supabase/tap/supabase`

Run app locally with supabase

```sh
supabase start
air # live reload go files, templates, and html
```

## References

- Primeagon talk on golang and htmx: <https://www.youtube.com/watch?v=x7v6SNIgJpE>
- Supabase-go package: <https://github.com/supabase-community/supabase-go>
- Local Development with Supabase <https://supabase.com/docs/guides/local-development>
