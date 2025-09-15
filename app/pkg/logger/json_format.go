package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"ragx/app/pkg/go-tls"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wk8/go-ordered-map/v2"
	"go.opentelemetry.io/otel/trace"
)

var DefaultJsonFormatter = &JSONFormatter{
	TimestampFormat: time.RFC3339,
}

// JSONFormatter formats logs into parsable json
type JSONFormatter struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	// The format to use is the same than for time.Format or time.Parse from the standard
	// library.
	// The standard Library already provides a set of predefined format.
	TimestampFormat string

	// DisableTimestamp allows disabling automatic timestamps in output
	DisableTimestamp bool

	// DisableHTMLEscape allows disabling html escaping in output
	DisableHTMLEscape bool

	// DataKey allows users to put all the log entry parameters into a nested dictionary at a given key.
	DataKey string

	// FieldMap allows users to customize the names of keys for default fields.
	// As an example:
	// formatter := &JSONFormatter{
	//   	FieldMap: FieldMap{
	// 		 FieldKeyTime:  "@timestamp",
	// 		 FieldKeyLevel: "@level",
	// 		 FieldKeyMsg:   "@message",
	// 		 FieldKeyFunc:  "@caller",
	//    },
	// }
	FieldMap FieldMap

	// CallerPrettyfier can be set by the user to modify the content
	// of the function and file keys in the json data when ReportCaller is
	// activated. If any of the returned value is the empty string the
	// corresponding key will be removed from json fields.
	CallerPrettyfier func(*runtime.Frame) (function string, file string)

	// PrettyPrint will indent all json logs
	PrettyPrint bool
}

type fieldKey string

// FieldMap allows customization of the key names for default fields.
type FieldMap map[fieldKey]string

func (f FieldMap) resolve(key fieldKey) string {
	if k, ok := f[key]; ok {
		return k
	}

	return string(key)
}

// FieldKeyTime 时间字段
var FieldKeyTime fieldKey = "ts"

var FieldKeyTraceID = "trace.id"

// Format renders a single log entry
func (f *JSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 使用有序映射
	data := orderedmap.New[string, interface{}]()

	// 按顺序添加字段
	if !f.DisableTimestamp {
		timestampFormat := f.TimestampFormat
		if timestampFormat == "" {
			timestampFormat = time.RFC3339
		}
		data.Set(f.FieldMap.resolve(FieldKeyTime), entry.Time.Format(timestampFormat))
	}
	data.Set(f.FieldMap.resolve(logrus.FieldKeyLevel), entry.Level.String())
	data.Set(f.FieldMap.resolve(logrus.FieldKeyMsg), entry.Message)

	// 添加 entry.Data 中的字段
	// 按 FieldsOrder 顺序添加
	for _, field := range FieldsOrder {
		if v, ok := entry.Data[field]; ok {
			data.Set(field, v)
			delete(entry.Data, field)
		}
	}
	// 添加其他字段
	for field, v := range entry.Data {
		data.Set(field, v)
	}

	if entry.HasCaller() {
		funcVal := entry.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
		if f.CallerPrettyfier != nil {
			funcVal, fileVal = f.CallerPrettyfier(entry.Caller)
		}
		if funcVal != "" {
			data.Set(f.FieldMap.resolve(logrus.FieldKeyFunc), funcVal)
		}
		if fileVal != "" {
			data.Set(f.FieldMap.resolve(logrus.FieldKeyFile), fileVal)
		}
	}
	// 如果没有 trace.id，则添加
	if _, ok := data.Get(FieldKeyTraceID); !ok {
		var traceID string
		ctx := entry.Context
		if ctx == nil {
			ctx, _ = tls.GetContext()
		}
		if ctx != nil {
			if span := trace.SpanContextFromContext(ctx); span.HasTraceID() {
				traceID = span.TraceID().String()
			}
		}
		data.Set(FieldKeyTraceID, traceID)
	}

	// log info
	if logInfo, ok := tls.GetLogInfo(); ok {
		data.Set("tenant_name", logInfo.TenantName)
		data.Set("site_name", logInfo.SiteName)
		data.Set("username", logInfo.Username)
	}

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	// 增加颜色打印
	levelColor := getColorByLevel(entry.Level)
	b.WriteString(fmt.Sprintf("\x1b[%dm", levelColor))

	encoder := json.NewEncoder(b)
	encoder.SetEscapeHTML(!f.DisableHTMLEscape)
	if f.PrettyPrint {
		encoder.SetIndent("", "  ")
	}

	if err := encoder.Encode(data); err != nil {
		return nil, fmt.Errorf("failed to marshal fields to JSON, %w", err)
	}
	// 重置颜色
	b.WriteString("\x1b[0m")
	return b.Bytes(), nil
}
