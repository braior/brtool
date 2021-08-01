package brtool

import (
	log "github.com/sirupsen/logrus"
)

var (
	logrusLogger *log.Logger
)

// BLogger blogger
type BLogger struct{}

// logrusEntry 返回logrusEntry
func (b *BLogger) logrusEntry(commonfFileds map[string]interface{}) *log.Entry {
	return logrusLogger.WithFields(log.Fields(commonfFileds))
}

// Debug Debug日志
func (b *BLogger) Debug(commonfFileds map[string]interface{}, message string) {
	b.logrusEntry(commonfFileds).Debug("%s", message)
}
