package registry

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

var repoPattern = regexp.MustCompile(`^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+$`)

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

func validate(raw any) (*Registry, error) {
	v := &validator{}

	top, ok := raw.(map[string]any)
	if !ok {
		v.add("", "", "registry must be a mapping with a top-level projects key")
		return nil, v.err()
	}
	v.checkKeys("", "", top, []string{"projects"})

	rawProjects, present := top["projects"]
	if !present {
		v.add("", "projects", "field required")
		return nil, v.err()
	}
	projects, ok := rawProjects.(map[string]any)
	if !ok {
		v.add("", "projects", "must be a mapping of project names to projects")
		return nil, v.err()
	}

	reg := &Registry{Projects: map[string]Project{}}
	for _, name := range sortedKeys(projects) {
		reg.Projects[name] = v.project(name, projects[name])
	}
	if err := v.err(); err != nil {
		return nil, err
	}
	return reg, nil
}

func (v *validator) project(name string, raw any) Project {
	fields, ok := raw.(map[string]any)
	if !ok {
		v.add(name, "", "must be a mapping")
		return Project{}
	}
	v.checkKeys(name, "", fields, []string{"docs", "openapi", "code"})

	var p Project
	p.Docs, _ = v.requiredString(name, "", fields, "docs")
	v.checkRelativePath(name, "docs", p.Docs)
	if _, present := fields["openapi"]; present {
		p.OpenAPI, _ = v.requiredString(name, "", fields, "openapi")
		v.checkRelativePath(name, "openapi", p.OpenAPI)
	}

	rawCode, present := fields["code"]
	if !present {
		v.add(name, "code", "at least one code source is required")
		return p
	}
	list, ok := rawCode.([]any)
	if !ok {
		v.add(name, "code", "must be a list of code sources")
		return p
	}
	if len(list) == 0 {
		v.add(name, "code", "at least one code source is required")
	}
	for i, item := range list {
		p.Code = append(p.Code, v.codeSource(name, fmt.Sprintf("code[%d]", i), item))
	}
	return p
}

func (v *validator) codeSource(project, path string, raw any) CodeSource {
	fields, ok := raw.(map[string]any)
	if !ok {
		v.add(project, path, "must be a mapping")
		return CodeSource{}
	}
	v.checkKeys(project, path, fields, []string{"repo", "ref"})

	var s CodeSource
	s.Repo, _ = v.requiredString(project, path, fields, "repo")
	if s.Repo != "" && !repoPattern.MatchString(s.Repo) {
		v.add(project, joinPath(path, "repo"), fmt.Sprintf("repo '%s' must use owner/name form", s.Repo))
	}
	s.Ref, _ = v.requiredString(project, path, fields, "ref")
	if _, present := fields["ref"]; present && strings.TrimSpace(s.Ref) == "" {
		v.add(project, joinPath(path, "ref"), "ref must not be empty")
	}
	return s
}

// checkRelativePath rejects absolute paths and paths that climb out of the docs repository with "..".
func (v *validator) checkRelativePath(project, field, value string) {
	if value == "" {
		return
	}
	if strings.HasPrefix(value, "/") || contains(strings.Split(value, "/"), "..") {
		v.add(project, field, fmt.Sprintf("path '%s' must stay within the repository: use a relative path without '..'", value))
	}
}

// requiredString returns the string at fields[key] and reports an error if it is missing or not a string.
func (v *validator) requiredString(project, path string, fields map[string]any, key string) (string, bool) {
	raw, present := fields[key]
	if !present {
		v.add(project, joinPath(path, key), "field required")
		return "", false
	}
	s, ok := raw.(string)
	if !ok {
		v.add(project, joinPath(path, key), "must be a string")
		return "", false
	}
	return s, true
}

// checkKeys reports fields that are not allowed. Credential-like names get their own message.
func (v *validator) checkKeys(project, path string, fields map[string]any, allowed []string) {
	for _, key := range sortedKeys(fields) {
		switch {
		case isCredentialKey(key):
			v.add(project, joinPath(path, key),
				fmt.Sprintf("field '%s' is not allowed; credentials belong in the secret store", key))
		case !contains(allowed, key):
			v.add(project, joinPath(path, key), "unknown field")
		}
	}
}

var credentialWords = map[string]bool{"token": true, "password": true, "secret": true, "key": true}

func isCredentialKey(key string) bool {
	parts := strings.FieldsFunc(strings.ToLower(key), func(r rune) bool { return r == '_' || r == '-' })
	for _, part := range parts {
		if credentialWords[part] {
			return true
		}
	}
	return false
}

type validator struct {
	errs []string
}

func (v *validator) add(project, field, message string) {
	var b strings.Builder
	if project != "" {
		b.WriteString(project + ": ")
	}
	if field != "" {
		b.WriteString(field + ": ")
	}
	b.WriteString(message)
	v.errs = append(v.errs, b.String())
}

func (v *validator) err() error {
	if len(v.errs) == 0 {
		return nil
	}
	return &ValidationError{Errors: v.errs}
}

func joinPath(parent, key string) string {
	if parent == "" {
		return key
	}
	return parent + "." + key
}

func sortedKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func contains(values []string, want string) bool {
	for _, v := range values {
		if v == want {
			return true
		}
	}
	return false
}
