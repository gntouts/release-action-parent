package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const testBinaryName = "echo_test"

func buildTestBinary(t *testing.T) string {
	t.Helper()

	// Create temporary directory
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, testBinaryName)

	// Build the binary
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")

	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build test binary: %v", err)
	}

	return binaryPath
}

func runBinary(t *testing.T, binaryPath string, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()

	cmd := exec.Command(binaryPath, args...)

	stdoutBytes, err := cmd.Output()
	stdout = string(stdoutBytes)

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			stderr = string(exitError.Stderr)
			exitCode = exitError.ExitCode()
		} else {
			t.Fatalf("Failed to run binary: %v", err)
		}
	}

	return stdout, stderr, exitCode
}

func TestFunctionalEcho(t *testing.T) {
	binaryPath := buildTestBinary(t)

	tests := []struct {
		name           string
		args           []string
		expectedStdout string
		expectedStderr string
		expectedExit   int
	}{
		{
			name:           "no arguments",
			args:           []string{},
			expectedStdout: "\n",
			expectedExit:   0,
		},
		{
			name:           "single argument",
			args:           []string{"hello"},
			expectedStdout: "hello\n",
			expectedExit:   0,
		},
		{
			name:           "multiple arguments",
			args:           []string{"hello", "world"},
			expectedStdout: "hello world\n",
			expectedExit:   0,
		},
		{
			name:           "arguments with spaces",
			args:           []string{"hello world", "test"},
			expectedStdout: "hello world test\n",
			expectedExit:   0,
		},
		{
			name:           "special characters",
			args:           []string{"hello!", "@#$%^&*()", "test"},
			expectedStdout: "hello! @#$%^&*() test\n",
			expectedExit:   0,
		},
		{
			name:           "empty string argument",
			args:           []string{""},
			expectedStdout: "\n",
			expectedExit:   0,
		},
		{
			name:           "mixed arguments",
			args:           []string{"hello", "", "world"},
			expectedStdout: "hello  world\n",
			expectedExit:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, exitCode := runBinary(t, binaryPath, tt.args...)

			if stdout != tt.expectedStdout {
				t.Errorf("stdout = %q, want %q", stdout, tt.expectedStdout)
			}
			if stderr != tt.expectedStderr {
				t.Errorf("stderr = %q, want %q", stderr, tt.expectedStderr)
			}
			if exitCode != tt.expectedExit {
				t.Errorf("exit code = %d, want %d", exitCode, tt.expectedExit)
			}
		})
	}
}

func TestFunctionalVersion(t *testing.T) {
	binaryPath := buildTestBinary(t)

	stdout, stderr, exitCode := runBinary(t, binaryPath, "--version")

	// Check exit code
	if exitCode != 0 {
		t.Errorf("--version exit code = %d, want 0", exitCode)
	}

	// Check stderr is empty
	if stderr != "" {
		t.Errorf("--version stderr = %q, want empty", stderr)
	}

	// Check stdout contains version (should be empty in test build, but should not error)
	stdout = strings.TrimSpace(stdout)
	// Since we're building without version injection, it should be empty
	if stdout != "" {
		t.Logf("Version output: %q", stdout)
	}
}

func TestFunctionalVersionWithOtherArgs(t *testing.T) {
	binaryPath := buildTestBinary(t)

	// --version should only work when it's the only argument
	stdout, stderr, exitCode := runBinary(t, binaryPath, "--version", "extra")

	if exitCode != 0 {
		t.Errorf("exit code = %d, want 0", exitCode)
	}
	if stderr != "" {
		t.Errorf("stderr = %q, want empty", stderr)
	}
	if stdout != "--version extra\n" {
		t.Errorf("stdout = %q, want %q", stdout, "--version extra\n")
	}
}

// Test with environment variables
func TestFunctionalWithEnv(t *testing.T) {
	binaryPath := buildTestBinary(t)

	cmd := exec.Command(binaryPath, "hello", "world")
	cmd.Env = append(os.Environ(), "TEST_VAR=test_value")

	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to run binary with env: %v", err)
	}

	expected := "hello world\n"
	if string(output) != expected {
		t.Errorf("output = %q, want %q", string(output), expected)
	}
}

// Benchmark functional test
func BenchmarkFunctionalEcho(b *testing.B) {
	binaryPath := buildTestBinary(&testing.T{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd := exec.Command(binaryPath, "hello", "world", "benchmark")
		cmd.Run()
	}
}
