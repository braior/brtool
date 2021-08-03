package brtool

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
)

var (
	logrusLogger *log.Logger
)

// BRLogger BRLogger
type BRLogger struct{}

// logrusEntry 返回logrusEntry
func (logger *BRLogger) logrusEntry(commonfFileds map[string]interface{}) *log.Entry {
	return logrusLogger.WithFields(log.Fields(commonfFileds))
}

// Debug Debug日志
func (logger *BRLogger) Debug(commonfFileds map[string]interface{}, message string) {
	logger.logrusEntry(commonfFileds).Debugf("%s", message)
}

// Info Inf日志
func (logger *BRLogger) Info(commonfFileds map[string]interface{}, message string) {
	logger.logrusEntry(commonfFileds).Infof("%s", message)
}

// Warn Warn日志
func (logger *BRLogger) Warn(commonFields map[string]interface{}, message string) {
	logger.logrusEntry(commonFields).Warnf("%s", message)
}

// Error Error日志
func (logger *BRLogger) Error(commonFields map[string]interface{}, message string) {
	logger.logrusEntry(commonFields).Errorf("%s", message)
}

// Fatal Fatal日志
func (logger *BRLogger) Fatal(commonFields map[string]interface{}, message string) {
	logger.logrusEntry(commonFields).Fatalf("%s", message)
}

// Panic Panic日志
func (logger *BRLogger) Panic(commonFields map[string]interface{}, message string) {
	logger.logrusEntry(commonFields).Panicf("%s", message)
}

// BRLog 定义
type BRLog struct {
	// log 路径
	LogPath string

	// 日志类型 josn|text default: json
	LogType string

	// 日志文件日期格式 default: %Y-%m-%d|%Y%m%d
	FileNameDateFormat string

	// 日志中日期时间格式 default: 2006-01-02 15:04:05
	TimestampFormat string

	// 是否分离不同级别的日志 default: true
	IsSeparateLevelLog bool

	// 日志级别 默认: log.InfoLevel
	LogLevel log.Level

	// 日志最长保存多久 default: 15天
	MaxAge time.Duration

	// 日志默认多长时间轮转一次 默认: 24小时
	RotationTime time.Duration
}

// NewBRLog 返回BRLog对象和目录创建失败error
func NewBRLog(logPath string) (*BRLog, error) {
	if logPath == "" {
		return nil, errors.New("the log path must be provided")
	}
	logDir := path.Dir(logPath)
	_, err := os.Stat(logDir)
	if !(err == nil || os.IsExist(err)) {
		err := os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("create log dir <%s> error: %s", logDir, err)
		}
	}

	return &BRLog{
		LogPath:            logPath,
		LogType:            "json",
		LogLevel:           log.InfoLevel,
		FileNameDateFormat: "%Y-%m-%d",
		TimestampFormat:    "2006-01-02 15:04:05",
		IsSeparateLevelLog: true,
		MaxAge:             15 * 24 * time.Hour,
		RotationTime:       24 * time.Hour,
	}, nil
}

// SetLogType 设置日志格式 json|text
func (b *BRLog) SetLogType(logType string) {
	b.LogType = logType
}

// SetMaxAge 设置最大保留时间 单位：天
func (b *BRLog) SetMaxAge(day time.Duration) {
	b.MaxAge = day * 24 * time.Hour
}

// SetRotationTime 设置日志多久轮转一次
// 单位: 天
func (b *BRLog) SetRotationTime(day time.Duration) {
	b.RotationTime = day * 24 * time.Hour
}

// SetLevel 设置log level
// debug|info|warn|error|fatal|panic
func (b *BRLog) SetLevel(level string) {
	switch strings.ToLower(level) {
	case "panic":
		b.LogLevel = log.PanicLevel
	case "fatal":
		b.LogLevel = log.FatalLevel
	case "error":
		b.LogLevel = log.ErrorLevel
	case "warn", "warning":
		b.LogLevel = log.WarnLevel
	case "info":
		b.LogLevel = log.InfoLevel
	default:
		b.LogLevel = log.DebugLevel
	}
}

// SetDateFormat 设置日期格式
// format "%Y-%m-%d" | "%Y%m%d"
func (b *BRLog) SetDateFormat(format string) {
	b.FileNameDateFormat = format
}

// SetSeparateLevelLog 设置是否分离不同级别的日志到不同的文件
func (b *BRLog) SetSeparateLevelLog(yesorno bool) {
	b.IsSeparateLevelLog = yesorno
}

// setNull 相当于/dev/null
func setNull() *bufio.Writer {
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return nil
	}
	return bufio.NewWriter(src)
}

// GetLogger getlogger
func (b *BRLog) GetLogger() (*BRLogger, error) {
	logrusLogger = log.New()
	switch b.LogType {
	case "text":
		logrusLogger.Formatter = &log.TextFormatter{
			TimestampFormat: b.TimestampFormat,
		}
	default:
		logrusLogger.Formatter = &log.JSONFormatter{
			TimestampFormat: b.TimestampFormat,
		}
	}

	logrusLogger.Level = b.LogLevel
	maxAge := rotatelogs.WithMaxAge(b.MaxAge)
	ratateDuration := rotatelogs.WithRotationTime(b.RotationTime)

	if b.IsSeparateLevelLog {
		debugFileName := b.LogPath + ".debug"
		debugWriter, err := rotatelogs.New(
			fmt.Sprintf("%s.%s", debugFileName, b.FileNameDateFormat),
			rotatelogs.WithLinkName(debugFileName),
			maxAge,
			ratateDuration,
		)
		if err != nil {
			return nil, err
		}

		infoFileName := b.LogPath + ".info"
		infoWriter, err := rotatelogs.New(
			fmt.Sprintf("%s.%s", infoFileName, b.FileNameDateFormat),
			rotatelogs.WithLinkName(infoFileName),
			maxAge,
			ratateDuration,
		)
		if err != nil {
			return nil, err
		}

		warningFileName := b.LogPath + ".warn"
		warningWriter, err := rotatelogs.New(
			fmt.Sprintf("%s.%s", warningFileName, b.FileNameDateFormat),
			rotatelogs.WithLinkName(warningFileName),
			maxAge,
			ratateDuration,
		)
		if err != nil {
			return nil, err
		}

		errorFileName := b.LogPath + ".error"
		errorWriter, err := rotatelogs.New(
			fmt.Sprintf("%s.%s", errorFileName, b.FileNameDateFormat),
			rotatelogs.WithLinkName(errorFileName),
			maxAge,
			ratateDuration,
		)
		if err != nil {
			return nil, err
		}

		// 文件 hook, 不同的级别 设置输出不同的文件
		fileHook := lfshook.NewHook(lfshook.WriterMap{
			log.DebugLevel: debugWriter, // 为不同级别设置不同的输出目的
			log.InfoLevel:  infoWriter,
			log.WarnLevel:  warningWriter,
			log.ErrorLevel: errorWriter,
			log.FatalLevel: errorWriter,
			log.PanicLevel: errorWriter,
		}, logrusLogger.Formatter)

		logrusLogger.Hooks.Add(fileHook)
	} else {
		writer, err := rotatelogs.New(
			fmt.Sprintf("%s.%s", b.LogPath, b.FileNameDateFormat),
			maxAge,
			ratateDuration,
		)
		if err != nil {
			return nil, err
		}

		fileHook := lfshook.NewHook(lfshook.WriterMap{
			log.DebugLevel: writer,
			log.InfoLevel:  writer,
			log.WarnLevel:  writer,
			log.ErrorLevel: writer,
			log.FatalLevel: writer,
			log.PanicLevel: writer,
		}, logrusLogger.Formatter)

		logrusLogger.Hooks.Add(fileHook)
	}

	if b.LogLevel != log.DebugLevel {
		if out := setNull(); out != nil {
			logrusLogger.Out = setNull()
		} else {
			logrusLogger.Out = os.Stdout
		}
	}
	return &BRLogger{}, nil
}
