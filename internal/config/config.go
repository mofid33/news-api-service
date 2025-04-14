package config

type Config struct {
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}
	JWT struct {
		Secret string
	}
	Server struct {
		Port string
	}
}

func Load() (*Config, error) {
	cfg := &Config{}

	// Default values
	cfg.DB.Host = "localhost"
	cfg.DB.Port = "5432"
	cfg.DB.User = "postgres"
	cfg.DB.Password = "postgres"
	cfg.DB.Name = "news_api"

	cfg.JWT.Secret = "your_secret_key"

	cfg.Server.Port = "8080"

	// Here you would typically load from environment variables or config file
	// For example, using github.com/spf13/viper

	return cfg, nil
}
