package log

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
