// Package logger
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-23
package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/gookit/color"
)

const prefix string = "[cci]"

const (
	tagDebug = "debug"
	tagInfo  = "info"
	tagWarn  = "warn"
	tagError = "error"
	tagFatal = "fatal"
)

// Level is a log level.
type Level int

// Log levels.
const (
	sDebug Level = iota
	sInfo
	sWarn
	sError
	sFatal
)

// Destination is a log destination.
type Destination int

const (
	// DestinationStdout writes logs to the standard output.
	DestinationStdout Destination = iota

	// DestinationFile writes logs to a file.
	DestinationFile

	// DestinationSyslog writes logs to the system logger.
	DestinationSyslog
)

// Logger is a log handler.
type Logger struct {
	level        Level
	destinations map[Destination]string

	mutex   sync.Mutex
	file    *os.File
	syslog  io.WriteCloser
	buffers map[Destination]bytes.Buffer
}

type LogConfig struct {
	Level   string
	Verbose bool
	Syslog  bool
	LogFile string
}

// New allocates a log handler.
func New(c LogConfig) (*Logger, error) {
	return newLogger(c.Level, c.Verbose, c.Syslog, c.LogFile)
}

func newLogger(level string, verbose bool, syslog bool, filePath string) (*Logger, error) {
	var destinations = map[Destination]string{}

	if verbose {
		destinations[DestinationStdout] = "stdout"
	}

	if len(filePath) > 0 {
		destinations[DestinationFile] = "file"
	}

	if syslog {
		destinations[DestinationSyslog] = "syslog"
	}

	return initialize(getLevel(level), destinations, filePath)
}

// Close closes a log handler.
func (l *Logger) Close() {
	if l.file != nil {
		l.file.Close()
	}

	if l.syslog != nil {
		l.syslog.Close()
	}
}

// Log writes a log entry.
func (l *Logger) Log(level Level, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	if _, ok := l.destinations[DestinationStdout]; ok {
		buff := l.getDestinationBuffer(DestinationStdout)
		logEntry(&buff, level, true, format, args)
		print(buff.String())
	}

	if _, ok := l.destinations[DestinationFile]; ok {
		buff := l.getDestinationBuffer(DestinationFile)
		logEntry(&buff, level, false, format, args)
		_, _ = l.file.Write(buff.Bytes())
	}

	if _, ok := l.destinations[DestinationSyslog]; ok {
		buff := l.getDestinationBuffer(DestinationSyslog)
		logEntry(&buff, level, false, format, args)
		_, _ = l.syslog.Write(buff.Bytes())
	}
}

// Info logs with the Info severity.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Info(v ...interface{}) {
	l.Log(sInfo, fmt.Sprint(v...))
}

// Infoln logs with the Info severity.
// Arguments are handled in the manner of fmt.Println.
func (l *Logger) Infoln(v ...interface{}) {
	l.Log(sInfo, fmt.Sprintln(v...))
}

// Infof logs with the Info severity.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Infof(format string, v ...interface{}) {
	l.Log(sInfo, fmt.Sprintf(format, v...))
}

// Warning logs with the Warning severity.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Warning(v ...interface{}) {
	l.Log(sWarn, fmt.Sprint(v...))
}

// Warningln logs with the Warning severity.
// Arguments are handled in the manner of fmt.Println.
func (l *Logger) Warningln(v ...interface{}) {
	l.Log(sWarn, fmt.Sprintln(v...))
}

// Warningf logs with the Warning severity.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Warningf(format string, v ...interface{}) {
	l.Log(sWarn, fmt.Sprintf(format, v...))
}

// Error logs with the ERROR severity.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Error(v ...interface{}) {
	l.Log(sError, fmt.Sprint(v...))
}

// Errorln logs with the ERROR severity.
// Arguments are handled in the manner of fmt.Println.
func (l *Logger) Errorln(v ...interface{}) {
	l.Log(sError, fmt.Sprintln(v...))
}

// Errorf logs with the Error severity.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Log(sError, fmt.Sprintf(format, v...))
}

// Fatal uses the default logger, logs with the sFatal severity,
// and ends with os.Exit(1).
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Fatal(v ...interface{}) {
	l.Log(sFatal, fmt.Sprint(v...))
	l.Close()
	os.Exit(1)
}

// Fatalln uses the default logger, logs with the sFatal severity,
// and ends with os.Exit(1).
// Arguments are handled in the manner of fmt.Println.
func (l *Logger) Fatalln(v ...interface{}) {
	l.Log(sFatal, fmt.Sprintln(v...))
	l.Close()
	os.Exit(1)
}

// Fatalf uses the default logger, logs with the sFatal severity,
// and ends with os.Exit(1).
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Log(sFatal, fmt.Sprintf(format, v...))
	l.Close()
	os.Exit(1)
}

func GetDefaultLevelName() string {
	return levelNames()[0]
}

func GetLevelName(level Level) string {
	return levelNames()[level]
}

func initialize(level Level, destinations map[Destination]string, filePath string) (*Logger, error) {
	var err error
	lh := Logger{
		level:        level,
		destinations: destinations,
	}

	if _, ok := destinations[DestinationFile]; ok {
		lh.file, err = os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			lh.Close()
			return nil, ErrCanNotOpenLogFile(filePath, err)
		}
	}

	if _, ok := destinations[DestinationSyslog]; ok {
		lh.syslog, err = NewSyslog(prefix)
		if err != nil {
			lh.Close()
			return nil, ErrCanNotInitSyslog(err)
		}
	}
	lh.buffers = map[Destination]bytes.Buffer{}
	for d := range destinations {
		var b bytes.Buffer
		lh.buffers[d] = b
	}

	return &lh, nil
}

// https://golang.org/src/log/log.go#L78
func itoa(i int, wid int) []byte {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	return b[bp:]
}

func bufferDate(now time.Time, buf *bytes.Buffer) {
	year, month, day := now.Date()
	buf.Write(itoa(year, 4))
	buf.WriteByte('/')
	buf.Write(itoa(int(month), 2))
	buf.WriteByte('/')
	buf.Write(itoa(day, 2))
	buf.WriteByte(' ')
}

func bufferClock(now time.Time, buf *bytes.Buffer) {
	hour, min, sec := now.Clock()
	buf.Write(itoa(hour, 2))
	buf.WriteByte(':')
	buf.Write(itoa(min, 2))
	buf.WriteByte(':')
	buf.Write(itoa(sec, 2))
	buf.WriteByte(' ')
}

func writeTime(buf *bytes.Buffer, doColor bool) {
	intBuffer := &bytes.Buffer{}

	now := time.Now()

	// date
	bufferDate(now, intBuffer)

	// time
	bufferClock(now, intBuffer)

	if doColor {
		buf.WriteString(color.RenderString(color.Gray.Code(), intBuffer.String()))
	} else {
		buf.WriteString(intBuffer.String())
	}
}

func (l *Logger) getDestinationBuffer(destination Destination) bytes.Buffer {
	return l.buffers[destination]
}
func writeLevel(buf *bytes.Buffer, level Level, doColor bool) {
	if doColor {
		buf.WriteString(color.RenderString(levelColor(level), levelPrefix(level)))
	} else {
		buf.WriteString(levelPrefix(level))
	}
}

func writeContent(buf *bytes.Buffer, format string, args []interface{}) {
	buf.Write([]byte(fmt.Sprintf(format, args...)))
	buf.WriteByte('\n')
}

func logEntry(buff *bytes.Buffer, level Level, doColor bool, format string, args ...interface{}) {
	buff.Reset()
	writeTime(buff, doColor)
	writeLevel(buff, level, doColor)
	writeContent(buff, format, args)
}

func getLevel(level string) Level {
	switch level {
	case tagInfo:
		return sInfo
	case tagWarn:
		return sWarn
	case tagError:
		return sError
	case tagFatal:
		return sError
	default:
		panic(fmt.Sprintln("unrecognized severity:", level))
	}
}

func levelNames() []string {
	return []string{tagDebug, tagInfo, tagWarn, tagError, tagFatal}
}

func levelColors() []string {
	return []string{color.Debug.Code(), color.Info.Code(), color.Warn.Code(), color.Error.Code(), color.Danger.Code()}
}

func levelColor(level Level) string {
	return levelColors()[level]
}

func levelPrefixes() []string {
	return []string{"D ", "I ", "W ", "E ", "F "}
}

func levelPrefix(level Level) string {
	return levelPrefixes()[level]
}
