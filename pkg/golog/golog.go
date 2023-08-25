// Provides a set of tools and syntax sugar for printing logs.
// Warning, this package is designed to work with Datadog. We're not using OTEL conventions.
package golog

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/JhonatanRSantos/gocore/pkg/goenv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type (
	logType string
	Options func(ctx context.Context, message string, logger *Logger, logType logType)
)

var (
	logInfo  logType = "info"
	logWarn  logType = "warn"
	logDebug logType = "debug"
	logError logType = "error"

	logColorRed    = "\033[31m"
	logColorGreen  = "\033[32m"
	logColorYellow = "\033[33m"
	logColorBlue   = "\033[34m"

	logEnvironment = goenv.Test

	logWriter           = os.Stdout
	logger      *Logger = nil
	initLogOnce sync.Once
)

type Logger struct {
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	debugLogger *log.Logger
	errorLogger *log.Logger
	zapLogger   *zap.Logger
}

// SetEnv Set the environment config
func SetEnv(env goenv.Env) {
	logEnvironment = env
}

// Log Creats a new logger
func Log() *Logger {
	init := func() {
		if logger == nil {
			if logEnvironment == goenv.Test || logEnvironment == goenv.Local {
				logger = &Logger{
					infoLogger:  log.New(logWriter, fmt.Sprintf("%v[  INFO  ] ", logColorGreen), log.LUTC),
					warnLogger:  log.New(logWriter, fmt.Sprintf("%v[  WARN  ] ", logColorYellow), log.LUTC),
					debugLogger: log.New(logWriter, fmt.Sprintf("%v[  DEBUG ] ", logColorBlue), log.LUTC),
					errorLogger: log.New(logWriter, fmt.Sprintf("%v[  ERROR ] ", logColorRed), log.LUTC),
				}
			} else {
				var err error
				zapConfig := zap.Config{
					Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
					Development: true,
					Encoding:    "json",
					EncoderConfig: zapcore.EncoderConfig{
						TimeKey:        "timestamp",
						LevelKey:       "status",
						NameKey:        zapcore.OmitKey,
						CallerKey:      zapcore.OmitKey,
						FunctionKey:    zapcore.OmitKey,
						MessageKey:     "message",
						StacktraceKey:  zapcore.OmitKey,
						LineEnding:     zapcore.DefaultLineEnding,
						EncodeLevel:    zapcore.CapitalLevelEncoder,
						EncodeTime:     zapcore.TimeEncoderOfLayout(time.RFC3339),
						EncodeDuration: zapcore.StringDurationEncoder,
						EncodeCaller:   zapcore.ShortCallerEncoder,
					},
					OutputPaths:      []string{"stderr"},
					ErrorOutputPaths: []string{"stderr"},
				}

				logger = &Logger{}

				if logEnvironment == goenv.Development || logEnvironment == goenv.Staging {
					logger.zapLogger, err = zapConfig.Build()
				} else {
					zapConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
					zapConfig.Development = false
					logger.zapLogger, err = zapConfig.Build()
				}

				if err != nil {
					panic(fmt.Errorf("failed to load logger configs. Cause: %s", err))
				}
				zap.ReplaceGlobals(logger.zapLogger.Sugar().Desugar())
			}
		}
	}
	initLogOnce.Do(init)
	return logger
}

// WithTags Add a group of tags into current log. Tags are used to provide more context when writing logs.
func WithTags(tags map[string]interface{}) Options {
	return func(ctx context.Context, message string, logger *Logger, logType logType) {
		allTags := []interface{}{}
		if logEnvironment == goenv.Test || logEnvironment == goenv.Local {
			bs, _ := json.MarshalIndent(tags, "", "\t")
			allTags = append(allTags, string(bs))
		} else {
			allTags = parseTags("", tags)
		}

		switch logType {
		case logInfo:
			logger.infoLog(ctx, message, allTags...)
		case logWarn:
			logger.warnLog(ctx, message, allTags...)
		case logDebug:
			logger.debugLog(ctx, message, allTags...)
		case logError:
			logger.errorLog(ctx, message, allTags...)
		default:
			logger.errorLogger.Println("- invalid log type")
		}
	}
}

// parseTags Convert a map of tags to a slice of tags that will be usad by zap
func parseTags(name string, tags map[string]interface{}) []interface{} {
	allTags := []interface{}{}
	for tagKey, tagValue := range tags {
		tagFullName := fmt.Sprintf("%s.%s", name, tagKey)

		if name == "" {
			tagFullName = tagKey
		}

		if subTags, isSubTags := tagValue.(map[string]interface{}); !isSubTags {
			value := tagValue
			if strings.ContainsAny(fmt.Sprint(value), "[]") {
				value = strings.Join(strings.Fields(fmt.Sprint(value)), ", ")
			}

			allTags = append(allTags, fmt.Sprintf("tags.%s", tagFullName), value)
		} else {
			allTags = append(allTags, parseTags(tagFullName, subTags)...)
		}
	}
	return allTags
}

// getTracerAndSpanID Get the tracer and span id
func getTracerAndSpanID(ctx context.Context) (traceID uint64, spanID uint64) {
	if ctx == nil {
		return traceID, spanID
	}

	defer func() {
		_ = recover()
	}()

	if span, ok := tracer.SpanFromContext(ctx); ok {
		return span.Context().TraceID(), span.Context().SpanID()
	}

	return 0, 0
}

// Info Prints using the info level
func (l *Logger) Info(ctx context.Context, message string, opts ...Options) {
	if len(opts) > 0 {
		for _, op := range opts {
			op(ctx, message, l, logInfo)
		}
	} else {
		l.infoLog(ctx, message)
	}
}

func (l *Logger) infoLog(ctx context.Context, message string, tags ...interface{}) {
	if logEnvironment == goenv.Test || logEnvironment == goenv.Local {
		if len(tags) > 0 {
			l.infoLogger.Printf("- %s\n%v\n", message, tags)
		} else {
			l.infoLogger.Printf("- %s\n", message)
		}
	} else {
		if traceID, spanID := getTracerAndSpanID(ctx); traceID != 0 && spanID != 0 {
			tags = append(tags, "dd.trace_id", traceID, "dd.span_id", spanID)
		}
		l.zapLogger.Sugar().Infow(message, tags...)
	}
}

// Warn Prints using the warn level
func (l *Logger) Warn(ctx context.Context, message string, opts ...Options) {
	if len(opts) > 0 {
		for _, op := range opts {
			op(ctx, message, l, logWarn)
		}
	} else {
		l.warnLog(ctx, message)
	}
}

func (l *Logger) warnLog(ctx context.Context, message string, tags ...interface{}) {
	if logEnvironment == goenv.Test || logEnvironment == goenv.Local {
		if len(tags) > 0 {
			l.warnLogger.Printf("- %s\n%v\n", message, tags)
		} else {
			l.warnLogger.Printf("- %s\n", message)
		}
	} else {
		if traceID, spanID := getTracerAndSpanID(ctx); traceID != 0 && spanID != 0 {
			tags = append(tags, "dd.trace_id", traceID, "dd.span_id", spanID)
		}
		l.zapLogger.Sugar().Warnw(message, tags...)
	}
}

// Debug Prints using the debug level
func (l *Logger) Debug(ctx context.Context, message string, opts ...Options) {
	if len(opts) > 0 {
		for _, op := range opts {
			op(ctx, message, l, logDebug)
		}
	} else {
		l.debugLog(ctx, message)
	}
}

func (l *Logger) debugLog(ctx context.Context, message string, tags ...interface{}) {
	if logEnvironment == goenv.Test || logEnvironment == goenv.Local {
		if len(tags) > 0 {
			l.debugLogger.Printf("- %s\n%v\n", message, tags)
		} else {
			l.debugLogger.Printf("- %s\n", message)
		}
	} else {
		if traceID, spanID := getTracerAndSpanID(ctx); traceID != 0 && spanID != 0 {
			tags = append(tags, "dd.trace_id", traceID, "dd.span_id", spanID)
		}
		l.zapLogger.Sugar().Debugw(message, tags...)
	}
}

// Error Prints using the error level
func (l *Logger) Error(ctx context.Context, message string, opts ...Options) {
	if len(opts) > 0 {
		for _, op := range opts {
			op(ctx, message, l, logError)
		}
	} else {
		l.errorLog(ctx, message)
	}
}

func (l *Logger) errorLog(ctx context.Context, message string, tags ...interface{}) {
	if logEnvironment == goenv.Test || logEnvironment == goenv.Local {
		if len(tags) > 0 {
			l.errorLogger.Printf("- %s\n%v\n", message, tags)
		} else {
			l.errorLogger.Printf("- %s\n", message)
		}
	} else {
		if traceID, spanID := getTracerAndSpanID(ctx); traceID != 0 && spanID != 0 {
			tags = append(tags, "dd.trace_id", traceID, "dd.span_id", spanID)
		}
		l.zapLogger.Sugar().Errorw(message, tags...)
	}
}
