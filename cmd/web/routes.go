package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	// инициализация нового рoутера,функцию "home" регистрируется как обработчик для URL-шаблона "/".
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/link", app.showLink)
	mux.HandleFunc("/link/create", app.createLink)

	//  обработчик HTTP-запросов к статическим файлам из папки "./ui/static"
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// обработчика для всех запросов, которые начинаются с "/static/". Мы убираем
	// префикс "/static" перед тем как запрос достигнет http.FileServer

	//mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
