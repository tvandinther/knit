package logging

import (
	"log"
	"os"
	"sync"
)

var logger *log.Logger
var once sync.Once

func GetInstance() *log.Logger {
	once.Do(func() {
		logger = log.New(os.Stdout, "", 0)
	})
	return logger
}
