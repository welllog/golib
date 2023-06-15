package strz

import "strings"

// KeyGenerator is a key generator.
type KeyGenerator struct {
	prefix    string
	delimiter string
}

// NewKeyGenerator creates a new key generator.
func NewKeyGenerator(delimiter string, parts ...string) KeyGenerator {
	prefix := strings.Join(parts, delimiter)
	if prefix != "" {
		prefix += delimiter
	}
	return KeyGenerator{
		prefix:    prefix,
		delimiter: delimiter,
	}
}

// Generate generates a key.
func (k KeyGenerator) Generate(parts ...string) string {
	return k.prefix + strings.Join(parts, k.delimiter)
}

// Spread returns the spread of the key.
func (k KeyGenerator) Spread() []string {
	prefixes := strings.Split(k.prefix, k.delimiter)
	if len(prefixes) <= 1 {
		return nil
	}
	return prefixes[:len(prefixes)-1]
}

// With returns a new key generator with the given parts.
func (k KeyGenerator) With(parts ...string) KeyGenerator {
	if len(parts) == 0 {
		return k
	}
	parts = append(k.Spread(), parts...)
	return NewKeyGenerator(k.delimiter, parts...)
}
