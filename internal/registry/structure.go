package registry

import (
	"fmt"
	"strings"
)

// validate checks the shape of the untyped YAML and builds the Registry. It returns every error it finds.
//
// The raw YAML arrives as nested map[string]any values. Each step asserts the type it expects
// with raw.(map[string]any). The two-value form returns ok=false instead of panicking when the
// assertion fails, and the code reports that as an error.
func validate(raw any) (*Registry, error) {
	v := &errorList{}

	top, ok := raw.(map[string]any)
	if !ok {
		v.add("", "", "registry must be a mapping with a top-level projects key")
		return nil, v.err()
	}
	// Reject unknown top-level keys before reading the rest, so a typo like "project:" is reported.
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

	// &Registry returns a pointer, so callers share one value. The map must be created here;
	// writing to a nil map panics in Go.
	reg := &Registry{Projects: map[string]Project{}}
	for _, name := range sortedKeys(projects) {
		reg.Projects[name] = v.project(name, projects[name])
	}
	// Collected errors are returned only after every project has been checked. That is how
	// one run reports all problems instead of stopping at the first.
	if err := v.err(); err != nil {
		return nil, err
	}
	return reg, nil
}

func (v *errorList) project(name string, raw any) Project {
	fields, ok := raw.(map[string]any)
	if !ok {
		v.add(name, "", "must be a mapping")
		return Project{}
	}
	v.checkKeys(name, "", fields, []string{"docs", "openapi", "code"})

	// The underscore discards the ok flag. A missing field is already recorded as an error by
	// requiredString, so the empty string it returns does not need a second check.
	var p Project
	p.Docs, _ = v.requiredString(name, "", fields, "docs")
	v.checkRelativePath(name, "docs", p.Docs)
	if _, present := fields["openapi"]; present {
		p.OpenAPI, _ = v.requiredString(name, "", fields, "openapi")
		v.checkRelativePath(name, "openapi", p.OpenAPI)
	}

	// A missing key and an empty list are reported as separate problems, so the message is specific.
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

func (v *errorList) codeSource(project, path string, raw any) CodeSource {
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
