# goban

`goban` is a Go(lang) linter that bans usage of user-supplied list of functions.

# Usage
`goban -cfg goban.cfg ./...`

# Config
Config is a newline-delimited list of banned symbols. Comment lines with pound
symbol (#).
Examples:

```conf
(*net/url.URL).Query # bans method `url.Query()` on type *net/url.URL
context.TODO # bans `context.TODO()`
```

# TODO
- [ ] Ban variables as well
- [ ] Support wildcards for rules
