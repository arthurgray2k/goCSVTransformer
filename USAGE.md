# Usage

`gocsvtransformer` offers a fast, robust CSV manipulation toolkit accessible both as a CLI application and a reusable Go library.

---

## CLI

### General Flags
Most commands support the following general flags:
- `-output <file>`: Output directly to a file (default: stdout).
- `-delimiter <char>`: Custom delimiter for CSV parsing (default: `,`).
- `-header=false`: Treat the first row as data rather than a header.
- `-stats`: Print throughput, bytes read/written, and duration at the end.
- `-verbose`: Enable debug logging via `slog`.

### Converting CSV to NDJSON

Convert a CSV file into a Newline Delimited JSON stream.

**Command:**
```bash
gocsvtransformer csv-to-ndjson [flags] [file]
```

**Examples:**
```bash
# Output to file, view statistics
gocsvtransformer csv-to-ndjson -output data.ndjson -stats data.csv
```

### Converting NDJSON to CSV

Convert a Newline Delimited JSON stream into a CSV file.

**Command:**
```bash
gocsvtransformer ndjson-to-csv [flags] [file]
```

### Converting CSV to JSON

Convert a CSV file into a JSON array of objects.

**Command:**
```bash
gocsvtransformer csv-to-json [flags] [file]
```

**Examples:**
```bash
gocsvtransformer csv-to-json -pretty -indent 4 -output data.json data.csv
```

### Converting JSON to CSV

Convert a JSON array of objects into a CSV file.

**Command:**
```bash
gocsvtransformer json-to-csv [flags] [file]
```

### Validation

Validate the integrity of a CSV file.

**Command:**
```bash
gocsvtransformer validate [flags] [file]
```

---

## Go Library

`gocsvtransformer` is built around streaming abstractions, avoiding large memory allocations even for massive files.

### Converting CSV to NDJSON Stream

```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/arthurgray2k/gocsvtransformer/pkg/gocsvtransformer"
)

func main() {
	f, _ := os.Open("data.csv")
	defer f.Close()

	out, _ := os.Create("data.ndjson")
	defer out.Close()

	opts := gocsvtransformer.DefaultOptions()
	transformer := gocsvtransformer.New(opts)
	
	if err := transformer.CSVToNDJSON(context.Background(), f, out); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
```
