package ports

import "gitag.ir/armogroup/armo/services/reality/modules/logger/engines"

type LoggerEngineType interface {
	*engines.LoggerEngineFile | *engines.LoggerEngineStdout | any
}
