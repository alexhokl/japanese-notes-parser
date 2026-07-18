# AGENTS.md

Go CLI (cobra + gorm/sqlite) that parses a markdown notes file (JLPT vocabulary
notes format, see README.md) into a SQLite database.

## Commands (via Taskfile.yml — use `task <name>`, not raw go commands, when possible)

- `task build` — `go build -o /dev/null` (compile check only, no binary kept)
- `task test` — `go test ./...`
- `task coverage` — `go test --cover ./...`
- `task coverage-html` — generates `coverage.out`/`coverage.html`
- `task install` — `go install`
- Single package/test: plain `go test ./parser/... -run TestName -v` works fine (no Taskfile shortcut for this)

## Structure

- `cmd/` — cobra commands (`create`, `list`, `list_parts_of_speech`, `root`). `create` is the main entrypoint: reads a notes file line-by-line, calls `ParseLine` (in `cmd/create.go`, not in `parser/`) which delegates regex matching to `parser/`, then upserts `database.Entry` rows via gorm.
- `parser/` — pure regex/string parsing helpers (`regex.go` has all the compiled regexes, `parser.go` has `ParseHeaderLine`/`ParseEnglish`). `parser_test.go` is the main test coverage for parsing logic.
- `database/` — gorm model (`Entry`) and `Upsert`/`AutoMigrate` operations against SQLite via `gorm.io/driver/sqlite`.

## Notes format being parsed

Markdown notes with `## <level>` (JLPT level, e.g. N5) headers, `### <partOfSpeech>` headers,
`- <japanese> - <english>` points, and `  * <japanese> - <english>` sub-points (sub-points link
back to the previous point's kanji via `LinkedKanji`). Level/part-of-speech state is cached
across lines while parsing (see `cachedLevel`/`cachedPartOfSpeech` in `cmd/create.go`).

## Gotchas

- `go.mod` has a commented-out `replace github.com/alexhokl/helper => ../helper` — uncomment
  locally only if developing against an unreleased version of the `helper` module.
- `English` and `Labels` fields on `database.Entry` are stored as `jsonb` via a custom
  `StringArray` type, not normalized tables.
