package lowlevel

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"runtime"

// 	"slots-nano/internal/dao/cache/logger/hook"

// 	"github.com/sirupsen/logrus"
// )

// // NewRunLogger creates a logger hooked by lumberjack.Logger,
// // which produce logs that are stored in aliyun's sls production.
// // It is for common use, such as "log.Println", and shows in stderr.
// func NewRunLogger(c map[string]interface{}) *Logger {
// 	logger := NewLogger()
// 	logger.SetReportCaller(true)
// 	logger.SetFormatter(&logrus.TextFormatter{
// 		ForceColors:            true,
// 		TimestampFormat:        "2006/01/02 15:04:05.0000000", // the "time" field configuratiom
// 		FullTimestamp:          true,
// 		DisableLevelTruncation: true, // log upgrade field configuration
// 		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
// 			return "", fmt.Sprintf(" %s:%d:", GetPackageFile(f.File), f.Line)
// 		},
// 	})
// 	logger.SetOutput(ioutil.Discard)
// 	if err := logger.ReadLevel(Run, c); err != nil { // Default is Info
// 		panic(err)
// 	}
// 	if err := logger.ReadHooks(Run, c, map[string]*hook.Processor{}); err != nil {
// 		panic(err)
// 	}
// 	return logger
// }
