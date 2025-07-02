package logger

import "gitag.ir/armogroup/armo/services/reality/modules/logger/ports"

type LoggerEngine[T ports.LoggerEngineType] struct {
	Types     []string
	Instances *[]T
}
