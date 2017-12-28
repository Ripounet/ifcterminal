package ifcterminal_test

import (
	"io"
	"reflect"
	"testing"

	"github.com/Ripounet/ifcterminal"
)

func TestInterfaceTerminalStruct(t *testing.T) {
	{
		var x *io.ReadWriter
		ty := reflect.TypeOf(x).Elem()
		ifcterminal.GenerateInterfaceTerminalStruct(ty)
	}
	{
		var x *io.WriterTo
		ty := reflect.TypeOf(x).Elem()
		ifcterminal.GenerateInterfaceTerminalStruct(ty)
	}
}
