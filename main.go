package main

import (
	"log"
	"net/http"
)

func main() {
	// Используется функция http.NewServeMux() для инициализации нового рoутера, затем
	// функцию "home" регистрируется как обработчик для URL-шаблона "/".
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/link", showLink)
	mux.HandleFunc("/link/create", createLink)

	// Используется функция http.ListenAndServe() для запуска нового веб-сервера.
	// Мы передаем два параметра: TCP-адрес сети для прослушивания (в данном случае это "localhost:8080")
	// и созданный рутер. Если вызов http.ListenAndServe() возвращает ошибку
	// мы используем функцию log.Fatal() для логирования ошибок. Обратите внимание
	// что любая ошибка, возвращаемая от http.ListenAndServe(), всегда non-nil.
	log.Println("Запуск веб-сервера на http://127.0.0.1:8080")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)

}

// Создается функция-обработчик "home", которая записывает байтовый слайс, содержащий
// текст "Привет из LinkBox" как тело ответа.

// обработчик главной страницы
func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Привет из LinkBox"))
}

// Обработчик для отображения содержимого заметки
func showLink(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Отображение заметки ..."))
}

// Обработчик для создание новой заметки
func createLink(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Форма для создания новой заметки ..."))
}
