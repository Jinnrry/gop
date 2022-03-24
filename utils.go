package gop

import (
	"fmt"
	"reflect"
	"strings"
	"time"
	"unsafe"
)

// GetPrivateField field value via field index
// TODO: we can use a LRU cache for the copy of the values, but it might be trivial for just testing.
func GetPrivateField(v reflect.Value, i int) reflect.Value {
	if v.Kind() != reflect.Struct {
		panic("expect v to be a struct")
	}

	copied := reflect.New(v.Type()).Elem()
	copied.Set(v)
	f := copied.Field(i)
	return reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
}

var float64Type = reflect.TypeOf(0.0)

// Compare returns the float value of x minus y
func Compare(x, y interface{}) float64 {
	if reflect.DeepEqual(x, y) {
		return 0
	}

	if x != nil && y != nil {
		xVal := reflect.Indirect(reflect.ValueOf(x))
		yVal := reflect.Indirect(reflect.ValueOf(y))

		if xVal.CanConvert(float64Type) && yVal.CanConvert(float64Type) {
			return xVal.Convert(float64Type).Float() - yVal.Convert(float64Type).Float()
		}

		if xt, ok := xVal.Interface().(time.Time); ok {
			if yt, ok := yVal.Interface().(time.Time); ok {
				return float64(xt.Sub(yt))
			}
		}
	}

	sa := fmt.Sprintf("%#v", x)
	sb := fmt.Sprintf("%#v", y)

	return float64(strings.Compare(sa, sb))
}
