package logger_helper

import (
	"log"
	"os"
)

func PrepareLogFile() (*os.File, error) {
	logFile, err := os.OpenFile("logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Println("Помилка при спробі відкрити лог файл: %v", err)
		return nil, err
	}
	log.SetOutput(logFile)
	return logFile, nil
}

func SoftLogClose(logFile *os.File) {
	if err := logFile.Close(); err != nil {
		log.Println("Помилка при спробі закрити лог файл: %v", err)
	}
}
