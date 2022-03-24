package gop

import (
	"fmt"
	"testing"
	"time"
	"unsafe"
)

func TestF(t *testing.T) {
	ref := "test"
	timeStamp, _ := time.Parse(time.RFC3339Nano, "2021-08-28T08:36:36.807908+08:00")
	v := []interface{}{
		nil,
		[]interface{}{true, false, uintptr(0x17), float32(100.121111133)},
		true, 10, int8(2), int32(100),
		float64(100.121111133),
		complex64(1 + 2i), complex128(1 + 2i),
		[3]int{1, 2},
		make(chan int),
		make(chan string, 3),
		func(string) int { return 10 },
		map[interface{}]interface{}{
			"test": 10,
			"a":    1,
		},
		unsafe.Pointer(&ref),
		struct {
			Int int
			str string
			M   map[int]int
		}{10, "ok", map[int]int{1: 0x20}},
		[]byte("aa\xe2"),
		[]byte("bytes\n\tbytes"),
		byte('a'),
		byte(1),
		'天',
		"\ntest",
		&ref,
		(*struct{ Int int })(nil),
		&struct{ Int int }{},
		&map[int]int{1: 2, 3: 4},
		&[]int{1, 2},
		&[2]int{1, 2},
		&[]byte{1, 2},
		timeStamp,
		time.Hour,
	}

	out := Sprint(v)

	expected := `[]interface {}/* len=31 cap=31 */{
    nil,
    []interface {}/* len=4 cap=4 */{
        true,
        false,
        uintptr(23),
        float32(100.12111),
    },
    true,
    10,
    int8(2),
    'd',
    float64(100.121111133),
    complex64(1+2i),
    1+2i,
    [3]int{
        1,
        2,
        0,
    },
    make(chan int),
    make(chan string, 3),
    (func(string) int)(nil),
    map[interface {}]interface {}/* len=2 */{
        "a": 1,
        "test": 10,
    },
    unsafe.Pointer(uintptr(` + fmt.Sprintf("%v", &ref) + `)),
    struct { Int int; str string; M map[int]int }/* len=3 */{
        Int: 10,
        str: "ok",
        M: map[int]int{
            1: 32,
        },
    },
    Base64("YWHi")/* len=3 */,
    []byte("" +
        "bytes\n" +
        "\tbytes")/* len=12 */,
    byte('a'),
    byte(0x1),
    '天',
    "" +
        "\n" +
        "test"/* len=5 */,
    Ptr("test").(*string),
    (*struct { Int int })(nil),
    &struct { Int int }{
        Int: 0,
    },
    &map[int]int/* len=2 */{
        1: 2,
        3: 4,
    },
    &[]int/* len=2 cap=2 */{
        1,
        2,
    },
    &[2]int{
        1,
        2,
    },
    Ptr([]byte("\x01\x02")/* len=2 */).(*[]uint8),
    Time("2021-08-28T08:36:36.807908+08:00"),
    Time.Duration("1h0m0s"),
}`

	if out != expected {
		t.Error("格式化失败！")
	}

	Print(v)
}
