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
	Port        string
	Environment string // dev/prod
	DBConfig    DBConfig
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
	Name       string  `json:"name" example:"Ivan"`
	Surname    string  `json:"surname" example:"Ivanov"`
	Patronymic *string `json:"patronymic,omitempty" example:"Ivanovich"`
}

type PersonEnriched struct {
	ID          int       `json:"id" example:"1"`
	Name        string    `json:"name" example:"Ivan"`
	Surname     string    `json:"surname" example:"Ivanov"`
	Patronymic  *string   `json:"patronymic,omitempty" example:"Ivanovich"`
	Age         int       `json:"age" example:"30"`
	Gender      string    `json:"gender" example:"male"`
	Nationality string    `json:"nationality" example:"RU"`
	CreatedAt   time.Time `json:"created_at" example:"2023-10-23T12:34:56Z`
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
