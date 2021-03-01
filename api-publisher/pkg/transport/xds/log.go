package xds

import (
	logs "github.com/sirupsen/logrus"

	"github.com/envoyproxy/go-control-plane/pkg/log"
)

var logger = log.LoggerFuncs{
	DebugFunc: func(s string, i ...interface{}) {
		logs.Debugf(s, i...)
	},
	InfoFunc: func(s string, i ...interface{}) {
		logs.Infof(s, i...)
	},
	ErrorFunc: func(s string, i ...interface{}) {
		logs.Errorf(s, i...)
	},
	WarnFunc: func(s string, i ...interface{}) {
		logs.Warnf(s, i...)
	},
}
