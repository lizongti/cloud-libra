package log

import "log"

type Logger interface {
	Tracef(format string, v ...interface{})
	Trace(v ...interface{})
	Traceln(v ...interface{})

	Debugf(format string, v ...interface{})
	Debug(v ...interface{})
	Debugln(v ...interface{})

	Infof(format string, v ...interface{})
	Info(v ...interface{})
	Infoln(v ...interface{})

	Printf(format string, v ...interface{})
	Print(v ...interface{})
	Println(v ...interface{})

	Warnf(format string, v ...interface{})
	Warn(v ...interface{})
	Warnln(v ...interface{})

	Warningf(format string, v ...interface{})
	Warning(v ...interface{})
	Warningln(v ...interface{})

	Errorf(format string, v ...interface{})
	Error(v ...interface{})
	Errorln(v ...interface{})

	Fatalf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalln(v ...interface{})

	Panicf(format string, v ...interface{})
	Panic(v ...interface{})
	Panicln(v ...interface{})
}

var logger Logger

// SetLogger must be called in init
func SetLogger(l Logger) {
	logger = l
}

func Printf(format string, v ...interface{}) {
	logger.Printf(format, v...)
}

func Print(v ...interface{}) {
	logger.Print(v...)
}

func Println(v ...interface{}) {
	logger.Println(v...)
}

func Fatalf(format string, v ...interface{}) {
	logger.Fatalf(format, v...)
}

func Fatal(v ...interface{}) {
	logger.Fatal(v...)
}

func Fatalln(v ...interface{}) {
	logger.Fatalln(v...)
}

func Panicf(format string, v ...interface{}) {
	logger.Panicf(format, v...)
}

func Panic(v ...interface{}) {
	logger.Panic(v...)
}

func Panicln(v ...interface{}) {
	logger.Panicln(v...)
}

func Tracef(format string, v ...interface{}) {
	logger.Tracef(format, v...)
}

func Trace(v ...interface{}) {
	logger.Trace(v...)
}

func Traceln(v ...interface{}) {
	logger.Traceln(v...)
}

func Debugf(format string, v ...interface{}) {
	logger.Debugf(format, v...)
}

func Debug(v ...interface{}) {
	logger.Debug(v...)
}

func Debugln(v ...interface{}) {
	logger.Debugln(v...)
}

func Infof(format string, v ...interface{}) {
	logger.Infof(format, v...)
}

func Info(v ...interface{}) {
	logger.Info(v...)
}

func Infoln(v ...interface{}) {
	logger.Infoln(v...)
}

func Warnf(format string, v ...interface{}) {
	logger.Warnf(format, v...)
}

func Warn(v ...interface{}) {
	logger.Warn(v...)
}

func Warnln(v ...interface{}) {
	logger.Warnln(v...)
}

func Warningf(format string, v ...interface{}) {
	logger.Warningf(format, v...)
}

func Warning(v ...interface{}) {
	logger.Warning(v...)
}

func Warningln(v ...interface{}) {
	logger.Warningln(v...)
}

func Errorf(format string, v ...interface{}) {
	logger.Errorf(format, v...)
}

func Error(v ...interface{}) {
	logger.Error(v...)
}

func Errorln(v ...interface{}) {
	logger.Errorln(v...)
}

type defaultLogger struct {
}

func (defaultLogger) Tracef(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (defaultLogger) Trace(v ...interface{}) {
	log.Print(v...)
}

func (defaultLogger) Traceln(v ...interface{}) {
	log.Println(v...)
}

func (defaultLogger) Debugf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (defaultLogger) Debug(v ...interface{}) {
	log.Print(v...)
}

func (defaultLogger) Debugln(v ...interface{}) {
	log.Println(v...)
}

func (defaultLogger) Infof(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (defaultLogger) Info(v ...interface{}) {
	log.Print(v...)
}

func (defaultLogger) Infoln(v ...interface{}) {
	log.Println(v...)
}

func (defaultLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (defaultLogger) Print(v ...interface{}) {
	log.Print(v...)
}

func (defaultLogger) Println(v ...interface{}) {
	log.Println(v...)
}

func (defaultLogger) Warnf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (defaultLogger) Warn(v ...interface{}) {
	log.Println(v...)
}

func (defaultLogger) Warnln(v ...interface{}) {
	log.Println(v...)
}

func (defaultLogger) Warningf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (defaultLogger) Warning(v ...interface{}) {
	log.Print(v...)
}

func (defaultLogger) Warningln(v ...interface{}) {
	log.Println(v...)
}

func (defaultLogger) Errorf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (defaultLogger) Error(v ...interface{}) {
	log.Print(v...)
}

func (defaultLogger) Errorln(v ...interface{}) {
	log.Println(v...)
}

func (defaultLogger) Fatalf(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}

func (defaultLogger) Fatal(v ...interface{}) {
	log.Fatal(v...)
}

func (defaultLogger) Fatalln(v ...interface{}) {
	log.Fatalln(v...)
}

func (defaultLogger) Panicf(format string, v ...interface{}) {
	log.Panicf(format, v...)
}

func (defaultLogger) Panic(v ...interface{}) {
	log.Panic(v...)
}

func (defaultLogger) Panicln(v ...interface{}) {
	log.Panicln(v...)
}
