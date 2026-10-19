package registry

import "strings"

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
