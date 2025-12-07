# gathercode

`gathercode` is a fast, concurrent CLI tool that collects `.go` and `.sql` files from local folders, then merges them into a single structured Markdown file.

Each file is written in this format:

```
-- repo/subfolder/file.go
<file contents>
---------------
```

This is useful for:

- LLM ingestion
- Code audits
- Documentation snapshots
- Offline code reviews

---

## Features

- Recursive folder scanning
- Concurrent collection
- `.go` and `.sql` filtering (customizable)
- Clean Markdown output
- External test suite
- Safe unzip for remote ZIP sources
- Streaming file writes (low memory usage)

---

## Project Structure

```
gathercode/
├── go.mod
├── cmd/
│   └── gather/
│       └── main.go
├── pkg/
│   ├── gather/
│   ├── github/
│   ├── utils/
│   └── writer/
└── tests/
```

---

## Build

```
go build ./cmd/gather
```

Binary output:

```
./gather
```

---

## Usage

### Basic local folder

```
./gather --inputs ./myproject --output out.md
```

### Multiple input folders

```
./gather --inputs ./projA,./projB --output combined.md
```

### Enable concurrency

```
./gather --inputs ./myproject --parallel 8
```

### Custom extensions

```
./gather --inputs ./myproject --ext .go,.sql,.txt
```

---

## CLI Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--inputs` | Comma-separated list of folders | required |
| `--output` | Output markdown file | `aggregated.md` |
| `--ext` | File extensions to include | `.go,.sql` |
| `--parallel` | Number of workers | `4` |
| `--timeout` | Max runtime | `10m` |
| `--include-hidden` | Include hidden files | `false` |

---

## Output Example

```
-- myrepo/cmd/main.go
package main

func main() {}
---------------
-- myrepo/db/schema.sql
CREATE TABLE users (...);
---------------
```

---

## Run Tests

```
go test ./... -v
```

All tests are located in:

```
/tests
```

and run as external integration tests.

---

## Makefile Usage

If you prefer using `make`, the project includes a simple `Makefile`:

```
make build # build the binary
make run # run the CLI via go run
make test # run all tests
make tidy # clean go.mod/go.sum
make fmt # format all Go files
make vet # run go vet
make clean # remove built binary
```

---

## Design Notes

- All file reads use streaming (`io.Copy`) to avoid memory spikes.
- ZIP extraction uses a disk-backed temp file (no in-memory archive load).
- Concurrency uses a bounded worker pool with `context.Context` cancellation.
- Tests validate:
  - Local scanning
  - Concurrent collection
  - ZIP fetch flow
  - Markdown writer
  - Hidden file filtering

---

## License

MIT