package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Создается функция-обработчик "home", которая записывает байтовый слайс, содержащий
// обработчик главной страницы
func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	//Инициализируем срез содержащий пути к двум файлам.
	// файл home.page.tmpl должен быть *первым* файлом в срезе.
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	// читаем файлы шаблона
	// если возникает ошибка, возвращаем 500 код
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	//мы используем метод Execute() для записи содержимого
	//шаблона в тело HTTP ответа. Последний параметр в Execute() предоставляет
	//возможность отправки динамических данных в шаблон.
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
	//w.Write([]byte("Привет из LinkBox"))
}

// Обработчик для отображения содержимого заметки
func showLink(w http.ResponseWriter, r *http.Request) {
	// извлекаем id из URL строки запроса используя метод r.URL.Query().Get().
	//Метод всегда будет возвращать значение параметра
	//в виде строки или пустую строку "", если нет совпадающего параметра;
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	// проверяем id,т к оно должно быть целым положительынм числом
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Отображение выбранной заметки с ID %d...", id)
	//w.Write([]byte("Отображение заметки ..."))
}

// Обработчик для создание новой заметки
func createLink(w http.ResponseWriter, r *http.Request) {
	//Создание новой заметки в базе данных является не идемпотентным действием,
	//которое изменяет состояние нашего сервера. Поэтому мы должны
	// ограничить маршрут, чтобы он отвечал только на POST-запросы.
	if r.Method != http.MethodPost {
		// для каждого ответа 405 «метод запрещен»,
		//чтобы пользователь знал, какие HTTP-методы поддерживаются для определенного URL.
		// указываем доствупные методы
		w.Header().Set("Allow", http.MethodPost)

		/*
			w.WriteHeader(405)
			w.Write([]byte("Доступен только POST-метод!"))
		*/

		http.Error(w, "Метод запрещен!", 405)
		return
	}
	w.Write([]byte("Форма для создания новой заметки ..."))
}
