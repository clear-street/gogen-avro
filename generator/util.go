package generator

import (
	"regexp"
	"strings"
	"unicode"
)

var numberSequence = regexp.MustCompile(`([a-zA-Z])(\d+)([a-zA-Z]?)`)
var numberReplacement = []byte(`$1 $2 $3`)

func addWordBoundariesToNumbers(s string) string {
	b := []byte(s)
	b = numberSequence.ReplaceAll(b, numberReplacement)
	return string(b)
}

// Converts a string to CamelCase
func toCamelInitCase(s string, initCase bool) string {
	s = addWordBoundariesToNumbers(s)
	s = strings.Trim(s, " ")
	n := ""
	capNext := initCase
	for _, v := range s {
		if v >= 'A' && v <= 'Z' {
			n += string(v)
		}
		if v >= '0' && v <= '9' {
			n += string(v)
		}
		if v >= 'a' && v <= 'z' {
			if capNext {
				n += strings.ToUpper(string(v))
			} else {
				n += string(v)
			}
		}
		if v == '_' || v == ' ' || v == '-' {
			capNext = true
		} else {
			capNext = false
		}
	}
	return n
}

// ToCamel converts a string to CamelCase
func toCamel(s string) string {
	return toCamelInitCase(s, true)
}

// ToPublicName returns a go-idiomatic public name. The Avro spec specifies names must start with [A-Za-z_] and contain [A-Za-z0-9_].
// The golang spec says valid identifiers start with [A-Za-z_] and contain [A-Za-z0-9], but the first character must be [A-Z] for the field to be public.
func ToPublicName(name string) string {
	// bleh: https://github.com/golang/go/wiki/CodeReviewComments#initialisms
	if strings.HasSuffix(name, "_id") {
		name = strings.TrimSuffix(name, "_id") + "ID"
	}
	if name == "id" {
		return "ID"
	}

	return toCamel(name)
	//return namer.ToPublicName(name)
}

// ToPublicSimpleName returns a go-idiomatic public name. The Avro spec
// specifies names must start with [A-Za-z_] and contain [A-Za-z0-9_].
// The golang spec says valid identifiers start with [A-Za-z_] and contain
// [A-Za-z0-9], but the first character must be [A-Z] for the field to be
// public.
func ToPublicSimpleName(name string) string {
	lastIndex := strings.LastIndex(name, ".")
	name = name[lastIndex+1:]
	return strings.Title(strings.Trim(name, "_"))
}

// ToSnake makes filenames snake-case, taken from https://gist.github.com/elwinar/14e1e897fdbe4d3432e1
func ToSnake(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}

	return string(out)
}
