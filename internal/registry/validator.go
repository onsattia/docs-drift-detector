package registry

import (
	"fmt"
	"strings"
)

// validator collects error messages so that one pass reports every problem in the file.
type validator struct {
	errs []string
}

// add records one problem, prefixed with the project and field path when they are known.
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

// err returns a ValidationError holding every recorded problem, or nil when there are none.
func (v *validator) err() error {
	if len(v.errs) == 0 {
		return nil
	}
	return &ValidationError{Errors: v.errs}
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

// checkRelativePath rejects absolute paths and paths that climb out of the docs repository with "..".
func (v *validator) checkRelativePath(project, field, value string) {
	if value == "" {
		return
	}
	if strings.HasPrefix(value, "/") || contains(strings.Split(value, "/"), "..") {
		v.add(project, field, fmt.Sprintf("path '%s' must stay within the repository: use a relative path without '..'", value))
	}
}

var credentialWords = map[string]bool{"token": true, "password": true, "secret": true, "key": true}

// isCredentialKey reports whether any word in the key, split on _ or -, is a credential word.
func isCredentialKey(key string) bool {
	parts := strings.FieldsFunc(strings.ToLower(key), func(r rune) bool { return r == '_' || r == '-' })
	for _, part := range parts {
		if credentialWords[part] {
			return true
		}
	}
	return false
}
