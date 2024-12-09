package logger

import (
	"fmt"
	"os"

	"github.com/olivere/elastic/v7"
	"github.com/praction-networks/quantum-ISP365/webapp/src/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func SetupLogger() error {

	cfg, err := config.LoggerEnvGet()

	if err != nil {
		return fmt.Errorf("failed to initialize logger env config: %w", err)
	}

	var level zapcore.Level
	switch cfg.LogLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "fatal":
		level = zap.FatalLevel
	case "panic":
		level = zap.PanicLevel
	case "dpanic":
		level = zap.DPanicLevel
	default:
		level = zap.InfoLevel // Fallback to info level if the log level is unrecognized
	}

	// Create encoder configuration
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // Human-readable time format

	var cores []zapcore.Core

	// Set console logging if enabled
	// Set console logging if enabled
	if cfg.Console {
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level)
		cores = append(cores, consoleCore)
	}

	// Set JSON logging to stderr if console logging is not enabled
	if !cfg.Console {
		jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)
		jsonCore := zapcore.NewCore(jsonEncoder, zapcore.AddSync(os.Stderr), level)
		cores = append(cores, jsonCore)
	}
	// ELK logging if Setup

	if cfg.ElkEnabled {
		elkURL := "http://" + cfg.ElkHost + ":" + cfg.ElkPort

		client, err := elastic.NewClient(
			elastic.SetURL(elkURL),
			elastic.SetSniff(false),
		)

		if err != nil {
			return fmt.Errorf("error creating Elasticsearch client: %v", err)
		}

		// Assuming you have a custom writer to send logs to ELK
		elkWriter := getElkWriter(client, cfg.ElkSearchIndex)

		elkCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(elkWriter), // Custom writer that sends logs to ELK
			level,
		)

		// Combine console and ELK cores
		cores = append(cores, elkCore)
	}

	var fields []zap.Field

	if deploymentName := os.Getenv("DEPLOYMENT_NAME"); deploymentName != "" {
		fields = append(fields, zap.String("deployment_name", deploymentName))
	}
	if podName := os.Getenv("POD_NAME"); podName != "" {
		fields = append(fields, zap.String("pod_name", podName))
	}
	if namespace := os.Getenv("NAMESPACE"); namespace != "" {
		fields = append(fields, zap.String("namespace", namespace))
	}

	combinedCore := zapcore.NewTee(cores...)
	log = zap.New(combinedCore, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))

	// Initialize the logger with added fields
	log = zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel), zap.Fields(fields...))

	return nil
}

// Sync flushes the logger before shutdown
func Sync() {
	if log != nil {
		_ = log.Sync() // Ensure buffered logs are flushed
	}
}

// Helper function to add fields in key-value pairs
func withFields(args []interface{}) []zap.Field {
	fields := []zap.Field{}

	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			key, ok := args[i].(string)
			if !ok {
				continue // Skip invalid keys that aren't strings
			}
			fields = append(fields, zap.Any(key, args[i+1]))
		}
	}

	return fields

}
