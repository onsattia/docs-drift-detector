// Package registry reads and validates sources.yaml, which maps documentation folders to the code they describe.
package registry

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

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
