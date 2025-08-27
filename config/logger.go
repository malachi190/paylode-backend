package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var GeneralLogger *log.Logger

var ErrorLogger *log.Logger

func init() {
	// absPath, err := filepath.Abs("../log")

	base, _ := os.Getwd()
	logDir := filepath.Join(base, "log")

	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Println("Error creating directory:", err)
		os.Exit(1)
	}

	generalLog, err := os.OpenFile(filepath.Join(logDir, "general-log.log"),
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}

	GeneralLogger = log.New(generalLog, "General Logger:\t", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(generalLog, "Error Logger:\t", log.Ldate|log.Ltime|log.Lshortfile)
}
