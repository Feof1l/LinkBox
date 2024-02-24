package main

import (
	"log"
	"net/http"
)

func main() {
	// инициализация нового рoутера,функцию "home" регистрируется как обработчик для URL-шаблона "/".
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/link", showLink)
	mux.HandleFunc("/link/create", createLink)

	//  обработчик HTTP-запросов к статическим файлам из папки "./ui/static"
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// обработчика для всех запросов, которые начинаются с "/static/". Мы убираем
	// префикс "/static" перед тем как запрос достигнет http.FileServer
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	// запуска нового веб-сервера.
	// параметры: TCP-адрес сети для прослушивания (в данном случае это "localhost:8080")
	// что любая ошибка, возвращаемая от http.ListenAndServe(), всегда non-nil.
	log.Println("Запуск веб-сервера на http://127.0.0.1:8080")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)

}
