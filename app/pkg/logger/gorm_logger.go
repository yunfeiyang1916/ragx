package logger

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"ragx/app/pkg/go-tls"

	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

var (
	DefaultGormLogger = NewGormLogger()
)

type GormLogger struct {
	logger.Config
}

func NewGormLogger() logger.Interface {
	return &GormLogger{
		Config: logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		},
	}
}

// LogMode log mode
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info print info
func (l GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		//l.Printf(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (l GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		//l.Printf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (l GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		//l.Printf(l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Trace print sql message
func (l GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}
	var traceID string
	span := trace.SpanFromContext(ctx)
	tt := span.SpanContext().TraceID()
	if tt.IsValid() {
		traceID = tt.String()
	} else {
		if ctx, ok := tls.GetContext(); ok {
			span = trace.SpanFromContext(ctx)
			traceID = span.SpanContext().TraceID().String()
		}
	}

	elapsed := time.Since(begin)
	info := SqlInfo{
		FileWithLine: utils.FileWithLineNum(),
		Duration:     float64(elapsed.Nanoseconds()/1e4) / 100.0,
		TraceID:      traceID,
	}
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		info.Sql = sql
		info.Rows = rows
		info.ErrorMsg = err.Error()

	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		info.Sql = sql
		info.Rows = rows
		info.ErrorMsg = fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
	case l.LogLevel == logger.Info:
		sql, rows := fc()
		info.Sql = sql
		info.Rows = rows
	default:
		return
	}
	info.SetCustomFormat("")
	fmt.Println(info.Output())
}

var (
	// Map from format's placeholders to printf verbs
	phfs map[string]string
	// Default format of log message
	defaultFmt = "\u001B[33m[%.2[1]fms] \u001B[34;1m[rows:%[2]d] %[3]s \u001B[31;1m%[4]s \u001B[34;1m[trace.id:%[5]s] \u001B[0m"
)

type SqlInfo struct {
	FileWithLine string
	Duration     float64
	Sql          string
	Rows         int64
	ErrorMsg     string
	TraceID      string

	Format string
}

func init() {
	initFormatPlaceholders()
}

// Initializes the map of placeholders
func initFormatPlaceholders() {
	phfs = map[string]string{
		"%{duration}":  "[%.2[1]fms]",
		"%{rows}":      "[rows:%[2]d]",
		"%{sql}":       "%[3]s",
		"%{error_msg}": "%[4]s",
		"%{trace_id}":  "[%[5]s]",
	}
}

func (r *SqlInfo) Output() string {
	msg := fmt.Sprintf(r.Format,
		//r.FileWithLine, // %[1] // %{file_with_line}
		r.Duration, // %[1] // %{duration}
		r.Rows,     // %[2] // %{rows}
		r.Sql,      // %[3] // %{sql}
		r.ErrorMsg, // %[4] // %{error_msg}
		r.TraceID,  // %[5] // %{trace_id}
	)
	// Ignore printf errors if len(args) > len(verbs)
	if i := strings.LastIndex(msg, "%!(EXTRA"); i != -1 {
		return msg[:i]
	}
	return msg
}
func (r *SqlInfo) SetCustomFormat(format string) {
	r.Format = parseFormat(format)
}

// Analyze and represent format string as printf format string and time format
func parseFormat(format string) (msgfmt string) {
	if len(format) < 6 /* (len of "%{sql} */ {
		return defaultFmt
	}
	idx := strings.IndexRune(format, '%')
	for idx != -1 {
		msgfmt += format[:idx]
		format = format[idx:]
		if len(format) > 2 {
			if format[1] == '{' {
				// end of curr verb pos
				if jdx := strings.IndexRune(format, '}'); jdx != -1 {
					// next verb pos
					idx = strings.Index(format[1:], "%{")
					// incorrect verb found ("...%{wefwef ...") but after
					// this, new verb (maybe) exists ("...%{inv %{verb}...")
					if idx != -1 && idx < jdx {
						msgfmt += "%%"
						format = format[1:]
						continue
					}
					// get verb and arg
					verb := ph2verb(format[:jdx+1])
					msgfmt += verb

					format = format[jdx+1:]
				} else {
					format = format[1:]
				}
			} else {
				msgfmt += "%%"
				format = format[1:]
			}
		}
		idx = strings.IndexRune(format, '%')
	}
	msgfmt += format
	return
}

func ph2verb(ph string) (verb string) {
	n := len(ph)
	if n < 4 {
		return ``
	}
	if ph[0] != '%' || ph[1] != '{' || ph[n-1] != '}' {
		return ``
	}
	return phfs[ph]
}
