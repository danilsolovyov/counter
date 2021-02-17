package main

import (
	"context"
	"fmt"
	//"log"
	"net/http"
	"os"
	"time"
)

// Chronometer - Меняет значение counter.count
func Chronometer(c *Connection) {
    for range time.Tick(time.Second) {
        c.UpdateScore()
    }
}

func main() {
    // Подключаемся к MongoDB
    connection, _ := ConnectToDb()

    // Задаем параметры веб-сервера
    port := os.Getenv("PORT")
    if port == "" {
        port = "3000"
    }

    mux := http.NewServeMux()

    // Задаем endpoint "/subscribe"
    mux.HandleFunc("/subscribe", func(rw http.ResponseWriter, r *http.Request) {
        rw.Header().Set("Content-Type", "text/event-stream")
        rw.Header().Set("Cache-Control", "no-cache")
        rw.Header().Set("Connection", "keep-alive")
        rw.Header().Set("Access-Control-Allow-Origin", "*")
        response := uint32(connection.GetScore().Count)
        fmt.Fprint(rw, response)
    })

    go http.ListenAndServe(":"+port, mux)

    // Запускаем счетчик
    Chronometer(connection)

    // Завершаем работу
    connection.Disconnect(context.TODO())
}
