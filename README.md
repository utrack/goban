# goban

`goban` is a Go(lang) linter that bans usage of user-supplied list of functions.

# Usage
`goban -cfg goban.cfg ./...`

# Config
Config is a newline-delimited list of banned symbols. Comments start with pound 
symbol (#).

Examples:
```conf
# bans method `url.Query()` on type *net/url.URL
(*net/url.URL).Query

# bans `context.TODO()`
context.TODO
```

If symbol has a comment on the same line - then it is printed along with the
report.

`fmt.Errorf # use pkg/errors instead` yields `/path/to/file/foo.go:145:15: fmt.Errorf is banned - use pkg/errors instead`

# TODO
- [ ] Ban variables as well
- [ ] Support wildcards for rules
