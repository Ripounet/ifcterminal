package ifcterminal

import (
	"reflect"
	"text/template"
	"bytes"
	"go/format"
)

func GenerateInterfaceTerminalStruct(ifc reflect.Type) []byte {
	if ifc == nil {
		panic("Please provide a non-nil go interface type")
	}
	if ifc.Kind() != reflect.Interface {
		panic("Please provide a go interface type, not a " + ifc.Kind().String())
	}

	var ifcTmpl = template.Must(template.New("").Funcs(template.FuncMap{
		"methods": methodsOf,
		"ins":     inArgumentsOf,
		"outs":    outArgumentsOf,
	}).Parse(ifcTmplString))

	var gencode bytes.Buffer

	err := ifcTmpl.Execute(&gencode, ifc)
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
type {{.Name}}_X struct {
	Methods struct {
		{{- range methods .}}
		{{.Name}} {{.Type}}
		{{- end}}
	}
}
{{range methods .}}
func (terminal {{$.Name}}_X) {{.Name}}({{- range $i, $arg := ins .Type -}}
	a{{$i}} {{$arg}},
{{- end}}) ({{- range outs .Type -}}
	{{.}}, 
{{- end}}) {
	return terminal.Methods.{{.Name}}({{- range $i, $arg := ins .Type -}}
		a{{$i}},  
	{{- end}})
}
{{end}}
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
