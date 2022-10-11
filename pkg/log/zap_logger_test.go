package log

var defaultTestConfig = &OutputConfig{

	Writer:    OutputConsole,
	Level:     "debug",
	Formatter: EncoderTypeJson,
}
