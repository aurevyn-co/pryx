// Package security provides secure logging utilities
package security

import (
	"fmt"
	"log"
	"os"
)

// SecureLogger wraps a standard logger with automatic redaction
type SecureLogger struct {
	logger   *log.Logger
	redactor *Redactor
	prefix   string
}

// NewSecureLogger creates a new secure logger
func NewSecureLogger(prefix string) *SecureLogger {
	return &SecureLogger{
		logger:   log.New(os.Stderr, prefix, log.LstdFlags),
		redactor: NewRedactor(),
		prefix:   prefix,
	}
}

// NewSecureLoggerWithWriter creates a secure logger with custom output
func NewSecureLoggerWithWriter(prefix string, output *os.File) *SecureLogger {
	return &SecureLogger{
		logger:   log.New(output, prefix, log.LstdFlags),
		redactor: NewRedactor(),
		prefix:   prefix,
	}
}

// Printf logs a redacted formatted message
func (l *SecureLogger) Printf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.logger.Print(l.redactor.RedactString(msg))
}

// Print logs a redacted message
func (l *SecureLogger) Print(v ...interface{}) {
	msg := fmt.Sprint(v...)
	l.logger.Print(l.redactor.RedactString(msg))
}

// Println logs a redacted line
func (l *SecureLogger) Println(v ...interface{}) {
	msg := fmt.Sprint(v...)
	l.logger.Println(l.redactor.RedactString(msg))
}

// Fatal logs a redacted fatal message
func (l *SecureLogger) Fatal(v ...interface{}) {
	msg := fmt.Sprint(v...)
	l.logger.Fatal(l.redactor.RedactString(msg))
}

// Fatalf logs a redacted formatted fatal message
func (l *SecureLogger) Fatalf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.logger.Fatal(l.redactor.RedactString(msg))
}

// Panic logs a redacted panic message
func (l *SecureLogger) Panic(v ...interface{}) {
	msg := fmt.Sprint(v...)
	l.logger.Panic(l.redactor.RedactString(msg))
}

// Panicf logs a redacted formatted panic message
func (l *SecureLogger) Panicf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.logger.Panic(l.redactor.RedactString(msg))
}

// SetPrefix sets the logger prefix
func (l *SecureLogger) SetPrefix(prefix string) {
	l.logger.SetPrefix(prefix)
	l.prefix = prefix
}

// SetRedactor sets a custom redactor
func (l *SecureLogger) SetRedactor(r *Redactor) {
	l.redactor = r
}

// GetRedactor returns the current redactor
func (l *SecureLogger) GetRedactor() *Redactor {
	return l.redactor
}

// Global secure logger instance
var secureLog = NewSecureLogger("")

// SecurePrintf logs with redaction using global logger
func SecurePrintf(format string, v ...interface{}) {
	secureLog.Printf(format, v...)
}

// SecurePrint logs with redaction using global logger
func SecurePrint(v ...interface{}) {
	secureLog.Print(v...)
}

// SecurePrintln logs with redaction using global logger
func SecurePrintln(v ...interface{}) {
	secureLog.Println(v...)
}

// SecureFatal logs fatal with redaction using global logger
func SecureFatal(v ...interface{}) {
	secureLog.Fatal(v...)
}

// SecureFatalf logs fatal formatted with redaction using global logger
func SecureFatalf(format string, v ...interface{}) {
	secureLog.Fatalf(format, v...)
}
