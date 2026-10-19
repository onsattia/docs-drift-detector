package registry

import "strings"

// errorList collects error messages so that one pass reports every problem in the file.
//
// Methods use a pointer receiver (v *errorList) so every call appends to the same slice.
// A value receiver would work on a copy and the messages would be lost.
type errorList struct {
	errs []string
}

// add records one problem, prefixed with the project and field path when they are known.
func (v *errorList) add(project, field, message string) {
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
func (v *errorList) err() error {
	if len(v.errs) == 0 {
		return nil
	}
	return &ValidationError{Errors: v.errs}
}
