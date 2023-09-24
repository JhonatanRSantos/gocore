package golog

import (
	"context"
	"io"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/JhonatanRSantos/gocore/pkg/goenv"
	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestLog(t *testing.T) {
	SetEnv(goenv.Test)
	logger = nil
	initLogOnce = sync.Once{}

	if Log() == nil {
		assert.FailNow(t, "expected non nil logger for test env")
	}

	SetEnv(goenv.Production)
	logger = nil
	initLogOnce = sync.Once{}

	if Log() == nil {
		assert.FailNow(t, "expected non nil logger for production env")
	}
}
func TestLogger(t *testing.T) {
	type testInput struct {
		ctx     context.Context
		message string
		opts    []Options
	}

	type test struct {
		name   string
		setup  func(t *testing.T) (*os.File, *os.File, *observer.ObservedLogs)
		input  testInput
		assert func(t *testing.T, input testInput, read, write *os.File, ob *observer.ObservedLogs)
	}

	tests := []test{
		// nolint:dupl
		{
			name: "should log using local info logger",
			setup: func(t *testing.T) (read *os.File, write *os.File, ob *observer.ObservedLogs) {
				SetEnv(goenv.Test)
				logger = nil
				initLogOnce = sync.Once{}
				read, write, err := os.Pipe()
				if err != nil {
					assert.FailNow(t, "failed to create pipe file. Cause: %v", err)
				}
				logWriter = write
				return read, write, nil
			},
			input: testInput{
				ctx:     context.Background(),
				message: "LOCAL INFO LOGGER",
				opts:    []Options{},
			},
			assert: func(t *testing.T, input testInput, read, write *os.File, ob *observer.ObservedLogs) {
				Log().Info(input.ctx, input.message, input.opts...)
				write.Close()
				out, err := io.ReadAll(read)
				read.Close()

				switch {
				case err != nil:
					assert.FailNow(t, "failed to read from stdout. Cause: %s", err)
				case !strings.Contains(string(out), "[  INFO  ]"):
					assert.FailNow(t, "the expected output must have: '[  INFO  ]' but got: '%s'", string(out))
				}
			},
		},
		// nolint:dupl
		{
			name: "should log using local warn logger",
			setup: func(t *testing.T) (read *os.File, write *os.File, ob *observer.ObservedLogs) {
				SetEnv(goenv.Test)
				logger = nil
				initLogOnce = sync.Once{}
				read, write, err := os.Pipe()
				if err != nil {
					assert.FailNow(t, "failed to create pipe file. Cause: %v", err)
				}
				logWriter = write
				return read, write, nil
			},
			input: testInput{
				ctx:     context.Background(),
				message: "LOCAL WARN LOGGER",
				opts:    []Options{},
			},
			assert: func(t *testing.T, input testInput, read, write *os.File, ob *observer.ObservedLogs) {
				Log().Warn(input.ctx, input.message, input.opts...)
				write.Close()
				out, err := io.ReadAll(read)
				read.Close()

				switch {
				case err != nil:
					assert.FailNow(t, "failed to read from stdout. Cause: %s", err)
				case !strings.Contains(string(out), "[  WARN  ]"):
					assert.FailNow(t, "the expected output must have: '[  WARN  ]' but got: '%s'", string(out))
				}
			},
		},
		// nolint:dupl
		{
			name: "should log using local debug logger",
			setup: func(t *testing.T) (read *os.File, write *os.File, ob *observer.ObservedLogs) {
				SetEnv(goenv.Test)
				logger = nil
				initLogOnce = sync.Once{}
				read, write, err := os.Pipe()
				if err != nil {
					assert.FailNow(t, "failed to create pipe file. Cause: %v", err)
				}
				logWriter = write
				return read, write, nil
			},
			input: testInput{
				ctx:     context.Background(),
				message: "LOCAL DEBUG LOGGER",
				opts:    []Options{},
			},
			assert: func(t *testing.T, input testInput, read, write *os.File, ob *observer.ObservedLogs) {
				Log().Debug(input.ctx, input.message, input.opts...)
				write.Close()
				out, err := io.ReadAll(read)
				read.Close()

				switch {
				case err != nil:
					assert.FailNow(t, "failed to read from stdout. Cause: %s", err)
				case !strings.Contains(string(out), "[  DEBUG ]"):
					assert.FailNow(t, "the expected output must have: '[  DEBUG ]' but got: '%s'", string(out))
				}
			},
		},
		// nolint:dupl
		{
			name: "should log using local error logger",
			setup: func(t *testing.T) (read *os.File, write *os.File, ob *observer.ObservedLogs) {
				SetEnv(goenv.Test)
				logger = nil
				initLogOnce = sync.Once{}
				read, write, err := os.Pipe()
				if err != nil {
					assert.FailNow(t, "failed to create pipe file. Cause: %v", err)
				}
				logWriter = write
				return read, write, nil
			},
			input: testInput{
				ctx:     context.Background(),
				message: "LOCAL ERROR LOGGER",
				opts:    []Options{},
			},
			assert: func(t *testing.T, input testInput, read, write *os.File, ob *observer.ObservedLogs) {
				Log().Error(input.ctx, input.message, input.opts...)
				write.Close()
				out, err := io.ReadAll(read)
				read.Close()

				switch {
				case err != nil:
					assert.FailNow(t, "failed to read from stdout. Cause: %s", err)
				case !strings.Contains(string(out), "[  ERROR ]"):
					assert.FailNow(t, "the expected output must have: '[  ERROR ]' but got: '%s'", string(out))
				}
			},
		},
		// nolint:dupl
		{
			name: "should log using production info logger",
			setup: func(t *testing.T) (read *os.File, write *os.File, ob *observer.ObservedLogs) {
				SetEnv(goenv.Production)
				logger = &Logger{}
				zl, ob := observer.New(zap.InfoLevel)
				logger.zapLogger = zap.New(zl)
				zap.ReplaceGlobals(logger.zapLogger.Sugar().Desugar())
				return nil, nil, ob
			},
			input: testInput{
				ctx:     context.Background(),
				message: "PRODUCTION INFO LOGGER",
				opts: []Options{WithTags(map[string]interface{}{
					"key": "value",
				})},
			},
			assert: func(t *testing.T, input testInput, read, write *os.File, ob *observer.ObservedLogs) {
				Log().Info(input.ctx, input.message, input.opts...)

				logs := ob.All()
				amountOfLogs := ob.Len()
				expectedTags := map[string]interface{}{
					"tags.key": "value",
				}

				switch {
				case amountOfLogs > 1:
					assert.FailNow(t, "invalid amount of logs. Expected 1, but got %d", amountOfLogs)
				case logs[0].Level != zap.InfoLevel:
					assert.FailNow(t, "invalid log level. Expected info, but got %s", logs[0].Level)
				case logs[0].Message != input.message:
					assert.FailNow(t, "invalid log message. Expected %s, but got %s", input.message, logs[0].Message)
				}

				for k, v := range logs[0].ContextMap() {
					if value, ok := expectedTags[k]; !ok || value != v {
						assert.FailNow(t, "the expected tags doesn't match")
					}
				}
			},
		},
		// nolint:dupl
		{
			name: "should log using production warn logger",
			setup: func(t *testing.T) (read *os.File, write *os.File, ob *observer.ObservedLogs) {
				SetEnv(goenv.Production)
				logger = &Logger{}
				zl, ob := observer.New(zap.WarnLevel)
				logger.zapLogger = zap.New(zl)
				zap.ReplaceGlobals(logger.zapLogger.Sugar().Desugar())
				return nil, nil, ob
			},
			input: testInput{
				ctx:     context.Background(),
				message: "PRODUCTION WARN LOGGER",
				opts: []Options{WithTags(map[string]interface{}{
					"key": "value",
				})},
			},
			assert: func(t *testing.T, input testInput, read, write *os.File, ob *observer.ObservedLogs) {
				Log().Warn(input.ctx, input.message, input.opts...)

				logs := ob.All()
				amountOfLogs := ob.Len()
				expectedTags := map[string]interface{}{
					"tags.key": "value",
				}

				switch {
				case amountOfLogs > 1:
					assert.FailNow(t, "invalid amount of logs. Expected 1, but got %d", amountOfLogs)
				case logs[0].Level != zap.WarnLevel:
					assert.FailNow(t, "invalid log level. Expected info, but got %s", logs[0].Level)
				case logs[0].Message != input.message:
					assert.FailNow(t, "invalid log message. Expected %s, but got %s", input.message, logs[0].Message)
				}

				for k, v := range logs[0].ContextMap() {
					if value, ok := expectedTags[k]; !ok || value != v {
						assert.FailNow(t, "the expected tags doesn't match")
					}
				}
			},
		},
		// nolint:dupl
		{
			name: "should log using production debug logger",
			setup: func(t *testing.T) (read *os.File, write *os.File, ob *observer.ObservedLogs) {
				SetEnv(goenv.Production)
				logger = &Logger{}
				zl, ob := observer.New(zap.DebugLevel)
				logger.zapLogger = zap.New(zl)
				zap.ReplaceGlobals(logger.zapLogger.Sugar().Desugar())
				return nil, nil, ob
			},
			input: testInput{
				ctx:     context.Background(),
				message: "PRODUCTION DEBUG LOGGER",
				opts: []Options{WithTags(map[string]interface{}{
					"key": "value",
				})},
			},
			assert: func(t *testing.T, input testInput, read, write *os.File, ob *observer.ObservedLogs) {
				Log().Debug(input.ctx, input.message, input.opts...)

				logs := ob.All()
				amountOfLogs := ob.Len()
				expectedTags := map[string]interface{}{
					"tags.key": "value",
				}

				switch {
				case amountOfLogs > 1:
					assert.FailNow(t, "invalid amount of logs. Expected 1, but got %d", amountOfLogs)
				case logs[0].Level != zap.DebugLevel:
					assert.FailNow(t, "invalid log level. Expected info, but got %s", logs[0].Level)
				case logs[0].Message != input.message:
					assert.FailNow(t, "invalid log message. Expected %s, but got %s", input.message, logs[0].Message)
				}

				for k, v := range logs[0].ContextMap() {
					if value, ok := expectedTags[k]; !ok || value != v {
						assert.FailNow(t, "the expected tags doesn't match")
					}
				}
			},
		},
		// nolint:dupl
		{
			name: "should log using production error logger",
			setup: func(t *testing.T) (read *os.File, write *os.File, ob *observer.ObservedLogs) {
				SetEnv(goenv.Production)
				logger = &Logger{}
				zl, ob := observer.New(zap.ErrorLevel)
				logger.zapLogger = zap.New(zl)
				zap.ReplaceGlobals(logger.zapLogger.Sugar().Desugar())
				return nil, nil, ob
			},
			input: testInput{
				ctx:     context.Background(),
				message: "PRODUCTION ERROR LOGGER",
				opts: []Options{WithTags(map[string]interface{}{
					"key": "value",
				})},
			},
			assert: func(t *testing.T, input testInput, read, write *os.File, ob *observer.ObservedLogs) {
				Log().Error(input.ctx, input.message, input.opts...)

				logs := ob.All()
				amountOfLogs := ob.Len()
				expectedTags := map[string]interface{}{
					"tags.key": "value",
				}

				switch {
				case amountOfLogs > 1:
					assert.FailNow(t, "invalid amount of logs. Expected 1, but got %d", amountOfLogs)
				case logs[0].Level != zap.ErrorLevel:
					assert.FailNow(t, "invalid log level. Expected info, but got %s", logs[0].Level)
				case logs[0].Message != input.message:
					assert.FailNow(t, "invalid log message. Expected %s, but got %s", input.message, logs[0].Message)
				}

				for k, v := range logs[0].ContextMap() {
					if value, ok := expectedTags[k]; !ok || value != v {
						assert.FailNow(t, "the expected tags doesn't match")
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			read, write, ob := tt.setup(t)
			tt.assert(t, tt.input, read, write, ob)
		})
	}
}
