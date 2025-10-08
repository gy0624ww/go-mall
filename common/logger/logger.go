package logger

import (
	"context"
	"fmt"
	"path"
	"runtime"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	once sync.Once
	f    *facade
)

type facade struct {
	_logger *zap.Logger
}

func (f *facade) log(ctx context.Context, lvl zapcore.Level, msg string, kv ...interface{}) {
	fields := makeLogFields(ctx, kv...)
	ce := f._logger.Check(lvl, msg)
	ce.Write(fields...)
}

func logFacade() *facade {
	once.Do(func() {
		f = &facade{
			_logger: _logger,
		}
	})
	return f
}

// 门面方法，方便使用
func Info(ctx context.Context, msg string, kv ...interface{}) {
	logFacade().log(ctx, zapcore.InfoLevel, msg, kv...)
}

func Debug(ctx context.Context, msg string, kv ...interface{}) {
	logFacade().log(ctx, zapcore.DebugLevel, msg, kv...)
}

func Warn(ctx context.Context, msg string, kv ...interface{}) {
	logFacade().log(ctx, zapcore.WarnLevel, msg, kv...)
}

func Error(ctx context.Context, msg string, kv ...interface{}) {
	logFacade().log(ctx, zapcore.ErrorLevel, msg, kv...)
}

func makeLogFields(ctx context.Context, kv ...interface{}) []zap.Field {
	// 保证要打印的日志信息成对出现
	if len(kv)%2 != 0 {
		kv = append(kv, "unknown")
	}
	// 日志行信息中添加追踪参数
	traceId, spanId, pSpanId := extractTraceInfo(ctx)
	kv = append(kv, "traceid", traceId, "spanid", spanId, "pspanid", pSpanId)

	// 增加日志调用者信息，方便查日志时定位程序位置
	funcName, file, line := getLoggerCallerInfo()
	kv = append(kv, "func", funcName, "file", file, "line", line)

	fields := make([]zap.Field, 0, len(kv)/2)
	for i := 0; i < len(kv); i += 2 {
		k := fmt.Sprintf("%v", kv[i])
		value := kv[i+1]
		fields = append(fields, convertToZapField(k, value))
	}
	return fields
}

// getLoggerCallerInfo 日志调用者信息 -- 方法名, 文件名, 行号
func getLoggerCallerInfo() (funcName, file string, line int) {
	// 第0层：当前函数getLoggerCallerInfo; 第1层:makeLogFields; 第2层:facade.log; 第3层: Info/Warn/Error门面方法 第4层: 调用者logger.Info/Warn/Error
	pc, file, line, ok := runtime.Caller(4) // 回溯拿调用日志方法的业务函数的信息
	if !ok {
		return
	}
	file = path.Base(file)
	funcName = runtime.FuncForPC(pc).Name()
	return
}

// 提取追踪信息的辅助方法
func extractTraceInfo(ctx context.Context) (traceId, spanId, pSpanId string) {
	if v := ctx.Value("traceid"); v != nil {
		traceId = ctx.Value("traceid").(string)
	}
	if v := ctx.Value("spanid"); v != nil {
		spanId = ctx.Value("spanid").(string)
	}
	if v := ctx.Value("pspanid"); v != nil {
		pSpanId = ctx.Value("pspanid").(string)
	}
	return
}

func convertToZapField(k string, value interface{}) zap.Field {
	switch v := value.(type) {
	case string:
		return zap.String(k, v)
	case int:
		return zap.Int(k, v)
	case bool:
		return zap.Bool(k, v)
	case int64:
		return zap.Int64(k, v)
	case float64:
		return zap.Float64(k, v)
	case uint64:
		return zap.Uint64(k, v)
	case uint32:
		return zap.Uint32(k, v)
	case uint16:
		return zap.Uint16(k, v)
	case uint8:
		return zap.Uint8(k, v)
	case int32:
		return zap.Int32(k, v)
	case int16:
		return zap.Int16(k, v)
	case int8:
		return zap.Int8(k, v)
	case uint:
		return zap.Uint(k, v)
	case float32:
		return zap.Float32(k, v)
	default:
		return zap.Any(k, v)
	}
}
