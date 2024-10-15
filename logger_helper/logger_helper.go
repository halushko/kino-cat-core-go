package logger_helper

import (
	"fmt"
	"log"
	"os"
)

//goland:noinspection GoUnusedExportedFunction
func SoftPrepareLogFile() *os.File {
	logFile, err := os.OpenFile("logs/logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Println(fmt.Sprintf("Помилка при спробі відкрити лог файл: %v", err))
		return nil
	}
	log.SetOutput(logFile)
	return logFile
}

//goland:noinspection GoUnusedExportedFunction
func SoftLogClose(logFile *os.File) {
	if logFile == nil {
		log.Println("SoftLogClose[] Лог файлу не існує")
	} else if err := logFile.Close(); err != nil {
		log.Println(fmt.Sprintf("[SoftLogClose] Помилка при спробі закрити лог файл: %v", err))
	}
}
