package log

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	crzap "sigs.k8s.io/controller-runtime/pkg/log/zap"
)

type LevelEnablerByName struct {
	zapcore.Core
	MinLevels    map[string]zapcore.Level
	DefaultLevel zapcore.Level
	MinLevel     zapcore.Level
}

func (c LevelEnablerByName) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	loggerName := entry.LoggerName

	minLevel := c.DefaultLevel

	for prefix, level := range c.MinLevels {
		if loggerName == prefix || (len(loggerName) > len(prefix) && loggerName[:len(prefix)] == prefix && loggerName[len(prefix)] == '.') {
			minLevel = level
			break
		}
	}

	if entry.Level >= minLevel {
		return ce.AddCore(entry, c)
	}

	return ce
}

func (c LevelEnablerByName) Enabled(level zapcore.Level) bool {
	return level >= c.MinLevel
}

func (c LevelEnablerByName) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	return c.Core.Write(ent, fields)
}

func New(opts *crzap.Options, levelMap map[string]string) (logr.Logger, error) {
	zapLog := crzap.NewRaw(crzap.UseFlagOptions(opts))

	zapLevels, err := toZapLevels(levelMap)

	if err != nil {
		return logr.Discard(), err
	}

	defaultLevel := zapLog.Level()
	minLevel := defaultLevel
	for _, l := range zapLevels {
		if l < minLevel {
			minLevel = l
		}
	}

	customCore := LevelEnablerByName{
		Core:         zapLog.Core(),
		MinLevels:    zapLevels,
		DefaultLevel: defaultLevel,
		MinLevel:     minLevel,
	}

	return zapr.NewLogger(zap.New(customCore)), nil
}

func toZapLevels(levelMap map[string]string) (map[string]zapcore.Level, error) {
	zapLevels := make(map[string]zapcore.Level)

	for k, v := range levelMap {
		l, err := toZapLevel(v)
		if err != nil {
			return nil, err
		}
		zapLevels[k] = l
	}

	return zapLevels, nil
}

var levelStrings = map[string]zapcore.Level{
	"debug":  zap.DebugLevel,
	"info":   zap.InfoLevel,
	"warn":   zap.WarnLevel,
	"error":  zap.ErrorLevel,
	"dpanic": zap.DPanicLevel,
	"panic":  zap.PanicLevel,
	"fatal":  zap.FatalLevel,
}

func toZapLevel(flagValue string) (zapcore.Level, error) {
	level, validLevel := levelStrings[strings.ToLower(flagValue)]
	if !validLevel {
		logLevel, err := strconv.Atoi(flagValue)
		if err != nil {
			return zap.FatalLevel, fmt.Errorf("invalid log level \"%s\"", flagValue)
		}
		if logLevel > 0 {
			intLevel := -1 * logLevel
			return zapcore.Level(int8(intLevel)), nil
		} else {
			return zap.FatalLevel, fmt.Errorf("invalid log level \"%s\"", flagValue)
		}
	} else {
		return level, nil
	}
}
