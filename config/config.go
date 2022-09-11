package config

import (
	"flag"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	DSN        string `toml:"dsn"`
	Port       string `toml:"port"`
	Server_URL string `toml:"server_url"`
	Logger     *zap.Logger
	LogLevel   string `toml:"log_level"`
}

// InitConfig() initializes the configuration
func InitConfig() (*Config, error) {
	var configPath string

	// The flag allows to specify the path to the folder with the configuration file.
	// When running without a flag, the default path is used
	flag.StringVar(&configPath, "config-path", "./config/config.toml", "path to file in .toml format")
	flag.Parse()

	var cfg = Config{}

	_, err := toml.DecodeFile(configPath, &cfg)
	if err != nil {
		log.Fatalf("can't load configuration file: %s", err)
	}
	// Logger settings
	atomicLevel := zap.NewAtomicLevel()

	switch cfg.LogLevel {
	case "info":
		{
			atomicLevel.SetLevel(zap.InfoLevel)
		}
	case "warning":
		{
			atomicLevel.SetLevel(zap.WarnLevel)
		}
	case "debug":
		{
			atomicLevel.SetLevel(zap.DebugLevel)
		}
	case "error":
		{
			atomicLevel.SetLevel(zap.ErrorLevel)

		}
	case "panic":
		{
			atomicLevel.SetLevel(zap.PanicLevel)
		}
	case "fatal":
		{
			atomicLevel.SetLevel(zap.FatalLevel)
		}
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atomicLevel,
	), zap.AddCaller())
	cfg.Logger = logger
	return &cfg, err
}
