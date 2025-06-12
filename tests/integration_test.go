package tests

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

const (
	EMPTY = "inputs/empty.txt"
	FILE1 = "inputs/file1.txt"
	FILE2 = "inputs/file2.txt"
	BLANK = "inputs/blank.txt"
)

func BuildCli(t *testing.T) string {
	projectDir := "/home/celal/Projects/Go/gocomm"
	binaryPath := filepath.Join(os.TempDir(), "gocomm-test")

	cmd := exec.Command("go", "build", "-o", binaryPath)
	cmd.Dir = projectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("failed to build gocomm: %v", err)
	}

	if _, err := os.Stat(binaryPath); err != nil {
		t.Fatalf("built binary not found: %v", err)
	}
	return binaryPath
}

func generate_bad_file() string {
	return uuid.New().String()
}

func TestDiesBadFile(t *testing.T) {
	execPath := BuildCli(t)
	badFile := generate_bad_file()
	cmd := exec.Command(execPath, badFile, FILE2)
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected error but got none")
	}
	expectedErr := fmt.Sprintf("Error: open %s: no such file or directory\n", badFile)
	if string(output) != expectedErr {
		t.Errorf("got %q, want %q", string(output), expectedErr)
	}
	cmd = exec.Command(execPath, FILE1, badFile)
	output, err = cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected error but got none")
	}
	expectedErr = fmt.Sprintf("Error: open %s: no such file or directory\n", badFile)
	if string(output) != expectedErr {
		t.Errorf("got %q, want %q", string(output), expectedErr)
	}
}
func TestDiesNoArgs(t *testing.T) {
	execPath := BuildCli(t)
	cmd := exec.Command(execPath)
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected error but got none")
	}
	expectedErr := "Error: accepts 2 arg(s), received 0\n"
	if string(output) != expectedErr {
		t.Errorf("got %q, want %q", string(output), expectedErr)
	}
}
func TestDiesBothStdin(t *testing.T) {
	execPath := BuildCli(t)
	cmd := exec.Command(execPath, "-", "-")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected error but got none")
	}
	expectedErr := "Error: Both input files cannot be STDIN (\"-\")\n"
	if string(output) != expectedErr {
		t.Errorf("got %q, want %q", string(output), expectedErr)
	}
}

func TestCli(t *testing.T) {
	execPath := BuildCli(t)
	tests := []struct {
		name     string
		args     []string
		expected string
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := exec.Command(execPath, tt.args...).CombinedOutput()
			if err != nil {
				t.Fatalf("gocomm failed: %v\nOutput: %s", err, string(output))
			}
			actual := string(output)
			expected := "not setted!"
			if !reflect.DeepEqual(actual, expected) {
				t.Fatalf("actual = %s, expected = %s", actual, expected)
			}
		})
	}

}
