package registry

import "strings"

// Registry is a validated sources.yaml file.
//
// Projects is a map, so looking up a project by name is direct. Go maps have no
// guaranteed iteration order, which is why the code sorts keys before printing or reporting.
type Registry struct {
	Projects map[string]Project
}

// Project maps one docs folder to the code repositories it describes.
//
// Code is a slice (Go's growable list, like a Python list). OpenAPI is a plain string;
// an empty string means "not set", because Go strings cannot be null.
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
//
// It implements Go's error interface (a type with an Error() string method), so it can be
// returned like any other error. Callers recover it with errors.As to read the list.
type ValidationError struct {
	Errors []string
}

func (e *ValidationError) Error() string {
	// Joining with newlines makes the error print one problem per line.
	return strings.Join(e.Errors, "\n")
}
