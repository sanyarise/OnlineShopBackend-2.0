package config

import (
	"flag"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	DSN        string `toml:"dsn" env:"DSN" envDefault:"postgres://shopteam:123@localhost:5432/shop?sslmode=disable"`
	Port       string `toml:"port" env:"PORT" envDefault:":8000"`
	Server_URL string `toml:"server_url" env:"SERVER_URL" envDefault:"http://localhost:8000"`
	Logger     *zap.Logger
	LogLevel   string `toml:"log_level" env:"LOG_LEVEL" envDefault:"debug"`
}

// InitConfig() initializes the configuration
func InitConfig() (*Config, error) {
	var configPath string

	// The flag allows to specify the path to the folder with the configuration file in .toml format.
	flag.StringVar(&configPath, "config-path", "", "path to file in .toml format")
	flag.Parse()

	var cfg = Config{}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("can't load environment variables: %s", err)
	}
	if configPath != "" {
		_, err := toml.DecodeFile(configPath, &cfg)
		if err != nil {
			log.Fatalf("can't load configuration file: %s", err)
		}
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
	return &cfg, nil
}
