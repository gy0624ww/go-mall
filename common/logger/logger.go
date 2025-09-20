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
	log  *logger
)

func Init() {
	once.Do(func() {
		log = &logger{
			_logger: _logger,
		}
	})
}

// 门面方法，方便使用
func Info(ctx context.Context, msg string, kv ...interface{}) {
	log.Info(ctx, msg, kv...)
}

func Debug(ctx context.Context, msg string, kv ...interface{}) {
	log.Debug(ctx, msg, kv...)
}

func Warn(ctx context.Context, msg string, kv ...interface{}) {
	log.Warn(ctx, msg, kv...)
}

func Error(ctx context.Context, msg string, kv ...interface{}) {
	log.Error(ctx, msg, kv...)
}

type logger struct {
	_logger *zap.Logger
}

func (l *logger) logWithContext(ctx context.Context, lvl zapcore.Level, msg string, kv ...interface{}) {
	// 保证要打印的日志信息成对出现
	if len(kv)%2 != 0 {
		kv = append(kv, "unknown")
	}
	// 日志行信息中添加追踪参数
	traceId, spanId, pSpanId := l.extractTraceInfo(ctx)
	kv = append(kv, "traceid", traceId, "spanid", spanId, "pspanid", pSpanId)

	// 增加日志调用者信息，方便查日志时定位程序位置
	funcName, file, line := l.getLoggerCallerInfo()
	kv = append(kv, "func", funcName, "file", file, "line", line)

	fields := make([]zap.Field, 0, len(kv)/2)
	for i := 0; i < len(kv); i += 2 {
		k := fmt.Sprintf("%v", kv[i])
		value := kv[i+1]
		fields = append(fields, l.convertToZapField(k, value))
	}
	ce := l._logger.Check(lvl, msg)
	ce.Write(fields...)
}

// 提取追踪信息的辅助方法
func (l *logger) extractTraceInfo(ctx context.Context) (traceId, spanId, pSpanId string) {
	if v := ctx.Value("traceid"); v != nil {
		if s, ok := v.(string); ok {
			traceId = s
		}
	}
	if v := ctx.Value("spanid"); v != nil {
		if s, ok := v.(string); ok {
			spanId = s
		}
	}
	if v := ctx.Value("pspanid"); v != nil {
		if s, ok := v.(string); ok {
			pSpanId = s
		}
	}
	return
}

func (l *logger) convertToZapField(k string, value interface{}) zap.Field {
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

func (l *logger) Debug(ctx context.Context, msg string, kv ...interface{}) {
	l.logWithContext(ctx, zapcore.DebugLevel, msg, kv...)
}

func (l *logger) Info(ctx context.Context, msg string, kv ...interface{}) {
	l.logWithContext(ctx, zapcore.InfoLevel, msg, kv...)
}

func (l *logger) Warn(ctx context.Context, msg string, kv ...interface{}) {
	l.logWithContext(ctx, zapcore.WarnLevel, msg, kv...)
}

func (l *logger) Error(ctx context.Context, msg string, kv ...interface{}) {
	l.logWithContext(ctx, zapcore.ErrorLevel, msg, kv...)
}

// getLoggerCallerInfo 日志调用者信息 -- 方法名，文件名，行号
func (l *logger) getLoggerCallerInfo() (funcName, file string, line int) {
	// 第0层：当前函数getLoggerCallerInfo; 第1层:l.log; 第2层:l.Debug/Info/Warn/Error 第3层: Info/Warn/Error门面方法 第4层: 调用者logger.Info/Warn/Error
	pc, file, line, ok := runtime.Caller(4)
	if !ok {
		return
	}
	file = path.Base(file)
	funcName = runtime.FuncForPC(pc).Name()
	return
}
