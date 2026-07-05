# gocsvtransformer

`gocsvtransformer` is a fast, scalable, and extensible CSV transformation utility and reusable Go library. 

It aims to provide both an idiomatic Go API and a convenient CLI tool for converting CSV files to/from structured formats like JSON and NDJSON, as well as providing strong validation features.

## Status

Currently in **Stage 4** (Production Ready):
- Go module initialized.
- Library API defined.
- Validation, CSV ↔ JSON, CSV ↔ NDJSON implemented.
- Robust Streaming I/O utilizing `io.Reader`/`io.Writer`.
- Performance optimizations using `sync.Pool` for zero-allocation mapping on massive data streams.
- Full statistics tracking (`--stats`) and Structured Logging (`--verbose`).

## Installation

```bash
go install github.com/arthurgray2k/gocsvtransformer/cmd/gocsvtransformer@latest
```

## Features

- **Stream Large Datasets**: Processing is done using Go's `io.Reader` and `io.Writer` streams.
- **CSV Validation**: Validates duplicate headers, empty headers, and row sizing before data is processed.
- **CSV ↔ JSON**: Convert between CSV files and JSON arrays.
- **CSV ↔ NDJSON**: Convert between CSV files and NDJSON streams.
- **High Performance**: Employs `sync.Pool` to avoid Garbage Collection lag during heavy row iteration.
- **Statistics**: Includes throughput, byte-counting, and elapsed time tracking.
- **Context Aware**: Full support for cancellation via `context.Context`.

## CLI Usage

See [USAGE.md](USAGE.md) for detailed examples.

Example:
```bash
gocsvtransformer csv-to-ndjson data.csv --output data.ndjson --stats --verbose
```

## Future Roadmap
- Built-in SQL/Parquet support
- Built-in Excel support

## License
MIT License
