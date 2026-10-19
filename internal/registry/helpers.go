package registry

import "sort"

// joinPath builds a field path such as code[0].repo.
func joinPath(parent, key string) string {
	if parent == "" {
		return key
	}
	return parent + "." + key
}

// sortedKeys returns the keys of m in order, so errors come out in a stable order.
//
// [T any] is a type parameter (Go generics). It lets one function accept a map of any value
// type. Without it we would need a copy of this function for each value type.
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
