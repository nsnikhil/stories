package reporters

import (
	"github.com/nsnikhil/stories/pkg/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewExternalLogFile(cfg config.LogFileConfig) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   cfg.GetFilePath(),
		MaxSize:    cfg.GetFileMaxSizeInMb(),
		MaxBackups: cfg.GetFileMaxBackups(),
		MaxAge:     cfg.GetFileMaxAge(),
		LocalTime:  cfg.GetFileWithLocalTimeStamp(),
	}
}
