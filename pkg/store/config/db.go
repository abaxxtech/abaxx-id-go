package config

// DBConfig holds database configuration
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewDefaultConfig returns a DBConfig with default values
func NewDefaultConfig() DBConfig {
	return DBConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "test",
		Password: "testonly",
		DBName:   "dwn_test",
		SSLMode:  "disable",
	}
}