package registry

import (
	"fmt"
	"regexp"
	"strings"
)

var repoPattern = regexp.MustCompile(`^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+$`)

// requiredString returns the string at fields[key] and reports an error if it is missing or not a string.
func (v *errorList) requiredString(project, path string, fields map[string]any, key string) (string, bool) {
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
func (v *errorList) checkKeys(project, path string, fields map[string]any, allowed []string) {
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
func (v *errorList) checkRelativePath(project, field, value string) {
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
