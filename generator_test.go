package ifcterminal_test

import (
	"io"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/Ripounet/ifcterminal"
)

func TestInterfaceTerminalStruct(t *testing.T) {
	{
		var x *io.Reader
		ty := reflect.TypeOf(x).Elem()
		code := ifcterminal.GenerateInterfaceTerminalStruct(ty)
		os.Stderr.Write(code)
	}
	{
		var x *io.ReadWriter
		ty := reflect.TypeOf(x).Elem()
		code := ifcterminal.GenerateInterfaceTerminalStruct(ty)
		os.Stderr.Write(code)
	}
	{
		var x *io.WriterTo
		ty := reflect.TypeOf(x).Elem()
		code := ifcterminal.GenerateInterfaceTerminalStruct(ty)
		os.Stderr.Write(code)
	}
	{
		var x *http.CloseNotifier
		ty := reflect.TypeOf(x).Elem()
		code := ifcterminal.GenerateInterfaceTerminalStruct(ty)
		os.Stderr.Write(code)
	}
	{
		var x *http.CookieJar
		ty := reflect.TypeOf(x).Elem()
		code := ifcterminal.GenerateInterfaceTerminalStruct(ty)
		os.Stderr.Write(code)
	}
}
