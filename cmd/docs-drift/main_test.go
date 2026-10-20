package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestBuildScriptInjectsVersionFromVersionFile builds through scripts/build.sh, the same path CI uses,
// and checks that the binary reports the version recorded in VERSION.
func TestBuildScriptInjectsVersionFromVersionFile(t *testing.T) {
	raw, err := os.ReadFile("../../VERSION")
	if err != nil {
		t.Fatal(err)
	}
	want := strings.TrimSpace(string(raw))

	bin := filepath.Join(t.TempDir(), "docs-drift")
	if out, err := exec.Command("sh", "../../scripts/build.sh", bin).CombinedOutput(); err != nil {
		t.Fatalf("build.sh failed: %v\n%s", err, out)
	}

	out, err := exec.Command(bin, "--version").Output()
	if err != nil {
		t.Fatalf("--version failed: %v", err)
	}
	if !strings.Contains(string(out), want) {
		t.Fatalf("--version = %q, want it to contain %q", out, want)
	}
}
