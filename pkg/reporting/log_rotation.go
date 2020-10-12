package reporters

import (
	"github.com/nsnikhil/stories/pkg/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

//TODO: HOW DO YOU DEAL WITH LOG ROTATION ?
func NewExternalLogFile(cfg config.LogFileConfig) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   cfg.GetFilePath(),
		MaxSize:    cfg.GetFileMaxSizeInMb(),
		MaxBackups: cfg.GetFileMaxBackups(),
		MaxAge:     cfg.GetFileMaxAge(),
		LocalTime:  cfg.GetFileWithLocalTimeStamp(),
	}
}
