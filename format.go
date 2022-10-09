package gop

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// stdout is the default stdout for gop.P .
var stdout io.Writer = os.Stdout

const indentUnit = "    "

// defaultTheme colors for Sprint
var defaultTheme = func(t Type) Color {
	switch t {
	case TypeName:
		return Cyan
	case Bool:
		return Blue
	case Rune, Byte, String:
		return Yellow
	case Number:
		return Green
	case Chan, Func, UnsafePointer:
		return Magenta
	case Len:
		return White
	case PointerCircular:
		return Red
	default:
		return None
	}
}

// noTheme colors for Sprint
var noTheme = func(t Type) Color {
	return None
}

// SprintWithColor  is a shortcut for format with color
// 返回结构体信息，带颜色
func SprintWithColor(v interface{}) string {
	return format(Tokenize(v), nil)
}

// Print pretty print the value list vs
func Print(vs ...interface{}) {
	list := []interface{}{}
	for _, v := range vs {
		list = append(list, format(Tokenize(v), nil))
	}
	_, _ = fmt.Fprintln(stdout, list...)
}

// PrintWithColor 输出带颜色的文字
func PrintWithColor(theme Color, vs ...interface{}) {
	list := []interface{}{}
	for _, v := range vs {
		list = append(list, format(Tokenize(v), func(t Type) Color {
			return theme
		}))
	}
	_, _ = fmt.Fprintln(stdout, list...)
}

// Sprint /* Sprint is a shortcut for format with plain color
func Sprint(v interface{}) string {
	return format(Tokenize(v), noTheme)
}

// format a list of tokens
func format(ts []*Token, theme func(Type) Color) string {
	if theme == nil {
		theme = defaultTheme
	}

	out := ""
	depth := 0
	for i, t := range ts {
		if oneOf(t.Type, SliceOpen, MapOpen, StructOpen) {
			depth++
		}
		if i < len(ts)-1 && oneOf(ts[i+1].Type, SliceClose, MapClose, StructClose) {
			depth--
		}

		color := theme(t.Type)
		s := ColorStr(color, t.Literal)

		switch t.Type {
		case SliceOpen, MapOpen, StructOpen:
			out += s + "\n"
		case SliceItem, MapKey, StructKey:
			out += strings.Repeat(indentUnit, depth)
		case Colon, InlineComma:
			out += s + " "
		case Comma:
			out += s + "\n"
		case SliceClose, MapClose, StructClose:
			out += strings.Repeat(indentUnit, depth) + s
		case String:
			out += ColorStr(color, readableStr(depth, t.Literal))
		default:
			out += s
		}
	}

	return out
}

func oneOf(t Type, list ...Type) bool {
	for _, el := range list {
		if t == el {
			return true
		}
	}
	return false
}

// To make line string block more human readable.
// Split newline into two strings, convert "\t" into tab.
// Such as foramt string: "line one \n\t line two" into:
//     "line one \n" +
//     "	 line two"
func readableStr(depth int, s string) string {
	//s, _ = replaceEscaped(s, '\t', "	")
	//
	//indent := strings.Repeat(indentUnit, depth+1)
	//if n, has := replaceEscaped(s, 'n', "\\n\" +\n"+indent+"\""); has {
	//	return "\"\" +\n" + indent + n
	//}

	//s = strings.ReplaceAll(s, "\\n", "\n")

	return s
}

// We use a simple state machine to replace escaped char like "\n"
func replaceEscaped(s string, escaped rune, new string) (string, bool) {
	type State int
	const (
		init State = iota
		prematch
		match
	)

	state := init
	out := ""
	buf := ""
	has := false

	onInit := func(r rune) {
		state = init
		out += buf + string(r)
		buf = ""
	}

	onPrematch := func() {
		state = prematch
		buf = "\\"
	}

	onEscape := func() {
		state = match
		out += new
		buf = ""
		has = true
	}

	for _, r := range s {
		switch state {
		case prematch:
			switch r {
			case escaped:
				onEscape()
			default:
				onInit(r)
			}

		case match:
			switch r {
			case '\\':
				onPrematch()
			default:
				onInit(r)
			}

		default:
			switch r {
			case '\\':
				onPrematch()
			default:
				onInit(r)
			}
		}
	}

	return out, has
}
