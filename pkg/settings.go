package pkg

import "log"

const (
	DateFormatTemplateYYYYMMDD   = "20060102"
	DateFormatTemplateDD_MM_YYYY = "02.01.2006"
	MaxNumTasks                  = 50
)

var (
	Logger *log.Logger
	JwtKey = []byte("your-secret-key-here")
)
