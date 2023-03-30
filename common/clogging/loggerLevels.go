package clogging

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"

	"go.uber.org/zap/zapcore"
)

// loggerNameRegexp 匹配形如 "xxx.xxx" 或 "xxx" 这样的字符串，其中，"xxx" 可以是字母，可以是数字，可以是下划线，可以是井号。可以是冒号，也可以是减号，其他的都不行。
var loggerNameRegexp = regexp.MustCompile(`^[[:alnum:]_#:-]+(\.[[:alnum:]_#:-]+)*$`)

type LoggerLevels struct {
	mutex sync.RWMutex
	levelCache map[string]zapcore.Level // 一个动态的缓冲区，将程序运行过程中遇到的日志记录器和对应的日志记录存储在这个缓冲区里。
	specs map[string]zapcore.Level
	defaultLevel zapcore.Level
	minLevel zapcore.Level
}

// DefaultLevel 为没有明确设置日志级别的记录器返回默认日志级别。
func (l *LoggerLevels) DefaultLevel() zapcore.Level {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.defaultLevel
}

// ActivateSpec 处理的字符串形式为："xx.yy.zz=level1:level2:aa.bb=level3"，其中等于号 `=` 的两边分别表示日志记录器名和日志级别，
// 等式表示为日志记录器设置对应的日志级别，`xx.yy.xx` 里的点号 `.` 可以看作是日志记录器的取名方式，类比网址，层层递进的那种关系，可
// 以将 `xx.yy` 看成是 `xx` 的子日志记录器的名称。`level2` 没有与之对应的日志记录器，那么 `level2` 会被看作是默认的日志级别，会赋
// 值给 LoggerLevel 的 defaultLevel 字段。`xx.yy.zz=level1` 与 `aa.bb=level3` 会被 LoggerLevel 的 specs 字段存储，其中 `xx.yy.zz`
// 与 `aa.bb` 会被当成 map 的 key，`level1` 和 `level3` 则会被当成对应的 value。`level1 level2 level3` 里级别最小的会被赋值给
// LoggerLevel 的 minLevel 字段。
func (l *LoggerLevels) ActivateSpec(spec string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	defaultLevel := zapcore.InfoLevel
	specs := map[string]zapcore.Level{}
	for _, field := range strings.Split(spec, ":") {
		split := strings.Split(field, "=")
		switch len(split) {
		case 1:
			if field != "" && !IsValidLevel(field) {
				return fmt.Errorf("invalid logging specification '%s': bad segment '%s'", spec, field)
			}
			defaultLevel = NameToLevel(field)
		case 2:
			if split[0] == "" {
				return fmt.Errorf("invalid logging specification '%s': no logger specified in segment '%s'", spec, field)
			}
			if field != "" && !IsValidLevel(split[1]) {
				return fmt.Errorf("invalid logging specification '%s': bad segment '%s'", spec, field)
			}
			level := NameToLevel(split[1])
			loggers := strings.Split(split[0], ",")
			for _, logger := range loggers {
				if !isValidLoggerName(strings.TrimSuffix(logger, ".")) {
					return fmt.Errorf("invalid logging specification '%s': bad logger name '%s'", spec, logger)
				}
				specs[logger] = level
			}
		default:
			return fmt.Errorf("invalid logging specification '%s': bad segment '%s'", spec, field)
		}
	}

	minLevel := defaultLevel
	for _, lvl := range specs {
		if lvl < minLevel {
			minLevel = lvl
		}
	}

	l.minLevel = minLevel
	l.defaultLevel = defaultLevel
	l.specs = specs
	l.levelCache = map[string]zapcore.Level{}
	return nil
}

func (l *LoggerLevels) Level(loggerName string) zapcore.Level {
	if level, ok := l.cachedLevel(loggerName); ok {
		return level
	}

	l.mutex.Lock()
	level := l.calculateLevel(loggerName)
	l.levelCache[loggerName] = level
	l.mutex.Unlock()

	return level
}

func (l *LoggerLevels) Spec() string {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	var fields []string
	for k, v := range l.specs {
		fields = append(fields, fmt.Sprintf("%s=%s", k, v))
	}

	sort.Strings(fields)
	fields = append(fields, l.defaultLevel.String())
	return strings.Join(fields, ":") // 从这里可以看出，spec 的形式是这样的 "logger.A,logger.B=info:logger.C=debug"
}

func (l *LoggerLevels) Enabled(lvl zapcore.Level) bool {
	l.mutex.RLock()
	enabled := l.minLevel.Enabled(lvl)
	l.mutex.RUnlock()
	return enabled
}

func (l *LoggerLevels) cachedLevel(loggerName string) (zapcore.Level, bool) {
	l.mutex.RLock()
	level, ok := l.levelCache[loggerName]
	l.mutex.RUnlock()
	return level, ok
}

func (l *LoggerLevels) calculateLevel(loggerName string) zapcore.Level {
	candidate := loggerName + "."
	for {
		if lvl, ok := l.specs[candidate]; ok {
			return lvl
		}

		idx := strings.LastIndex(candidate, ".")
		if idx <= 0 {
			return l.defaultLevel
		}
		candidate = candidate[:idx]
	}
}

func isValidLoggerName(loggerName string) bool {
	return loggerNameRegexp.MatchString(loggerName)
}