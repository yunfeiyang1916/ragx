package logger

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"runtime"
	"sort"
	"strings"
	"time"

	"ragx/app/pkg/utils/cast"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	log *logrus.Logger
}

func NewLogger(logger *logrus.Logger) log.Logger {
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(DefaultJsonFormatter)
	return &Logger{
		log: logger,
	}
}

// 转成logrus.Logger
func NewLogrus(log log.Logger) *logrus.Logger {
	return logrus.StandardLogger()
}

func (l *Logger) Log(level log.Level, keyvals ...any) (err error) {
	var (
		logrusLevel logrus.Level
		fields      logrus.Fields = make(map[string]any)
		msg         string
	)

	switch level {
	case log.LevelDebug:
		logrusLevel = logrus.DebugLevel
	case log.LevelInfo:
		logrusLevel = logrus.InfoLevel
	case log.LevelWarn:
		logrusLevel = logrus.WarnLevel
	case log.LevelError:
		logrusLevel = logrus.ErrorLevel
	case log.LevelFatal:
		logrusLevel = logrus.FatalLevel
	default:
		logrusLevel = logrus.DebugLevel
	}

	if logrusLevel > l.log.Level {
		return
	}

	if len(keyvals) == 0 {
		return nil
	}
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "")
	}
	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			continue
		}
		if key == logrus.FieldKeyMsg {
			msg, _ = keyvals[i+1].(string)
			continue
		}
		fields[key] = keyvals[i+1]
	}

	if len(fields) > 0 {
		l.log.WithFields(fields).Log(logrusLevel, msg)
	} else {
		l.log.Log(logrusLevel, msg)
	}

	return
}

// 字段打印顺序
var FieldsOrder = []string{"req_method", "req_uri", "latency", "code", "reason", "start", "real_ip", "http_status", "req_body", "resp_body"}

var DefaultNestedFormatter = &NestedLogrusFormatter{
	TimestampFormat: time.RFC3339,
	CallerFirst:     true,
}

// NestedLogrusFormatter - logrus formatter, implements logrus.Formatter
type NestedLogrusFormatter struct {
	// TimestampFormat - default: time.StampMilli = "Jan _2 15:04:05.000"
	TimestampFormat string

	// HideKeys - show [fieldValue] instead of [fieldKey:fieldValue]
	HideKeys bool

	// NoColors - disable colors
	NoColors bool

	// NoFieldsColors - apply colors only to the level, default is level + fields
	NoFieldsColors bool

	// NoFieldsSpace - no space between fields
	NoFieldsSpace bool

	// ShowFullLevel - show a full level [WARNING] instead of [WARN]
	ShowFullLevel bool

	// NoUppercaseLevel - no upper case for level value
	NoUppercaseLevel bool

	// TrimMessages - trim whitespaces on messages
	TrimMessages bool

	// CallerFirst - print caller info first
	CallerFirst bool

	// CustomCallerFormatter - set custom formatter for caller info
	CustomCallerFormatter func(*runtime.Frame) string
}

// Format an log entry
func (f *NestedLogrusFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	levelColor := getColorByLevel(entry.Level)

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.StampMilli
	}

	// output buffer
	b := &bytes.Buffer{}

	// write time
	b.WriteString(entry.Time.Format(timestampFormat))

	// write level
	var level string
	if f.NoUppercaseLevel {
		level = entry.Level.String()
	} else {
		level = strings.ToUpper(entry.Level.String())
	}

	if f.CallerFirst {
		f.writeCaller(b, entry)
	}

	if !f.NoColors {
		fmt.Fprintf(b, "\x1b[%dm", levelColor)
	}

	b.WriteString(" [")
	if f.ShowFullLevel {
		b.WriteString(level)
	} else {
		b.WriteString(level[:4])
	}
	b.WriteString("]")

	if !f.NoFieldsSpace {
		b.WriteString(" ")
	}

	if !f.NoColors && f.NoFieldsColors {
		b.WriteString("\x1b[0m")
	}

	// write fields
	if len(FieldsOrder) == 0 {
		f.writeFields(b, entry)
	} else {
		f.writeOrderedFields(b, entry)
	}

	if f.NoFieldsSpace {
		b.WriteString(" ")
	}

	if !f.NoColors && !f.NoFieldsColors {
		b.WriteString("\x1b[0m")
	}

	// write message
	if f.TrimMessages {
		b.WriteString(strings.TrimSpace(entry.Message))
	} else {
		b.WriteString(entry.Message)
	}

	if !f.CallerFirst {
		f.writeCaller(b, entry)
	}

	b.WriteByte('\n')

	return b.Bytes(), nil
}

func (f *NestedLogrusFormatter) writeCaller(b *bytes.Buffer, entry *logrus.Entry) {
	if entry.HasCaller() {
		if f.CustomCallerFormatter != nil {
			fmt.Fprintf(b, f.CustomCallerFormatter(entry.Caller))
		} else {
			dir, file := path.Split(entry.Caller.File)
			dir = path.Base(dir)
			file = path.Join(dir, file)
			fmt.Fprintf(
				b,
				" (%s:%d)",
				file,
				entry.Caller.Line,
			)

			//fmt.Fprintf(
			//	b,
			//	" (%s:%d %s)",
			//	entry.Caller.File,
			//	entry.Caller.Line,
			//	entry.Caller.Function,
			//)
		}
	}
}

func (f *NestedLogrusFormatter) writeFields(b *bytes.Buffer, entry *logrus.Entry) {
	if len(entry.Data) != 0 {
		fields := make([]string, 0, len(entry.Data))
		for field := range entry.Data {
			fields = append(fields, field)
		}

		sort.Strings(fields)

		for _, field := range fields {
			f.writeField(b, entry, field)
		}
	}
}

func (f *NestedLogrusFormatter) writeOrderedFields(b *bytes.Buffer, entry *logrus.Entry) {
	length := len(entry.Data)
	foundFieldsMap := map[string]bool{}
	for _, field := range FieldsOrder {
		if _, ok := entry.Data[field]; ok {
			foundFieldsMap[field] = true
			length--
			f.writeField(b, entry, field)
		}
	}

	if length > 0 {
		notFoundFields := make([]string, 0, length)
		for field := range entry.Data {
			if foundFieldsMap[field] == false {
				notFoundFields = append(notFoundFields, field)
			}
		}

		sort.Strings(notFoundFields)

		for _, field := range notFoundFields {
			f.writeField(b, entry, field)
		}
	}
}

func (f *NestedLogrusFormatter) writeField(b *bytes.Buffer, entry *logrus.Entry, field string) {
	val, _ := cast.ToStringJson(entry.Data[field])
	if f.HideKeys {
		fmt.Fprintf(b, "[%v]", val)
	} else {
		fmt.Fprintf(b, "[%s:%v]", field, val)
	}

	if !f.NoFieldsSpace {
		b.WriteString(" ")
	}
}

const (
	colorRed    = 31
	colorYellow = 33
	colorBlue   = 36
	colorGray   = 37
)

func getColorByLevel(level logrus.Level) int {
	switch level {
	case logrus.DebugLevel, logrus.TraceLevel:
		return colorGray
	case logrus.WarnLevel:
		return colorYellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return colorRed
	default:
		return colorBlue
	}
}
