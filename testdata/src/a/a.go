package a

import (
	"io"
	"net/http"
)

// A is defined all fields.
var A = Value{
	F1: 1,
	F2: 2,
	F3: 3,
	F4: 4,
	F5: 5,
	F6: 6,
	f1: 1,
	f2: 2,
	f3: 3,
	f4: 4,
	f5: 5,
	f6: 6,
}

// B is defined all fields only values.
var B = Value{
	1, 2, 3, 4, 5, 6,
	1, 2, 3, 4, 5, 6,
}

// C is defined required field.
var C = Value{
	F1: 1,
	F2: 2,
	F3: 3,
	f1: 1,
	f2: 2,
	f3: 3,
}

// D is empty struct literal.
var D = Value{} // want "incomplete struct: Value requires F1, F2, F3, f1, f2, f3"

// E is partial defined fields.
var E = Value{ // want "incomplete struct: Value requires f1, f2, f3"
	F1: 1,
	F2: 2,
	F3: 3,
}

var ArrA = []*Value{
	{}, // cannot check field name omitted composit literal.
}

var ArrB = []string{"arr1"} // array literal

var X1 = io.PipeWriter{}

var X2 = http.Client{} // want "incomplete struct: http.Client requires Transport, CheckRedirect, Jar, Timeout"

// partial
var X3 = http.Client{} // if comments starts with `partial`, skip this
