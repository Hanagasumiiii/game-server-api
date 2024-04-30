package config

type Config struct {
	Server   Server
	Postgres Postgres
}

type (
	Server struct {
		Port int
	}
	Postgres struct {
		Host     string
		Port     int
		User     string
		Password string
		Database string
	}
)

func LoadConfig() (*Config, error) {
	return &Config{
		Server: Server{
			Port: 8080,
		},
		Postgres: Postgres{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "postgres",
			Database: "postgres",
		},
	}, nil
}
