// Package pkg
package pkg

import (
	"database/sql"
	"log"
)

const (
	DateFormatTemplateYYYYMMDD = "20060102"
	DateFormatTemplateDDMMYYYY = "02.01.2006"
	MaxNumTasks                = 50
)

var (
	Logger               *log.Logger
	Port                 string
	DBFile               string
	DB                   *sql.DB
	ExpectedPassword     string
	ExpectedHashPassword string
	JwtKey               = []byte("your-secret-key-here")
)
