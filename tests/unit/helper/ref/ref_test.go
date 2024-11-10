package ref

import (
	"fmt"
	"github.com/mandarine-io/Backend/internal/helper/ref"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func Test_RefUtils_SafeRef(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		datas := []any{
			"test",
			1,
			1.1,
			true,
			struct{ Name string }{Name: "test"},
			nil,
		}

		for _, data := range datas {
			t.Run(fmt.Sprintf("%T", data), func(t *testing.T) {
				ptr := ref.SafeRef(data)
				assert.Equal(t, data, *ptr)
			})
		}
	})
}

func Test_RefUtils_SafeDeref(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		str := "test"
		num := 1
		prec := 1.1
		boolean := true
		structure := struct{ Name string }{Name: "test"}

		t.Run(fmt.Sprintf("%T", &str), func(t *testing.T) {
			value := ref.SafeDeref(&str)
			assert.Equal(t, str, value)
		})

		t.Run(fmt.Sprintf("%T", &num), func(t *testing.T) {
			value := ref.SafeDeref(&num)
			assert.Equal(t, num, value)
		})

		t.Run(fmt.Sprintf("%T", &prec), func(t *testing.T) {
			value := ref.SafeDeref(&prec)
			assert.Equal(t, prec, value)
		})

		t.Run(fmt.Sprintf("%T", &boolean), func(t *testing.T) {
			value := ref.SafeDeref(&boolean)
			assert.Equal(t, boolean, value)
		})

		t.Run(fmt.Sprintf("%T", &structure), func(t *testing.T) {
			value := ref.SafeDeref(&structure)
			assert.Equal(t, structure, value)
		})

		t.Run(fmt.Sprintf("%T", nil), func(t *testing.T) {
			var nilPtr *string = nil

			value := ref.SafeDeref(nilPtr)
			assert.Equal(t, value, reflect.Zero(reflect.TypeOf(value)).Interface())
		})
	})
}
