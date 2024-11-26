// Модуль для объявления конфигурации агента
package flags

import (
	"flag"
)

var (
	Action   string
	DataType string
	Name     string
	Username string
	Password string
	Meta     string
)

// ParseFlags читает аргументы переданные при старте агента
func ParseFlags() {
	flag.StringVar(&Action, "a", "", "Action")
	flag.StringVar(&DataType, "t", "cred", "Service type")
	flag.StringVar(&Name, "n", "", "Name")
	flag.StringVar(&Username, "u", "", "Username")
	flag.StringVar(&Password, "p", "", "Password")
	flag.StringVar(&Meta, "m", "", "Meta")
	flag.Parse()
}
