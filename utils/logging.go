package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger - глобальный логгер
var Logger = logrus.New()

func init() {
	// Установка формата JSON для логов
	Logger.SetFormatter(&logrus.JSONFormatter{})

	// Получение текущего рабочего каталога
	cwd, err := os.Getwd()
	if err != nil {
		Logger.WithField("error", err).Error("Failed to get current working directory")
	} else {
		Logger.WithField("cwd", cwd).Info("Current working directory")
	}

	// Открытие файла для записи логов
	filePath := "app.log" // Укажите абсолютный путь, если нужно
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// Если файл не открылся, логи пишутся в консоль
		Logger.SetOutput(os.Stdout)
		Logger.WithField("error", err).Warn("Failed to log to file, using default stderr")
	} else {
		// Если файл успешно открылся, логи пишутся в файл
		Logger.SetOutput(file)
		Logger.WithField("filePath", filePath).Info("Logging to file initialized")
	}
}
