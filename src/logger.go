package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var logFile *os.File               // Файл для логирования
var filePath = "../../logfile.txt" // Путь к файлу для логирования

// Функция, которая отвечает за логгирование(в консоль и файл) запросов, приходящих на сервер
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().UTC()

		inner.ServeHTTP(w, r)

		fmt.Fprintf(
			logFile,
			"%s\t%s\t%s\t%s\t%s\n",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
			start)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)

	})
}

func LoggerInit() {
	var err error

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		logFile, err = os.Create(filePath) // создадим файл для логов, если такого еще не существует
	}

	logFile, err = os.OpenFile(filePath, os.O_APPEND | os.O_WRONLY, 0666)
	if err != nil {
		panic("LogFile is not woking!")
		return
	}
	defer logFile.Close()
}
