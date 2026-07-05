package cli_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/arthurgray2k/gocsvtransformer/internal/cli"
)

func TestApp_Run_HelpAndVersion(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected int
		outSub   string
	}{
		{"no args", []string{}, 1, "Usage:"},
		{"help command", []string{"help"}, 0, "Usage:"},
		{"version command", []string{"version"}, 0, "gocsvtransformer v"},
		{"unknown command", []string{"unknown"}, 1, "Unknown command: unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var outBuf bytes.Buffer
			var errBuf bytes.Buffer
			
			app := &cli.App{
				Stdin:  strings.NewReader(""),
				Stdout: &outBuf,
				Stderr: &errBuf,
			}
			
			exitCode := app.Run(context.Background(), tt.args)
			if exitCode != tt.expected {
				t.Errorf("expected exit code %d, got %d", tt.expected, exitCode)
			}
			
			combined := outBuf.String() + errBuf.String()
			if !strings.Contains(combined, tt.outSub) {
				t.Errorf("expected output to contain %q, got: %s", tt.outSub, combined)
			}
		})
	}
}

func TestApp_Run_AllCommands(t *testing.T) {
	tests := []struct {
		command string
		input   string
		expect  int
	}{
		{"validate", "id,name\n1,alice\n", 0},
		{"csv-to-json", "id,name\n1,alice\n", 0},
		{"csv-to-ndjson", "id,name\n1,alice\n", 0},
		{"json-to-csv", "[{\"id\":\"1\"}]", 0},
		{"ndjson-to-csv", "{\"id\":\"1\"}\n", 0},
	}

	for _, tt := range tests {
		t.Run(tt.command, func(t *testing.T) {
			var outBuf bytes.Buffer
			var errBuf bytes.Buffer
			
			app := cli.NewApp()
			app.Stdin = strings.NewReader(tt.input)
			app.Stdout = &outBuf
			app.Stderr = &errBuf
			
			exitCode := app.Run(context.Background(), []string{tt.command, "-stats"}) // test stats execution too
			if exitCode != tt.expect {
				t.Errorf("expected exit code %d, got %d. err: %s", tt.expect, exitCode, errBuf.String())
			}
		})
	}
}
