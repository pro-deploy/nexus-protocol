package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Auth     AuthConfig     `mapstructure:"auth"`
	AI       AIConfig       `mapstructure:"ai"`
	Metrics  MetricsConfig  `mapstructure:"metrics"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
	Cache    CacheConfig    `mapstructure:"cache"`
}

type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
	MaxConns int    `mapstructure:"max_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type AuthConfig struct {
	Provider      string        `mapstructure:"provider"` // jwt or keycloak
	JWTSecret     string        `mapstructure:"jwt_secret"`
	JWTExpiry     time.Duration `mapstructure:"jwt_expiry"`
	RefreshExpiry time.Duration `mapstructure:"refresh_expiry"`

	// Keycloak configuration
	KeycloakURL      string `mapstructure:"keycloak_url"`
	KeycloakRealm    string `mapstructure:"keycloak_realm"`
	KeycloakClientID string `mapstructure:"keycloak_client_id"`
	KeycloakSecret   string `mapstructure:"keycloak_secret"`
}

type AIConfig struct {
	Provider        string `mapstructure:"provider"`
	APIKey          string `mapstructure:"api_key"`
	BaseURL         string `mapstructure:"base_url"`
	Model           string `mapstructure:"model"`
	MaxTokens       int    `mapstructure:"max_tokens"`
	Temperature     float32 `mapstructure:"temperature"`
	Timeout         time.Duration `mapstructure:"timeout"`
}

type MetricsConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Port    int    `mapstructure:"port"`
	Path    string `mapstructure:"path"`
}

type RateLimitConfig struct {
	Enabled         bool          `mapstructure:"enabled"`
	RequestsPerMin  int           `mapstructure:"requests_per_min"`
	BurstSize       int           `mapstructure:"burst_size"`
	CleanupInterval time.Duration `mapstructure:"cleanup_interval"`
}

type CacheConfig struct {
	Enabled bool          `mapstructure:"enabled"`
	TTL     time.Duration `mapstructure:"ttl"`
	MaxSize int           `mapstructure:"max_size"`
}

// Load reads configuration from environment variables and config files
func Load() (*Config, error) {
	// Set defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")
	viper.SetDefault("server.idle_timeout", "120s")

	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.max_conns", 25)

	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)

	viper.SetDefault("auth.jwt_expiry", "24h")
	viper.SetDefault("auth.refresh_expiry", "168h") // 7 days

	viper.SetDefault("ai.provider", "openai")
	viper.SetDefault("ai.model", "gpt-4")
	viper.SetDefault("ai.max_tokens", 4096)
	viper.SetDefault("ai.temperature", 0.7)
	viper.SetDefault("ai.timeout", "30s")

	viper.SetDefault("metrics.enabled", true)
	viper.SetDefault("metrics.port", 9090)
	viper.SetDefault("metrics.path", "/metrics")

	viper.SetDefault("rate_limit.enabled", true)
	viper.SetDefault("rate_limit.requests_per_min", 1000)
	viper.SetDefault("rate_limit.burst_size", 100)
	viper.SetDefault("rate_limit.cleanup_interval", "1m")

	viper.SetDefault("cache.enabled", true)
	viper.SetDefault("cache.ttl", "5m")
	viper.SetDefault("cache.max_size", 10000)

	// Environment variables
	viper.SetEnvPrefix("NEXUS")
	viper.AutomaticEnv()

	// Config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		// Config file is optional, continue with defaults
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate required fields
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &config, nil
}

func validateConfig(cfg *Config) error {
	if cfg.Auth.JWTSecret == "" {
		return fmt.Errorf("auth.jwt_secret is required")
	}

	if cfg.Database.User == "" {
		return fmt.Errorf("database.user is required")
	}

	if cfg.Database.DBName == "" {
		return fmt.Errorf("database.dbname is required")
	}

	if cfg.AI.APIKey == "" {
		return fmt.Errorf("ai.api_key is required")
	}

	return nil
}

// GetDSN returns PostgreSQL DSN string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// GetRedisAddr returns Redis address string
func (c *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
