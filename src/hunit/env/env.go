package env

import (
	"os"
)

var ExprDebug = os.Getenv("HUNIT_EXPR_DEBUG") != ""
