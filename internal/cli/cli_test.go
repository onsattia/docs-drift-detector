package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func run(t *testing.T, args ...string) (int, string, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := Execute(args, &stdout, &stderr, "1.2.3")
	return code, stdout.String(), stderr.String()
}

func TestHelpListsUsageAndExitsZero(t *testing.T) {
	code, stdout, _ := run(t, "--help")

	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if !strings.Contains(stdout, "Usage") {
		t.Fatalf("help output missing usage section:\n%s", stdout)
	}
}

func TestVersionPrintsVersionAndExitsZero(t *testing.T) {
	code, stdout, _ := run(t, "--version")

	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if !strings.Contains(stdout, "1.2.3") {
		t.Fatalf("version output = %q, want it to contain 1.2.3", stdout)
	}
}

func TestUnknownCommandPrintsUsageToStderrAndExitsNonZero(t *testing.T) {
	code, stdout, stderr := run(t, "nope")

	if code == 0 {
		t.Fatal("exit code = 0, want non-zero")
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty for an error", stdout)
	}
	if !strings.Contains(stderr, "Usage") {
		t.Fatalf("stderr missing usage message:\n%s", stderr)
	}
}

func TestRegistryValidatePrintsSummaryAndExitsZero(t *testing.T) {
	code, stdout, stderr := run(t, "registry", "validate", "../../examples/sources.yaml")

	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr:\n%s", code, stderr)
	}
	if !strings.Contains(stdout, "billing-docs") || !strings.Contains(stdout, "acme/billing-api@main") {
		t.Fatalf("summary missing project or code source:\n%s", stdout)
	}
}

func TestRegistryValidateReportsErrorsAndExitsNonZero(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sources.yaml")
	if err := os.WriteFile(path, []byte("projects:\n  billing-docs:\n    docs: docs/billing\n    code: []\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	code, stdout, stderr := run(t, "registry", "validate", path)

	if code == 0 {
		t.Fatal("exit code = 0, want non-zero")
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty on validation error", stdout)
	}
	if !strings.Contains(stderr, "billing-docs") || !strings.Contains(stderr, "at least one code source") {
		t.Fatalf("stderr missing the validation error:\n%s", stderr)
	}
}
