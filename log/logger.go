package log

import "log"

// DefaultLogger represents the default log interface.
type DefaultLogger interface {
	Printf(format string, v ...interface{})
	Print(v ...interface{})
	Println(v ...interface{})

	Fatalf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalln(v ...interface{})

	Panicf(format string, v ...interface{})
	Panic(v ...interface{})
	Panicln(v ...interface{})
}

var defaultLogger DefaultLogger = log.Default()

// Logger represents the standard log interface.
type Logger interface {
	DefaultLogger

	Tracef(format string, v ...interface{})
	Trace(v ...interface{})
	Traceln(v ...interface{})

	Debugf(format string, v ...interface{})
	Debug(v ...interface{})
	Debugln(v ...interface{})

	Infof(format string, v ...interface{})
	Info(v ...interface{})
	Infoln(v ...interface{})

	Warnf(format string, v ...interface{})
	Warn(v ...interface{})
	Warnln(v ...interface{})

	Warningf(format string, v ...interface{})
	Warning(v ...interface{})
	Warningln(v ...interface{})

	Errorf(format string, v ...interface{})
	Error(v ...interface{})
	Errorln(v ...interface{})
}

type LocalLogger struct {
	Printf    func(format string, v ...interface{})
	Print     func(v ...interface{})
	Println   func(v ...interface{})
	Fatalf    func(format string, v ...interface{})
	Fatal     func(v ...interface{})
	Fatalln   func(v ...interface{})
	Panicf    func(format string, v ...interface{})
	Panic     func(v ...interface{})
	Panicln   func(v ...interface{})
	Tracef    func(format string, v ...interface{})
	Trace     func(v ...interface{})
	Traceln   func(v ...interface{})
	Debugf    func(format string, v ...interface{})
	Debug     func(v ...interface{})
	Debugln   func(v ...interface{})
	Infof     func(format string, v ...interface{})
	Info      func(v ...interface{})
	Infoln    func(v ...interface{})
	Warnf     func(format string, v ...interface{})
	Warn      func(v ...interface{})
	Warnln    func(v ...interface{})
	Warningf  func(format string, v ...interface{})
	Warning   func(v ...interface{})
	Warningln func(v ...interface{})
	Errorf    func(format string, v ...interface{})
	Error     func(v ...interface{})
	Errorln   func(v ...interface{})
}

func Set(logger Logger) {
	localLogger.Printf = logger.Printf
	localLogger.Print = logger.Print
	localLogger.Println = logger.Println
	localLogger.Fatal = logger.Fatal
	localLogger.Fatalln = logger.Fatalln
	localLogger.Panicf = logger.Panicf
	localLogger.Panic = logger.Panic
	localLogger.Panicln = logger.Panicln
	localLogger.Tracef = logger.Tracef
	localLogger.Trace = logger.Trace
	localLogger.Traceln = logger.Traceln
	localLogger.Debugf = logger.Debugf
	localLogger.Debug = logger.Debug
	localLogger.Debugln = logger.Debugln
	localLogger.Infof = logger.Infof
	localLogger.Info = logger.Info
	localLogger.Infoln = logger.Infoln
	localLogger.Warnf = logger.Warnf
	localLogger.Warn = logger.Warn
	localLogger.Warnln = logger.Warnln
	localLogger.Warningf = logger.Warningf
	localLogger.Warning = logger.Warning
	localLogger.Warningln = logger.Warningln
	localLogger.Errorf = logger.Errorf
	localLogger.Error = logger.Error
	localLogger.Errorln = logger.Errorln
	localLogger.Fatalf = logger.Fatalf
}

var localLogger = &LocalLogger{
	Printf:    defaultLogger.Printf,
	Print:     defaultLogger.Print,
	Println:   defaultLogger.Println,
	Fatalf:    defaultLogger.Fatalf,
	Fatal:     defaultLogger.Fatal,
	Fatalln:   defaultLogger.Fatalln,
	Panicf:    defaultLogger.Panicf,
	Panic:     defaultLogger.Panic,
	Panicln:   defaultLogger.Panicln,
	Tracef:    defaultLogger.Printf,
	Trace:     defaultLogger.Print,
	Traceln:   defaultLogger.Println,
	Debugf:    defaultLogger.Printf,
	Debug:     defaultLogger.Print,
	Debugln:   defaultLogger.Println,
	Infof:     defaultLogger.Printf,
	Info:      defaultLogger.Print,
	Infoln:    defaultLogger.Println,
	Warnf:     defaultLogger.Printf,
	Warn:      defaultLogger.Print,
	Warnln:    defaultLogger.Println,
	Warningf:  defaultLogger.Printf,
	Warning:   defaultLogger.Print,
	Warningln: defaultLogger.Println,
	Errorf:    defaultLogger.Printf,
	Error:     defaultLogger.Print,
	Errorln:   defaultLogger.Println,
}

func Printf(format string, v ...interface{}) {
	localLogger.Printf(format, v...)
}

func Print(v ...interface{}) {
	localLogger.Print(v...)
}

func Println(v ...interface{}) {
	localLogger.Println(v...)
}

func Fatalf(format string, v ...interface{}) {
	localLogger.Fatalf(format, v...)
}

func Fatal(v ...interface{}) {
	localLogger.Fatal(v...)
}

func Fatalln(v ...interface{}) {
	localLogger.Fatalln(v...)
}

func Panicf(format string, v ...interface{}) {
	localLogger.Panicf(format, v...)
}

func Panic(v ...interface{}) {
	localLogger.Panic(v...)
}

func Panicln(v ...interface{}) {
	localLogger.Panicln(v...)
}

func Tracef(format string, v ...interface{}) {
	localLogger.Tracef(format, v...)
}

func Trace(v ...interface{}) {
	localLogger.Trace(v...)
}

func Traceln(v ...interface{}) {
	localLogger.Traceln(v...)
}

func Debugf(format string, v ...interface{}) {
	localLogger.Debugf(format, v...)
}

func Debug(v ...interface{}) {
	localLogger.Debug(v...)
}

func Debugln(v ...interface{}) {
	localLogger.Debugln(v...)
}

func Infof(format string, v ...interface{}) {
	localLogger.Infof(format, v...)
}

func Info(v ...interface{}) {
	localLogger.Info(v...)
}

func Infoln(v ...interface{}) {
	localLogger.Infoln(v...)
}

func Warnf(format string, v ...interface{}) {
	localLogger.Warnf(format, v...)
}

func Warn(v ...interface{}) {
	localLogger.Warn(v...)
}

func Warnln(v ...interface{}) {
	localLogger.Warnln(v...)
}

func Warningf(format string, v ...interface{}) {
	localLogger.Warningf(format, v...)
}

func Warning(v ...interface{}) {
	localLogger.Warning(v...)
}

func Warningln(v ...interface{}) {
	localLogger.Warningln(v...)
}

func Errorf(format string, v ...interface{}) {
	localLogger.Errorf(format, v...)
}

func Error(v ...interface{}) {
	localLogger.Error(v...)
}

func Errorln(v ...interface{}) {
	localLogger.Errorln(v...)
}
