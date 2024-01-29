package logger

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	TimeLayout = "01-02 15:04:05"
)

var (
	once sync.Once
	sl   *zap.SugaredLogger

	sl2 *zap.SugaredLogger
)

type Interface interface {
	Printf(string, ...any)
	Println(...any)
}

type loggerFuncs struct {
	printfFunc  func(string, ...any)
	printlnFunc func(...any)
}

var _ Interface = (*loggerFuncs)(nil)

func (l *loggerFuncs) Printf(format string, args ...any) {
	if l.printfFunc != nil {
		l.printfFunc(format, args...)
	}
}
func (l *loggerFuncs) Println(args ...any) {
	if l.printlnFunc != nil {
		l.printlnFunc(args...)
	}
}

func InitLogger() {
	once.Do(func() {
		config := zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(TimeLayout)
		l, _ := config.Build()
		sl = l.Sugar()

		// adapter包了一层，打印文件位置时往上查找
		l2, _ := config.Build(zap.AddCallerSkip(2))
		sl2 = l2.Sugar()
	})
}

func Sugar() *zap.SugaredLogger {
	if sl != nil {
		return sl
	}
	InitLogger()
	return sl
}

func Flush() {
	if sl != nil {
		sl.Sync()
	}
	if sl2 != nil {
		sl2.Sync()
	}
}

func GetInfoLogger() Interface {
	return &loggerFuncs{
		printfFunc: func(format string, args ...any) {
			sl2.Infof(format, args...)
		},
		printlnFunc: func(args ...any) {
			sl2.Infoln(args...)
		},
	}
}

func GetWarnLogger() Interface {
	return &loggerFuncs{
		printfFunc: func(format string, args ...any) {
			sl2.Warnf(format, args...)
		},
		printlnFunc: func(args ...any) {
			sl2.Warnln(args...)
		},
	}
}
