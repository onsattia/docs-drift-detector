// Package registry reads and validates sources.yaml, which maps documentation folders to the code they describe.
package registry

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Registry is a validated sources.yaml file.
type Registry struct {
	Projects map[string]Project
}

// Project maps one docs folder to the code repositories it describes.
type Project struct {
	Docs    string
	OpenAPI string
	Code    []CodeSource
}

// CodeSource is one code repository that a project's docs describe.
type CodeSource struct {
	Repo string
	Ref  string
}

// ValidationError lists every problem found in a registry file, one per entry.
type ValidationError struct {
	Errors []string
}

func (e *ValidationError) Error() string {
	return strings.Join(e.Errors, "\n")
}

// Load reads and validates the registry at path. It reports every problem in one error.
func Load(path string) (*Registry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read registry %s: %w", path, err)
	}

	// The YAML is untyped until validated, so it is decoded into generic maps first.
	var raw any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("%s: invalid YAML: %w", path, err)
	}
	return validate(raw)
}
