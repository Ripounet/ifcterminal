package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/importer"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

func usage() {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "  "+os.Args[0]+" Ifc")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Ifc is a fully-qualified interface name, like net/http.CloseNotifier")
	os.Exit(1)
}

// Problem: reflect and go/types don't play together.
//
// Approach 1: use go/types, generate go code, compile it, run it (it will use reflect).
//
// Approach 2: use go/types only.
//
// Approach 3: hijack go doc, do some regexp.

func main() {
	flag.Usage = usage
	flag.Parse()

	if len(flag.Args()) != 1 {
		usage()
	}
	ifcPath := flag.Arg(0)

	parts := strings.Split(ifcPath, ".")
	if len(parts) == 1 {
		fmt.Fprintln(os.Stderr, "Package name is mandatory, as in net/http.CloseNotifier")
		os.Exit(1)
	}
	if len(parts) > 2 {
		fmt.Fprintln(os.Stderr, "Packages must be slash-separated, as in net/http.CloseNotifier")
		os.Exit(1)
	}

	packagePath, ifcName := parts[0], parts[1]

	p, err := importer.Default().Import(packagePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	obj := p.Scope().Lookup(ifcName)
	if obj == nil {
		fmt.Fprintln(os.Stderr, "Interface", ifcName, "not found in package", packagePath)
		os.Exit(1)
	}
	ty := obj.Type()
	if ty == nil {
		fmt.Fprintln(os.Stderr, "Could not determine type of", ifcPath)
		os.Exit(1)
	}

	// log.Println(obj)
	// log.Println(ty)
	// log.Println(ty.Underlying())
	// log.Println(ty.Underlying().Underlying())

	code := generateRunner(packagePath, ifcName)
	err = runCode(code)
	if err != nil {
		panic(err)
	}
	fmt.Println()
}

//
// generateRunner, runCode, runnerTmplString :
// this is Approach 1.
//

func generateRunner(ifcPackage, ifcName string) (code []byte) {
	runnerTmpl := template.Must(template.New("").Parse(runnerTmplString))

	var buffer bytes.Buffer

	err := runnerTmpl.Execute(&buffer, struct {
		IfcPackage, IfcName string
	}{
		IfcPackage: ifcPackage,
		IfcName:    ifcName,
	})
	if err != nil {
		panic(err)
	}

	return buffer.Bytes()
}

func runCode(code []byte) error {
	dir, err := ioutil.TempDir("", "ifcterminal-")
	if err != nil {
		return err
	}
	// log.Println("Temp dir ", dir)
	defer os.RemoveAll(dir)

	fpath := filepath.Join(dir, "runner.go")
	// log.Println("Temp file ", fpath)
	err = ioutil.WriteFile(fpath, code, 0777)
	if err != nil {
		return err
	}
	log.Println("Written to ", fpath)
	command := exec.Command("go", "run", fpath)
	command.Stdout = os.Stdout
	return command.Run()
}

const runnerTmplString = `
package main

import (
	"os"
	"reflect"
	"github.com/Ripounet/ifcterminal"
	p "{{.IfcPackage}}"
)

func main() {
	var x *p.{{.IfcName}}
	ty := reflect.TypeOf(x).Elem()
	code := ifcterminal.GenerateInterfaceTerminalStruct(ty)
	os.Stdout.Write(code)
}
`
