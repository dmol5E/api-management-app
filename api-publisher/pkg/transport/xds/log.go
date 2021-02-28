package xds

import (
	logs "log"

	"github.com/envoyproxy/go-control-plane/pkg/log"
)

var logger = log.LoggerFuncs{
	DebugFunc: func(s string, i ...interface{}) {
		logs.Printf("[DEBUG]"+s, i...)
	},
	InfoFunc: func(s string, i ...interface{}) {
		logs.Printf("[INFO]"+s, i...)
	},
	ErrorFunc: func(s string, i ...interface{}) {
		logs.Printf("[ERROR]"+s, i...)
	},
	WarnFunc: func(s string, i ...interface{}) {
		logs.Printf("[WARN]"+s, i...)
	},
}
