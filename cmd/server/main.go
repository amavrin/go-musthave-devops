package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		defer r.Body.Close() // Закрываем тело запроса после чтения

		fmt.Println("Получен запрос:", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			log.Print("error writing response:", err)
		}
	} else {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/", handler) // Регистрируем обработчик для корневого пути

	fmt.Println("Сервер запущен на порту 8080...")
	err := http.ListenAndServe(":8080", nil) // Запускаем сервер на порту 8080
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}
