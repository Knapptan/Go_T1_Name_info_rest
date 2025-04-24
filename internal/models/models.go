package models

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
}

type Config struct {
	Port     string
	DBConfig DBConfig
}
