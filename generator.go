package ifcterminal

import (
	"bytes"
	"go/format"
	"reflect"
	"text/template"
)

//
// Using Reflection.
//

func GenerateInterfaceTerminalStruct(ifc reflect.Type) []byte {
	if ifc == nil {
		panic("Please provide a non-nil go interface type")
	}
	if ifc.Kind() != reflect.Interface {
		panic("Please provide a go interface type, not a " + ifc.Kind().String())
	}

	ifcTmpl := template.Must(template.New("").Funcs(template.FuncMap{
		"methods": methodsOf,
		"ins":     inArgumentsOf,
		"outs":    outArgumentsOf,
	}).Parse(ifcTmplString))

	var gencode bytes.Buffer

	err := ifcTmpl.Execute(&gencode, struct {
		Ifc        reflect.Type
		TypeSuffix string
	}{
		Ifc:        ifc,
		TypeSuffix: "_X",
	})
	if err != nil {
		panic(err)
	}

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
