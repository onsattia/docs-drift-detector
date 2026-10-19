package registry

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeRegistry(t *testing.T, text string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "sources.yaml")
	if err := os.WriteFile(path, []byte(text), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

// validationErrors returns the individual errors from a failed Load, or fails the test.
func validationErrors(t *testing.T, err error) []string {
	t.Helper()
	var verr *ValidationError
	if !errors.As(err, &verr) {
		t.Fatalf("error = %v, want *ValidationError", err)
	}
	return verr.Errors
}

func TestValidRegistryExposesProjectsAndCodeSources(t *testing.T) {
	path := writeRegistry(t, `
projects:
  billing-docs:
    docs: docs/billing
    code:
      - repo: acme/billing-api
        ref: main
`)

	reg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	project := reg.Projects["billing-docs"]
	if project.Docs != "docs/billing" {
		t.Fatalf("Docs = %q, want docs/billing", project.Docs)
	}
	if len(project.Code) != 1 || project.Code[0].Repo != "acme/billing-api" || project.Code[0].Ref != "main" {
		t.Fatalf("Code = %+v, want one acme/billing-api@main source", project.Code)
	}
}

func TestProjectWithoutCodeSourceIsRejected(t *testing.T) {
	path := writeRegistry(t, `
projects:
  billing-docs:
    docs: docs/billing
    code: []
`)

	_, err := Load(path)

	errs := validationErrors(t, err)
	if !anyContains(errs, "billing-docs", "at least one code source") {
		t.Fatalf("errors = %q, want one naming billing-docs and the missing code source", errs)
	}
}

func TestAllErrorsAreReportedInOneRun(t *testing.T) {
	path := writeRegistry(t, `
projects:
  billing-docs:
    docs: docs/billing
    code:
      - repo: not-a-repo
        ref: main
  search-docs:
    docs: docs/search
    code:
      - repo: acme/search-api
`)

	_, err := Load(path)

	errs := validationErrors(t, err)
	if len(errs) != 2 {
		t.Fatalf("got %d errors, want 2: %q", len(errs), errs)
	}
	if !anyContains(errs, "billing-docs", "code[0].repo") {
		t.Fatalf("errors = %q, want a billing-docs code[0].repo error", errs)
	}
	if !anyContains(errs, "search-docs", "code[0].ref") {
		t.Fatalf("errors = %q, want a search-docs code[0].ref error", errs)
	}
}

func anyContains(errs []string, parts ...string) bool {
	for _, e := range errs {
		all := true
		for _, p := range parts {
			if !strings.Contains(e, p) {
				all = false
				break
			}
		}
		if all {
			return true
		}
	}
	return false
}

func TestCredentialFieldsAreRejectedWithoutEchoingValue(t *testing.T) {
	const secretValue = "s3cr3t-value-123"
	for _, field := range []string{"token", "password", "secret", "api_key"} {
		t.Run(field, func(t *testing.T) {
			path := writeRegistry(t, `
projects:
  billing-docs:
    docs: docs/billing
    code:
      - repo: acme/billing-api
        ref: main
        `+field+`: `+secretValue+`
`)

			_, err := Load(path)

			errs := validationErrors(t, err)
			message := strings.Join(errs, "\n")
			if !strings.Contains(message, field) {
				t.Fatalf("errors = %q, want the field name %q", errs, field)
			}
			if strings.Contains(message, secretValue) {
				t.Fatalf("errors echo the secret value: %q", errs)
			}
		})
	}
}

func TestDocsAndOpenAPIPathsMustStayWithinRepository(t *testing.T) {
	tests := []struct {
		name  string
		field string
		value string
	}{
		{"docs parent escape", "docs", "../other-repo/docs"},
		{"docs absolute", "docs", "/etc/docs"},
		{"openapi parent escape", "openapi", "../billing-api/openapi.yaml"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			docs := "docs/billing"
			openapi := ""
			if tt.field == "docs" {
				docs = tt.value
			} else {
				openapi = "    openapi: " + tt.value + "\n"
			}
			path := writeRegistry(t, "projects:\n  billing-docs:\n    docs: "+docs+"\n"+openapi+
				"    code:\n      - repo: acme/billing-api\n        ref: main\n")

			_, err := Load(path)

			errs := validationErrors(t, err)
			if !anyContains(errs, "billing-docs", tt.field, "must stay within the repository") {
				t.Fatalf("errors = %q, want a billing-docs %s path error", errs, tt.field)
			}
		})
	}
}

func TestInvalidRegistriesAreRejectedWithSpecificErrors(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{"top level is not a mapping", "- billing-docs\n", "registry must be a mapping"},
		{"projects key missing", "other: 1\n", "projects: field required"},
		{"projects is not a mapping", "projects: nope\n", "projects: must be a mapping"},
		{"project is not a mapping", "projects:\n  billing-docs: nope\n", "billing-docs: must be a mapping"},
		{"unknown project field", "projects:\n  billing-docs:\n    docs: docs/billing\n    extra: 1\n" +
			"    code:\n      - repo: acme/billing-api\n        ref: main\n", "billing-docs: extra: unknown field"},
		{"docs missing", "projects:\n  billing-docs:\n    code:\n      - repo: acme/billing-api\n        ref: main\n",
			"billing-docs: docs: field required"},
		{"code missing", "projects:\n  billing-docs:\n    docs: docs/billing\n",
			"billing-docs: code: at least one code source is required"},
		{"code is not a list", "projects:\n  billing-docs:\n    docs: docs/billing\n    code: nope\n",
			"billing-docs: code: must be a list"},
		{"code item is not a mapping", "projects:\n  billing-docs:\n    docs: docs/billing\n    code:\n      - nope\n",
			"billing-docs: code[0]: must be a mapping"},
		{"unknown code source field", "projects:\n  billing-docs:\n    docs: docs/billing\n    code:\n" +
			"      - repo: acme/billing-api\n        ref: main\n        branch: x\n", "billing-docs: code[0].branch: unknown field"},
		{"ref is blank", "projects:\n  billing-docs:\n    docs: docs/billing\n    code:\n" +
			"      - repo: acme/billing-api\n        ref: \"  \"\n", "billing-docs: code[0].ref: ref must not be empty"},
		{"repo is not a string", "projects:\n  billing-docs:\n    docs: docs/billing\n    code:\n" +
			"      - repo: [a]\n        ref: main\n", "billing-docs: code[0].repo: must be a string"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Load(writeRegistry(t, tt.text))

			errs := validationErrors(t, err)
			if !anyContains(errs, tt.want) {
				t.Fatalf("errors = %q, want one containing %q", errs, tt.want)
			}
		})
	}
}

func TestInvalidYAMLAndMissingFileAreReportedAsReadErrors(t *testing.T) {
	_, err := Load(writeRegistry(t, "projects: [unclosed\n"))
	if err == nil || !strings.Contains(err.Error(), "invalid YAML") {
		t.Fatalf("error = %v, want an invalid YAML error", err)
	}

	_, err = Load(filepath.Join(t.TempDir(), "missing.yaml"))
	if err == nil || !strings.Contains(err.Error(), "read registry") {
		t.Fatalf("error = %v, want a read registry error", err)
	}
}
