package generator

import (
	"path/filepath"
	"sort"

	"github.com/clear-street/gogen-avro/imprt"
)

// Package represents the output package
type Package struct {
	root  string
	name  string
	files map[string]*File
}

func NewPackage(root, name string) *Package {
	return &Package{root: root, name: name, files: make(map[string]*File)}
}

func (p *Package) Root() string {
	return p.root
}

func (p *Package) Name() string {
	return p.name
}

func (p *Package) WriteFiles(targetDir string) error {
	for _, name := range p.Files() {
		f := p.files[name]
		err := f.WriteFile(imprt.Pkg(p.root, p.name), filepath.Join(targetDir, f.name))
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Package) Files() []string {
	files := make([]string, 0, len(p.files))
	for file, _ := range p.files {
		files = append(files, file)
	}
	sort.Strings(files)
	return files
}

func (p *Package) File(name string) (*File, bool) {
	file, ok := p.files[name]
	return file, ok
}

func (p *Package) AddHeader(file, header string) {
	f, ok := p.files[file]
	if !ok {
		f = NewFile(file)
		p.files[file] = f
	}

	f.headers = append(f.headers, header)
}

func (p *Package) AddFunction(file, str, name, def string) {
	f, ok := p.files[file]
	if !ok {
		f = NewFile(file)
		p.files[file] = f
	}
	f.functions[FunctionName{str, name}] = def
}

func (p *Package) AddStruct(file, name, def string) {
	f, ok := p.files[file]
	if !ok {
		f = NewFile(file)
		p.files[file] = f
	}
	f.structs[name] = def
}

func (p *Package) AddImport(file, name string) {
	f, ok := p.files[file]
	if !ok {
		f = NewFile(file)
		p.files[file] = f
	}
	f.imports[name] = 1
}

func (p *Package) AddConstant(file, name string, value interface{}) {
	f, ok := p.files[file]
	if !ok {
		f = NewFile(file)
		p.files[file] = f
	}
	f.constants[name] = value
}

func (p *Package) HasStruct(file, name string) bool {
	f, ok := p.files[file]
	if !ok {
		return false
	}
	_, ok = f.structs[name]
	return ok
}

func (p *Package) HasFunction(file, str, name string) bool {
	f, ok := p.files[file]
	if !ok {
		return false
	}
	_, ok = f.functions[FunctionName{str, name}]
	return ok
}

func (p *Package) HasImport(file, name string) bool {
	f, ok := p.files[file]
	if !ok {
		return false
	}
	_, ok = f.imports[name]
	return ok
}
