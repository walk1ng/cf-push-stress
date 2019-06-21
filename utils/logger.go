package utils

import (
	"log"
	"os"
)

var logger = log.New(os.Stdout, "[info]", log.Lshortfile)
