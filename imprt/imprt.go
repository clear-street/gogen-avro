package imprt

import (
	"fmt"
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

func IsRootPkg(root, ns string) bool {
	ar := strings.Split(ns, ".")
	return len(ar) <= 1
}

func Pkg(root, ns string) string {
	ar := strings.Split(ns, ".")
	br := strings.Split(root, "/")
	if len(ar) <= 1 {
		return br[len(br)-1]
	}

	return ar[len(ar)-1]
}

func Type(root, ns, typename string) string {
	pkg := Pkg(root, ns)
	if typename[0] == '*' {
		return fmt.Sprintf("*%v.%v", pkg, typename[1:])
	}
	return fmt.Sprintf("%v.%v", pkg, typename)
}

func UniqName(root, ns, name string) string {
	pkg := Pkg(root, ns)
	return fmt.Sprintf("%v%v", strings.Title(pkg), name)
}
