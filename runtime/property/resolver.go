package property

import (
	"fmt"
)

// Resolver interface provides methods for resolving properties.
type Resolver interface {
	ContainsProperty(name string) bool
	Property(name string) (any, bool)
	PropertyOrDefault(name string, defaultValue any) any
	ResolvePlaceholders(text string) string
	ResolveRequiredPlaceholders(text string) (string, error)
}

// SourcesResolver is an implementation of the Resolver interface.
// It resolves properties from the given sources.
type SourcesResolver struct {
	sources *Sources
}

// NewSourcesResolver creates a new SourcesResolver with the given sources.
func NewSourcesResolver(sources *Sources) *SourcesResolver {
	if sources == nil {
		panic("nil sources")
	}

	return &SourcesResolver{
		sources: sources,
	}
}

// ContainsProperty checks if the given property name exists in the sources.
func (r *SourcesResolver) ContainsProperty(name string) bool {
	return r.sources.Contains(name)
}

// Property returns the value of the given property name from the sources.
func (r *SourcesResolver) Property(name string) (any, bool) {
	for _, source := range r.sources.ToSlice() {
		if value, ok := source.Property(name); ok {
			return value, true
		}
	}

	return nil, false
}

// PropertyOrDefault returns the value of the given property name from the sources.
// If the property does not exist, it returns the default value.
func (r *SourcesResolver) PropertyOrDefault(name string, defaultValue any) any {
	for _, source := range r.sources.ToSlice() {
		if value, ok := source.Property(name); ok {
			return value.(string)
		}
	}

	return defaultValue
}

// ResolvePlaceholders resolves placeholders in the given text.
// If a placeholder cannot be resolved, it continues to resolve other placeholders.
func (r *SourcesResolver) ResolvePlaceholders(s string) string {
	result, _ := r.resolveRequiredPlaceHolders(s, true)
	return result
}

// ResolveRequiredPlaceholders resolves placeholders in the given text.
// If a placeholder cannot be resolved, it returns an error.
func (r *SourcesResolver) ResolveRequiredPlaceholders(s string) (string, error) {
	return r.resolveRequiredPlaceHolders(s, false)
}

// resolveRequiredPlaceHolders resolves placeholders in the given text.
// If continueOnError is true, it continues to resolve placeholders even if a placeholder cannot be resolved.
func (r *SourcesResolver) resolveRequiredPlaceHolders(s string, continueOnError bool) (string, error) {
	var buf []byte

	i := 0
	for j := 0; j < len(s); j++ {
		if s[j] == '$' && j+1 < len(s) {
			if buf == nil {
				buf = make([]byte, 0, 2*len(s))
			}

			buf = append(buf, s[i:j]...)
			name, w := r.getPlaceholderName(s[j+1:])

			if name == "" && w > 0 {
			} else if name == "" {
				buf = append(buf, s[j])
			} else {
				value, ok := r.Property(name)

				if !ok && !continueOnError {
					return "", fmt.Errorf("cannot resolve placeholder '%s'", s[j:i+w+1])
				}

				stringValue, canConvert := value.(string)
				if !canConvert && !continueOnError {
					return "", fmt.Errorf("string values can only be used as placeholder '%s'", s[j:i+w+1])
				}

				if continueOnError {
					buf = append(buf, s[j:i+w+1]...)
				} else {
					buf = append(buf, stringValue...)
				}
			}

			j += w
			i = j + 1
		}
	}

	if buf == nil {
		return s, nil
	}

	return string(buf) + s[i:], nil
}

func (r *SourcesResolver) getPlaceholderName(s string) (string, int) {
	switch {
	case s[0] == '{':
		if len(s) > 2 && isSpecialVar(s[1]) && s[2] == '}' {
			return s[1:2], 3
		}

		for i := 1; i < len(s); i++ {
			if s[i] == '}' {
				if i == 1 {
					return "", 2
				}
				return s[1:i], i + 1
			}
		}
		return "", 1
	case isSpecialVar(s[0]):
		return s[0:1], 1
	}

	var i int
	for i = 0; i < len(s) && isAlphaNum(s[i]); i++ {
	}

	return s[:i], i
}

func isSpecialVar(c uint8) bool {
	switch c {
	case '*', '#', '$', '@', '!', '?', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	}
	return false
}

func isAlphaNum(c uint8) bool {
	return c == '_' || '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}
