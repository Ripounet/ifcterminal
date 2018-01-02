package ifcterminal

import (
	"bytes"
	"crypto"
	"go/format"
	"io"
	"net/http"
	"reflect"
	"strings"
	"text/template"
)

//
// Using Reflection.
//

// GenerateInterfaceTerminalStruct produces a full go source file containing an implementation struct
// of a given interface.
//
// typeSuffix: if feel like avoiding type name conflicts:  (interface "Reader", suffix "_X") -> struct "Reader_X"
// targetPackage: set if you want to save the result file in a custom package. Default value is same (leaf) package name as ifc.
func GenerateInterfaceTerminalStruct(ifc reflect.Type, typeSuffix string, targetPackage string) []byte {
	if ifc == nil {
		panic("Please provide a non-nil go interface type")
	}
	if ifc.Kind() != reflect.Interface {
		panic("Please provide a go interface type, not a " + ifc.Kind().String())
	}
	if targetPackage == "" {
		parts := strings.Split(ifc.PkgPath(), "/")
		targetPackage = parts[len(parts)-1]
	}

	ifcTmpl := template.Must(template.New("").Funcs(template.FuncMap{
		"methods": methodsOf,
		"ins":     inArgumentsOf,
		"outs":    outArgumentsOf,
	}).Parse(ifcTmplString))

	var gencode bytes.Buffer

	err := ifcTmpl.Execute(&gencode, struct {
		Ifc           reflect.Type
		TypeSuffix    string
		TargetPackage string
		Imports       []string
	}{
		Ifc:           ifc,
		TypeSuffix:    typeSuffix,
		TargetPackage: targetPackage,
		Imports:       findDependencies(ifc),
	})
	if err != nil {
		panic(err)
	}
	//log.Println(string(gencode.Bytes()))

	// gofmt
	pretty, err := format.Source(gencode.Bytes())
	if err != nil {
		panic(err)
	}
	return pretty
}

// A struct can't have a Method and a Field with the same name.
// So we introduce an intermediate name "Methods".
const ifcTmplString = `
package {{.TargetPackage}}

import(
	{{range .Imports -}}
	"{{.}}"
	{{end}}
)

type {{.Ifc.Name}}{{.TypeSuffix}} struct {
	Methods struct {
		{{- range methods .Ifc}}
		{{.Name}} {{.Type}}
		{{- end}}
	}
}
{{range methods .Ifc}}
func (terminal {{$.Ifc.Name}}{{$.TypeSuffix}}) {{.Name}}({{- range $i, $arg := ins .Type -}}
	a{{$i}} {{$arg}},
{{- end}}) ({{- range $i, $arg := outs .Type -}}
	out{{$i}} {{$arg}},
{{- end}}) {
	if terminal.Methods.{{.Name}} == nil {
		return
	}
	{{if outs .Type}}return {{end}}terminal.Methods.{{.Name}}({{- range $i, $arg := ins .Type -}}
		a{{$i}},  
	{{- end}})
}
{{end}}

// This proves (by compiling) that {{.Ifc.Name}}{{.TypeSuffix}} implements {{.Ifc}}
func init() {
	var x {{.Ifc}} = {{.Ifc.Name}}{{.TypeSuffix}}{}
	_ = x
}

`

func methodsOf(t reflect.Type) (methods []reflect.Method) {
	for i := 0; i < t.NumMethod(); i++ {
		methods = append(methods, t.Method(i))
	}
	return methods
}

func inArgumentsOf(t reflect.Type) (ins []reflect.Type) {
	for i := 0; i < t.NumIn(); i++ {
		ins = append(ins, t.In(i))
	}
	return ins
}

func outArgumentsOf(t reflect.Type) (outs []reflect.Type) {
	for i := 0; i < t.NumOut(); i++ {
		outs = append(outs, t.Out(i))
	}
	return outs
}

func findDependencies(ifc reflect.Type) (packages []string) {
	importSet := make(map[string]bool)
	importSet[ifc.PkgPath()] = true

	var traverse func(t reflect.Type)
	traverse = func(t reflect.Type) {
		if t.Name() != "" {
			// A named type is self-sufficient
			importSet[t.PkgPath()] = true
			return
		}

		switch t.Kind() {
		case reflect.Ptr, reflect.Array, reflect.Chan, reflect.Slice:
			traverse(t.Elem())
			return
		case reflect.Map:
			traverse(t.Key())
			traverse(t.Elem())
			return
		case reflect.Func:
			for i := 0; i < t.NumIn(); i++ {
				argType := t.In(i)
				traverse(argType)
			}
			for i := 0; i < t.NumOut(); i++ {
				outType := t.Out(i)
				traverse(outType)
			}
			return
		default:
			// built-in primitives don't need imports
		}
	}

	for k := 0; k < ifc.NumMethod(); k++ {
		m := ifc.Method(k).Type
		traverse(m)
	}

	// Map to list
	for p := range importSet {
		if p == "" {
			continue
		}
		packages = append(packages, p)
	}
	return packages
}

// SomeFunctionalIfc : Just for the tests
type SomeFunctionalIfc interface {
	Instrument(func(http.CookieJar, io.Reader)) (int, map[reflect.Kind]string, func() crypto.Signer)
}

// SomeRecursiveType : Just for the tests
type SomeRecursiveType interface {
	Wrap(string) SomeRecursiveType
}
