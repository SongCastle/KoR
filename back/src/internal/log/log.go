package log

import (
	"io"
	"log"
	"strconv"

	"github.com/SongCastle/KoR/internal/env"
)

const (
	DEBUG        = iota
	INFO
	WARN
	ERROR
	FATAL
	DEBUG_PREFIX = "[DEBUG] "
	INFO_PREFIX  = "[INFO] "
	WARN_PREFIX  = "[WARNING] "
	ERROR_PREFIX = "[ERROR] "
	FATAL_PREFIX = "[FATAL] "
	LOG_FLAG     = log.LstdFlags
)

var (
	writer io.Writer
	debug, info, warn, err, fatal *logger
	Level  uint
)

type any = interface{}

type logger struct {
	*log.Logger
	level uint
}

func (l *logger) Putf(format string, v ...any) {
	if Level <= l.level {
		l.Printf(format + "\n", v...)
	}
}

func new(writer io.Writer, prefix string, flag int, level uint) *logger {
	return &logger{Logger: log.New(writer, prefix, flag), level: level}
}

func init() {
	// 出力先 (標準出力)
	writer = log.Writer()
	// Logger
	debug  = new(writer, DEBUG_PREFIX, LOG_FLAG, DEBUG)
	info   = new(writer, INFO_PREFIX,  LOG_FLAG, INFO)
	warn   = new(writer, WARN_PREFIX,  LOG_FLAG, WARN)
	err    = new(writer, ERROR_PREFIX, LOG_FLAG, ERROR)
	fatal  = new(writer, FATAL_PREFIX, LOG_FLAG, FATAL)
	// ログレベル
	setLogLevel()
}

func setLogLevel() {
	k := "KOR_LOG_LEVEL"
	sl := env.Get(k)
	if sl != "" {
		l, err := strconv.ParseUint(sl, 10, 64)
		if err == nil {
			if DEBUG <= l && l <= FATAL {
				logf("Set log level: %d", l)
				Level = uint(l)
				return
			}
		}
		logf("Invalid log level (%s): %d", k, l)
	}
	logf("Set default log level: %d", DEBUG)
	Level = DEBUG
}

// used only in this package
func logf(format string, v ...any) {
	log.Printf("[LOG] " + format + "\n", v...)
}

func Debugf(format string, v ...any) {
	debug.Putf(format, v...)
}

func Infof(format string, v ...any) {
	info.Putf(format, v...)
}

func Warnf(format string, v ...any) {
	warn.Putf(format, v...)
}

func Errorf(format string, v ...any) {
	err.Putf(format, v...)
}

func Fatalf(format string, v ...any) {
	fatal.Putf(format, v...)
}
