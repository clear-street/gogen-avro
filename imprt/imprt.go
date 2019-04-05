package imprt

import (
	"path/filepath"
	"strings"
)

func Path(root, ns string) string {
	ar := strings.Split(ns, ".")
	if len(ar) <= 1 {
		return root
	}

	path := filepath.Join(root)
	for _, v := range ar[1:] {
		path = filepath.Join(path, v)
	}

	return path
}

func Pkg(root, ns string) string {
	ar := strings.Split(ns, ".")
	if len(ar) <= 1 {
		return root
	}

	return ar[len(ar)-1]
}
