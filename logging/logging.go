package logging

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	root    *zap.SugaredLogger
	logfile *os.File // not-nil if we opened a local file
	errfile *os.File // not-nil if we opened a local file
)

func Logger() *zap.SugaredLogger {
	return root
}

func NamedLogger(name string) *zap.SugaredLogger {
	if root == nil {
		panic("root logger was not setup")
	}
	return root.Named(name)
}

func Setup(debug bool, logfileName string, errfileName string) {
	if root != nil {
		return
	}

	if logfileName == "" {
		logfileName = "stdout"
	} else if errfileName == "" {
		errfileName = logfileName
	}
	if errfileName == "" {
		errfileName = "stderr"
	}

	cfg := zap.Config{
		Level:            zap.NewAtomicLevel(),
		Encoding:         "console",
		OutputPaths:      []string{logfileName},
		ErrorOutputPaths: []string{errfileName},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			NameKey:     "name",
			EncodeLevel: zapcore.LowercaseLevelEncoder,
		},
	}

	if debug {
		cfg.Level.SetLevel(zap.DebugLevel)
	} else {
		cfg.Level.SetLevel(zap.InfoLevel)
	}

	root = zap.Must(cfg.Build()).Sugar()
}

// alternate "setup" function used to suppress log output durinng test
func Disable() {
	root = zap.NewNop().Sugar()
}

func Teardown() {
	if root == nil {
		return
	}

	root.Sync()
	if logfile != nil {
		logfile.Close()
	}
	if errfile != nil {
		errfile.Close()
	}
}
