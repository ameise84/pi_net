package net

import "github.com/ameise84/logger"

var _gLogger logger.Logger

func init() {
	_gLogger = logger.DefaultLogger()
}

func SetLogger(log logger.Logger) {
	_gLogger = log
}
