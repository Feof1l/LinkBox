package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {

	//флаг командной строки, значение по умолчанию: 8080 + справка по флагу
	addr := flag.String("addr", ":8080", "Сетевой адрес HTTP")
	// извлекаем флаг из командной строки
	flag.Parse()

	// создаем новый логер для вывода информационных сообщенйи в поток stdout c припиской info
	infoLOg := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	//аналогично для логов с ошибками, такеж включим вывод фйла и номера  строки, где произошла ошибка
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Инициализируем новую структуру с зависимостями приложения.
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLOg,
	}

	//Инициализируем новую структуру http.Server, устанавливаем поля Addr и Handler,
	// поле ErrorLog, чтобы сервер использовал наш логгер вместо стандартного
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(), // Вызов нового метода app.routes()
	}

	// запуска нового веб-сервера.
	// параметры: TCP-адрес сети для прослушивания (в данном случае это "localhost:8080")
	// что любая ошибка, возвращаемая от http.ListenAndServe(), всегда non-nil.
	// flag.String() вовзращает указатель, поэтому нам нужно убрать ссылку
	infoLOg.Printf("Запуск веб-сервера на %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)

}

// структура `application` для хранения зависимостей всего веб-приложения.
// внедряем паттерн Dependency Injection - В целом рекомендуется внедрять зависимости в обработчики.
// Это делает код более явным, менее подверженным ошибкам  более простым для модульного тестирования, чем в случае использования глобальных переменных.
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
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
