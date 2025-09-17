package log

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	httpTransport "github.com/go-kratos/kratos/v2/transport/http"
	"go.opentelemetry.io/otel/trace"
	"os"
	"ragx/app/pkg/go-tls"
	"ragx/app/pkg/utils"
	"ragx/app/pkg/utils/cast"
	"regexp"
	"time"
)

// Redacter defines how to log an object
type Redacter interface {
	Redact() string
}

// 是否不包装内部错误
var notWrapInternalError = cast.ToBool(os.Getenv("app.not_wrap_internal_error"))

// Server is an server logging middleware.
func Server(logger log.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var (
				code   int32
				reason string
				//kind      string
				operation string
				reqUri    string
				method    string
				realIp    string
				traceID   string
			)
			startTime := time.Now()
			// 获取traceID
			if span := trace.SpanContextFromContext(ctx); span.HasTraceID() {
				traceID = span.TraceID().String()
			} else {
				// 说明没有配置tracerProvider,直接生成一个traceID用于问题定位
				traceID = utils.UniqueID()
				traceIDHex, err := trace.TraceIDFromHex(traceID)
				if err != nil {
					log.Errorf("server middleware set trace_id error,err=%v", err)
				} else {
					sc := trace.NewSpanContext(trace.SpanContextConfig{
						TraceID:    traceIDHex,
						TraceFlags: 01,
						Remote:     false,
					})
					ctx = trace.ContextWithSpanContext(ctx, sc)
				}
			}
			// 将ctx设置到线程本地缓存中，用于日志中获取traceID
			tls.SetContext(ctx)
			defer tls.Flush()
			if info, ok := transport.FromServerContext(ctx); ok {
				realIp = utils.GetClientIP(info.RequestHeader())
				//kind = info.Kind().String()
				operation = info.Operation()
				// 如果是http请求，获取uri
				if httpTr, ok2 := httpTransport.RequestFromServerContext(ctx); ok2 {
					reqUri = httpTr.URL.Path
					method = httpTr.Method
				}
				// 响应头设置traceID，用于问题快速定位
				info.ReplyHeader().Set("X-Trace-ID", traceID)
			}
			reply, err = handler(ctx, req)
			level, _ := extractError(err)
			args := extractArgs(req)
			// 入参超过1024个字符，不在日志中打印
			if len(args) > 1024 {
				args = ""
			}
			log.NewHelper(log.WithContext(ctx, logger)).Log(level,
				"start", startTime.Format("2006-01-02 15:04:05.999"),
				"msg", "request access logging",
				"req_method", method,
				"req_uri", reqUri,
				"latency", time.Since(startTime).String(),
				"code", code,
				"reason", reason,
				"real_ip", realIp,
				"req_body", args,
				//"kind", "server",
				//"component", kind,
				"operation", operation,
				//"stack", stack,
			)
			return
		}
	}
}

// Client is a client logging middleware.
func Client(logger log.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var (
				code      int32
				reason    string
				kind      string
				operation string
			)
			startTime := time.Now()
			if info, ok := transport.FromClientContext(ctx); ok {
				kind = info.Kind().String()
				operation = info.Operation()
			}
			reply, err = handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				code = se.Code
				reason = se.Reason
			}
			level, stack := extractError(err)
			log.NewHelper(log.WithContext(ctx, logger)).Log(level,
				"kind", "client",
				"component", kind,
				"operation", operation,
				"args", extractArgs(req),
				"code", code,
				"reason", reason,
				"stack", stack,
				"latency", time.Since(startTime).String(),
			)
			return
		}
	}
}

// extractArgs returns the string of the req
func extractArgs(req interface{}) string {
	if redacter, ok := req.(Redacter); ok {
		return redacter.Redact()
	}
	if stringer, ok := req.(fmt.Stringer); ok {
		reqStr := stringer.String()

		// 定义需要脱敏的字段及其正则表达式
		sensitiveFields := map[string]*regexp.Regexp{
			// 匹配 password 字段，支持单引号、双引号包裹或无引号的情况
			"password": regexp.MustCompile(`(?i)(password:\s*['"]?)([^'"\s]+)(['"]?)`),
			// 匹配 email 字段，格式如 "email: example@example.com"，支持引号
			//"email": regexp.MustCompile(`(?i)(email:\s*['"]?)([^'"\s@]+)@([^'"\s]+)(['"]?)`),
			// 可按需添加更多敏感字段，如 phone 等
			// "phone":    regexp.MustCompile(`(?i)(phone:\s*['"]?)(1[3-9]\d{2})\d{4}(\d{4})(['"]?)`),
		}

		// 对每个敏感字段进行脱敏处理
		for _, re := range sensitiveFields {
			switch re.String() {
			//case sensitiveFields["email"].String():
			//	// 对 email 脱敏，保留邮箱前两位和域名
			//	reqStr = re.ReplaceAllString(reqStr, "$1$2***@$3$4")
			default:
				// 其他字段统一替换为 ***
				reqStr = re.ReplaceAllString(reqStr, "$1***$3")
			}
		}

		return reqStr
	}
	return fmt.Sprintf("%+v", req)
}

// extractError returns the string of the error
func extractError(err error) (log.Level, string) {
	if err != nil {
		return log.LevelError, fmt.Sprintf("%+v", err)
	}
	return log.LevelInfo, ""
}
