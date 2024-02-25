package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
)

func main() {

	//флаг командной строки, значение по умолчанию: 8080 + справка по флагу
	addr := flag.String("addr", ":8080", "Сетевой адрес HTTP")
	// извлекаем флаг из командной строки
	flag.Parse()

	// инициализация нового рoутера,функцию "home" регистрируется как обработчик для URL-шаблона "/".
	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/link", showLink)
	mux.HandleFunc("/link/create", createLink)

	//  обработчик HTTP-запросов к статическим файлам из папки "./ui/static"
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// обработчика для всех запросов, которые начинаются с "/static/". Мы убираем
	// префикс "/static" перед тем как запрос достигнет http.FileServer

	//mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// запуска нового веб-сервера.
	// параметры: TCP-адрес сети для прослушивания (в данном случае это "localhost:8080")
	// что любая ошибка, возвращаемая от http.ListenAndServe(), всегда non-nil.
	// flag.String() вовзращает указатель, поэтому нам нужно убрать ссылку
	log.Printf("Запуск веб-сервера на %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)

}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}
