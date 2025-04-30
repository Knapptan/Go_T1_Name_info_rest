package models

import (
	"fmt"
	"time"
)

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

func (cfg *Config) DatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DBConfig.User,
		cfg.DBConfig.Password,
		cfg.DBConfig.Host,
		cfg.DBConfig.Port,
		cfg.DBConfig.DbName,
	)
}

type PersonInput struct {
	Name       string  `json:"name"`
	Surname    string  `json:"surname"`
	Patronymic *string `json:"patronymic,omitempty"`
}

type PersonEnriched struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Surname     string    `json:"surname"`
	Patronymic  *string   `json:"patronymic,omitempty"`
	Age         int       `json:"age"`
	Gender      string    `json:"gender"`
	Nationality string    `json:"nationality"`
	CreatedAt   time.Time `json:"created_at"`
}

type Filters struct {
	Name   string
	Gender string
	Limit  int
	Offset int
}

type PersonUpdate struct {
	Name        *string `json:"name,omitempty"`
	Surname     *string `json:"surname,omitempty"`
	Patronymic  *string `json:"patronymic,omitempty"`
	Age         *int    `json:"age,omitempty"`
	Gender      *string `json:"gender,omitempty"`
	Nationality *string `json:"nationality,omitempty"`
}

type Enrichment struct {
	Age         int
	Gender      string
	Nationality string
}
