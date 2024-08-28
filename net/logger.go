package net

import "github.com/ameise84/logger"

var _gLogger logger.Log

func init() {
	_gLogger = logger.DefaultLogger()
}

func SetLogger(log logger.Logger) {
	_gLogger = log
}
