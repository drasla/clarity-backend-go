package fnMySQL

import (
	"fmt"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	Charset  string
}

func (c Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", c.User, c.Password, c.Host, c.Port, c.DBName, c.Charset)
}
