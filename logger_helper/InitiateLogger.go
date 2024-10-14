package logger_helper

import (
	"log"
	"os"
)

func SoftPrepareLogFile() *os.File {
	logFile, err := os.OpenFile("logs/logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Println("Помилка при спробі відкрити лог файл: %v", err)
		return nil
	}
	log.SetOutput(logFile)
	return logFile
}

func SoftLogClose(logFile *os.File) {
	if logFile == nil {
		log.Println("Лог файлу не існує")
	} else if err := logFile.Close(); err != nil {
		log.Println("Помилка при спробі закрити лог файл: %v", err)
	}
}
