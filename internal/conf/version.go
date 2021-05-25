package conf

import (
	"fmt"
)

var Major = 0
var Min = 1
var Rel = 1
var BVersion = fmt.Sprintf("%v.%v(rel %v)", Major, Min, Rel)
