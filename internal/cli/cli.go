package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"

	internalio "github.com/arthurgray2k/gocsvtransformer/internal/io"
	"github.com/arthurgray2k/gocsvtransformer/pkg/gocsvtransformer"
)

const (
	version = "0.1.0"
)

// App encapsulates the CLI execution.
type App struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// NewApp creates a new CLI App with standard IO streams.
func NewApp() *App {
	return &App{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

// Run executes the CLI.
func (a *App) Run(ctx context.Context, args []string) int {
	if len(args) < 1 {
		a.usage()
		return 1
	}

	command := args[0]
	cmdArgs := args[1:]

	switch command {
	case "validate":
		return a.runValidate(ctx, cmdArgs)
	case "csv-to-json":
		return a.runCSVToJSON(ctx, cmdArgs)
	case "json-to-csv":
		return a.runJSONToCSV(ctx, cmdArgs)
	case "csv-to-ndjson":
		return a.runCSVToNDJSON(ctx, cmdArgs)
	case "ndjson-to-csv":
		return a.runNDJSONToCSV(ctx, cmdArgs)
	case "version":
		fmt.Fprintf(a.Stdout, "gocsvtransformer v%s\n", version)
		return 0
	case "help":
		a.usage()
		return 0
	default:
		fmt.Fprintf(a.Stderr, "Unknown command: %s\n", command)
		a.usage()
		return 1
	}
}

func (a *App) usage() {
	fmt.Fprintf(a.Stderr, `gocsvtransformer - A fast, scalable CSV transformation utility.

Usage:
  gocsvtransformer <command> [arguments]

Commands:
  validate        Validate a CSV file
  csv-to-json     Convert CSV to JSON
  json-to-csv     Convert JSON to CSV
  csv-to-ndjson   Convert CSV to NDJSON
  ndjson-to-csv   Convert NDJSON to CSV
  version         Print version information
  help            Show this help message

Use "gocsvtransformer <command> -h" for more information about a command.
`)
}

type commandConfig struct {
	opts   *gocsvtransformer.Options
	input  io.Reader
	output io.Writer
	stats  bool
}

func (a *App) parseCommonFlags(fs *flag.FlagSet, args []string) (*commandConfig, error) {
	delimiterStr := fs.String("delimiter", ",", "Delimiter character")
	headerRow := fs.Bool("header", true, "Indicate if the first row is a header")
	pretty := fs.Bool("pretty", false, "Pretty print output (for JSON)")
	indent := fs.Int("indent", 2, "Indentation level for pretty print")
	outputFile := fs.String("output", "", "Output file (default: stdout)")
	verbose := fs.Bool("verbose", false, "Enable verbose/debug logging")
	stats := fs.Bool("stats", false, "Print processing statistics upon completion")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	// Setup logging
	logLevel := slog.LevelInfo
	if *verbose {
		logLevel = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(a.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	var input io.Reader = a.Stdin
	if fs.NArg() > 0 {
		file, err := os.Open(fs.Arg(0))
		if err != nil {
			slog.Error("Error opening input file", "error", err)
			return nil, err
		}
		input = file
	}

	var output io.Writer = a.Stdout
	if *outputFile != "" {
		file, err := os.Create(*outputFile)
		if err != nil {
			slog.Error("Error creating output file", "error", err)
			return nil, err
		}
		output = file
	}

	delims := []rune(*delimiterStr)
	if len(delims) != 1 {
		err := fmt.Errorf("invalid delimiter: must be a single character")
		slog.Error(err.Error())
		return nil, err
	}

	opts := gocsvtransformer.Options{
		Delimiter: delims[0],
		HeaderRow: *headerRow,
		Pretty:    *pretty,
		Indent:    *indent,
	}

	return &commandConfig{
		opts:   &opts,
		input:  input,
		output: output,
		stats:  *stats,
	}, nil
}

func (a *App) printStats(tr *internalio.TrackingReader, tw *internalio.TrackingWriter) {
	elapsed := tr.Elapsed()
	bytesRead := tr.BytesRead()
	bytesWritten := tw.BytesWritten()

	// Avoid division by zero
	sec := elapsed.Seconds()
	if sec <= 0 {
		sec = 1
	}

	throughput := float64(bytesRead) / 1024 / 1024 / sec

	slog.Info("Processing complete",
		slog.Int64("bytes_read", bytesRead),
		slog.Int64("bytes_written", bytesWritten),
		slog.String("elapsed", elapsed.String()),
		slog.String("throughput", fmt.Sprintf("%.2f MB/s", throughput)),
	)
}

func (a *App) runValidate(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	fs.SetOutput(a.Stderr)

	cfg, err := a.parseCommonFlags(fs, args)
	if err != nil {
		return 1
	}

	slog.Debug("Starting validation")

	tr := internalio.NewTrackingReader(cfg.input)
	transformer := gocsvtransformer.New(*cfg.opts)

	if err := transformer.ValidateCSV(ctx, tr); err != nil {
		slog.Error("Validation failed", "error", err)
		return 1
	}

	fmt.Fprintf(a.Stdout, "CSV is valid.\n")

	if cfg.stats {
		a.printStats(tr, internalio.NewTrackingWriter(io.Discard))
	}
	return 0
}

func (a *App) runCSVToJSON(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("csv-to-json", flag.ContinueOnError)
	fs.SetOutput(a.Stderr)

	cfg, err := a.parseCommonFlags(fs, args)
	if err != nil {
		return 1
	}

	slog.Debug("Starting CSV to JSON conversion")

	tr := internalio.NewTrackingReader(cfg.input)
	tw := internalio.NewTrackingWriter(cfg.output)
	transformer := gocsvtransformer.New(*cfg.opts)

	if err := transformer.CSVToJSON(ctx, tr, tw); err != nil {
		slog.Error("Conversion failed", "error", err)
		return 1
	}

	if cfg.stats {
		a.printStats(tr, tw)
	}
	return 0
}

func (a *App) runJSONToCSV(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("json-to-csv", flag.ContinueOnError)
	fs.SetOutput(a.Stderr)

	cfg, err := a.parseCommonFlags(fs, args)
	if err != nil {
		return 1
	}

	slog.Debug("Starting JSON to CSV conversion")

	tr := internalio.NewTrackingReader(cfg.input)
	tw := internalio.NewTrackingWriter(cfg.output)
	transformer := gocsvtransformer.New(*cfg.opts)

	if err := transformer.JSONToCSV(ctx, tr, tw); err != nil {
		slog.Error("Conversion failed", "error", err)
		return 1
	}

	if cfg.stats {
		a.printStats(tr, tw)
	}
	return 0
}

func (a *App) runCSVToNDJSON(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("csv-to-ndjson", flag.ContinueOnError)
	fs.SetOutput(a.Stderr)

	cfg, err := a.parseCommonFlags(fs, args)
	if err != nil {
		return 1
	}

	slog.Debug("Starting CSV to NDJSON conversion")

	tr := internalio.NewTrackingReader(cfg.input)
	tw := internalio.NewTrackingWriter(cfg.output)
	transformer := gocsvtransformer.New(*cfg.opts)

	if err := transformer.CSVToNDJSON(ctx, tr, tw); err != nil {
		slog.Error("Conversion failed", "error", err)
		return 1
	}

	if cfg.stats {
		a.printStats(tr, tw)
	}
	return 0
}

func (a *App) runNDJSONToCSV(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("ndjson-to-csv", flag.ContinueOnError)
	fs.SetOutput(a.Stderr)

	cfg, err := a.parseCommonFlags(fs, args)
	if err != nil {
		return 1
	}

	slog.Debug("Starting NDJSON to CSV conversion")

	tr := internalio.NewTrackingReader(cfg.input)
	tw := internalio.NewTrackingWriter(cfg.output)
	transformer := gocsvtransformer.New(*cfg.opts)

	if err := transformer.NDJSONToCSV(ctx, tr, tw); err != nil {
		slog.Error("Conversion failed", "error", err)
		return 1
	}

	if cfg.stats {
		a.printStats(tr, tw)
	}
	return 0
}
