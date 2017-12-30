package ifcterminal_test

import (
	"io"
	"net/http"
	"reflect"
	"testing"
	"os"

	"github.com/Ripounet/ifcterminal"
)

func TestInterfaceTerminalStruct(t *testing.T) {
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
}
