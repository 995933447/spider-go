package logger

type LogWriter interface {
	Emergency(message interface{}, context interface{}) error
	Alert(message interface{}, context interface{}) error
	Info(message interface{}, context interface{}) error
	Debug(message interface{}, context interface{}) error
	Error(message interface{}, context interface{}) error
	Critical(message interface{}, context interface{}) error
	Warning(message interface{}, context interface{})error
	Notice(message interface{}, context interface{}) error
}

type Logger struct {
	LogWriterList []LogWriter
}

func (logger *Logger) RegisterWriter(w LogWriter) {
	logger.LogWriterList = append(logger.LogWriterList, w)
}

func (logger *Logger) Emergency(message interface{}, context interface{}) {
	for _, writer := range logger.LogWriterList {
		writer.Emergency(message, context)
	}
}

func (logger *Logger) Alert(message interface{}, context interface{}) {
	for _, writer := range logger.LogWriterList {
		writer.Alert(message, context)
	}
}

func (logger *Logger) Info(message interface{}, context interface{}) {
	for _, writer := range logger.LogWriterList {
		writer.Info(message, context)
	}
}

func (logger *Logger) Debug(message interface{}, context interface{}) {
	for _, writer := range logger.LogWriterList {
		writer.Debug(message, context)
	}
}

func (logger *Logger) Error(message interface{}, context interface{}) {
	for _, writer := range logger.LogWriterList {
		writer.Error(message, context)
	}
}

func (logger *Logger) Critical(message interface{}, context interface{}) {
	for _, writer := range logger.LogWriterList {
		writer.Critical(message, context)
	}
}

func (logger *Logger) Warning(message interface{}, context interface{}) {
	for _, writer := range logger.LogWriterList {
		writer.Warning(message, context)
	}
}

func (logger *Logger) Notice(message interface{}, context interface{}) {
	for _, writer := range logger.LogWriterList {
		writer.Notice(message, context)
	}
}