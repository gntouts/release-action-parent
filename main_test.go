package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestProcessArgs(t *testing.T) {
	// Save original version for restoration
	originalVersion := version
	defer func() { version = originalVersion }()

	// Set test version
	version = "1.2.3-abc123"

	tests := []struct {
		name         string
		args         []string
		wantOutput   string
		wantExit     bool
		wantExitCode int
	}{
		{
			name:         "version flag",
			args:         []string{"--version"},
			wantOutput:   "1.2.3-abc123",
			wantExit:     true,
			wantExitCode: 0,
		},
		{
			name:         "empty args",
			args:         []string{},
			wantOutput:   "",
			wantExit:     false,
			wantExitCode: 0,
		},
		{
			name:         "single argument",
			args:         []string{"hello"},
			wantOutput:   "hello",
			wantExit:     false,
			wantExitCode: 0,
		},
		{
			name:         "multiple arguments",
			args:         []string{"hello", "world", "test"},
			wantOutput:   "hello world test",
			wantExit:     false,
			wantExitCode: 0,
		},
		{
			name:         "arguments with spaces",
			args:         []string{"hello world", "test"},
			wantOutput:   "hello world test",
			wantExit:     false,
			wantExitCode: 0,
		},
		{
			name:         "version with other args",
			args:         []string{"--version", "extra"},
			wantOutput:   "--version extra",
			wantExit:     false,
			wantExitCode: 0,
		},
		{
			name:         "special characters",
			args:         []string{"hello!", "@#$%", "test"},
			wantOutput:   "hello! @#$% test",
			wantExit:     false,
			wantExitCode: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, shouldExit, exitCode := processArgs(tt.args)

			if output != tt.wantOutput {
				t.Errorf("processArgs() output = %q, want %q", output, tt.wantOutput)
			}
			if shouldExit != tt.wantExit {
				t.Errorf("processArgs() shouldExit = %v, want %v", shouldExit, tt.wantExit)
			}
			if exitCode != tt.wantExitCode {
				t.Errorf("processArgs() exitCode = %v, want %v", exitCode, tt.wantExitCode)
			}
		})
	}
}

func TestEchoOutput(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple string",
			input: "hello",
			want:  "hello\n",
		},
		{
			name:  "empty string",
			input: "",
			want:  "\n",
		},
		{
			name:  "string with spaces",
			input: "hello world test",
			want:  "hello world test\n",
		},
		{
			name:  "string with special characters",
			input: "hello! @#$% test",
			want:  "hello! @#$% test\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Call function
			echoOutput(tt.input)

			// Restore stdout
			w.Close()
			os.Stdout = old

			// Read captured output
			var buf bytes.Buffer
			io.Copy(&buf, r)
			got := buf.String()

			if got != tt.want {
				t.Errorf("echoOutput() = %q, want %q", got, tt.want)
			}
		})
	}
}

// Benchmark tests
func BenchmarkProcessArgs(b *testing.B) {
	args := []string{"hello", "world", "this", "is", "a", "test"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processArgs(args)
	}
}

func BenchmarkProcessArgsVersion(b *testing.B) {
	version = "1.0.0-test"
	args := []string{"--version"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processArgs(args)
	}
}
