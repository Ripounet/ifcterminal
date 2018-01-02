package ifcterminal_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/Ripounet/ifcterminal"
)

func TestInterfaceTerminalStruct(t *testing.T) {
	{
		var x *io.Reader
		ty := reflect.TypeOf(x).Elem()
		code := ifcterminal.GenerateInterfaceTerminalStruct(ty, "", "")
		err := checkCode(code)
		if err != nil {
			t.Errorf("error in compiling generated code: %v \n\n%s", err, string(code))
		}
	}
	{
		var x *io.ReadWriter
		ty := reflect.TypeOf(x).Elem()
		code := ifcterminal.GenerateInterfaceTerminalStruct(ty, "", "mypackage")
		packageLOC := []byte("package mypackage")
		if !bytes.Contains(code, packageLOC) {
			t.Errorf("Expected %q in\n\n %s", string(packageLOC), string(code))
		}
		err := checkCode(code)
		if err != nil {
			t.Errorf("error in compiling generated code: %v \n\n%s", err, string(code))
		}
	}
	{
		var x *io.WriterTo
		ty := reflect.TypeOf(x).Elem()
		code := ifcterminal.GenerateInterfaceTerminalStruct(ty, "", "")
		err := checkCode(code)
		if err != nil {
			t.Errorf("error in compiling generated code: %v \n\n%s", err, string(code))
		}
	}
	{
		var x *http.CloseNotifier
		ty := reflect.TypeOf(x).Elem()
		code := ifcterminal.GenerateInterfaceTerminalStruct(ty, "_X", "")
		typeDecl := []byte("type CloseNotifier_X struct")
		if !bytes.Contains(code, typeDecl) {
			t.Errorf("Expected %q in\n\n %s", string(typeDecl), string(code))
		}
		err := checkCode(code)
		if err != nil {
			t.Errorf("error in compiling generated code: %v \n\n%s", err, string(code))
		}
	}
	{
		var x *http.CookieJar
		ty := reflect.TypeOf(x).Elem()
		code := ifcterminal.GenerateInterfaceTerminalStruct(ty, "", "")
		err := checkCode(code)
		if err != nil {
			t.Errorf("error in compiling generated code: %v \n\n%s", err, string(code))
		}
	}
	{
		var x *ifcterminal.SomeFunctionalIfc
		ty := reflect.TypeOf(x).Elem()
		code := ifcterminal.GenerateInterfaceTerminalStruct(ty, "", "")
		err := checkCode(code)
		if err != nil {
			t.Errorf("error in compiling generated code: %v \n\n%s", err, string(code))
		}
	}
	{
		// Actually, recursive types don't seem to be a problem, because they are named...
		var x *ifcterminal.SomeRecursiveType
		ty := reflect.TypeOf(x).Elem()
		code := ifcterminal.GenerateInterfaceTerminalStruct(ty, "", "")
		err := checkCode(code)
		if err != nil {
			t.Errorf("error in compiling generated code: %v \n\n%s", err, string(code))
		}
	}
}

// Test if the generated code "seems legit"
func checkCode(code []byte) error {
	dir, err := ioutil.TempDir("", "ifcterminal-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	fpath := filepath.Join(dir, "test.go")
	err = ioutil.WriteFile(fpath, code, 0777)
	if err != nil {
		return err
	}
	command := exec.Command("go", "build", fpath)
	return command.Run()
}
