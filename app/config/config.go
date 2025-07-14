package config

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

type CacheConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DB       int
	PoolSize int
}

type DbConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

func (d *DbConfig) DSN() string {
	var builder strings.Builder
	appender := func(key, value string) {
		if value == "" {
			return
		}
		builder.WriteString(fmt.Sprintf("%s=%s ", key, value))
	}
	if d.Host == "" {
		panic("DB host not set")
	}
	appender("host", d.Host)
	if d.Port == "" {
		panic("DB port not set")
	}
	appender("port", d.Port)
	appender("user", d.User)
	appender("password", d.Password)
	appender("database", d.Name)
	appender("sslmode", d.SSLMode)
	return builder.String()
}

type JWTConfig struct {
	Secret                 []byte
	TokenExpiration        time.Duration
	RefreshTokenExpiration time.Duration
}

type ServerConfig struct {
	Host     string
	Port     string
	LogLevel *slog.LevelVar
}

func (s *ServerConfig) ServeAddr() string {
	return fmt.Sprintf(":%s", s.Port)
}

type AppConfig struct {
	JWTConfig
	ServerConfig
	DbConfig
	CacheConfig
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvAsInt64(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}

func getEnvLevelVar(key string, fallback slog.Level) *slog.LevelVar {
	var levelVar slog.LevelVar
	level := fallback
	if value, ok := os.LookupEnv(key); ok {
		level.UnmarshalText([]byte(value))
	}
	levelVar.Set(level)

	return &levelVar
}

func jwtFromEnv() JWTConfig {
	return JWTConfig{
		Secret:                 []byte(getEnv("JWT_SECRET", "MY_SECRET_ENV_TOKEN")),
		TokenExpiration:        time.Second * time.Duration(getEnvAsInt64("JWT_TOKEN_EXPIRATION", 3600*12)),
		RefreshTokenExpiration: time.Second * time.Duration(getEnvAsInt64("JWT_REFRESH_TOKEN_EXPIRATION", 3600*24)),
	}
}

// load database configuration from environment variables.
func dbFromEnv() DbConfig {
	return DbConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		Name:     getEnv("DB_NAME", "TicketApp"),
		Password: getEnv("DB_PASSWORD", "Password123!"),
		User:     getEnv("DB_USER", "user"),
		SSLMode:  getEnv("DB_SSLMODE", "prefer"),
	}
}

// load cache configuration from environment variables.
func cacheFromEnv() CacheConfig {
	return CacheConfig{
		Host:     getEnv("CACHE_HOST", "localhost"),
		Port:     getEnv("CACHE_PORT", "5432"),
		Username: getEnv("CACHE_USERNAME", "default"),
		Password: getEnv("CACHE_PASSWORD", "redispw"),
		DB:       int(getEnvAsInt64("CACHE_DB", 0)),
		PoolSize: int(getEnvAsInt64("CACHE_POOL_SIZE", 10)),
	}
}

func serverFromEnv() ServerConfig {
	return ServerConfig{
		Host:     getEnv("HOST", "http://localhost"),
		Port:     getEnv("PORT", "8080"),
		LogLevel: getEnvLevelVar("LOG", slog.LevelInfo),
	}
}

var (
	config AppConfig
	once   sync.Once
)

// Care, this function is designed to be used in a testing environment
func CustomConfig(loader func() AppConfig) AppConfig {
	once.Do(func() {
		config = loader()
	})
	return config
}

// LoadConfigFromYaml loads database configuration from yaml file
func LoadConfigFromYaml(path string) (*AppConfig, error) {
	fIn, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fIn.Close()
	fileContent, err := io.ReadAll(fIn)
	if err != nil {
		return nil, err
	}
	var cfg AppConfig
	if err = yaml.Unmarshal(fileContent, &cfg); err != nil {
		return nil, err
	}
	return &cfg, err
}

func Config() AppConfig {
	once.Do(func() {
		config = AppConfig{
			JWTConfig:    jwtFromEnv(),
			ServerConfig: serverFromEnv(),
			DbConfig:     dbFromEnv(),
			CacheConfig:  cacheFromEnv(),
		}
	})
	return config
}
