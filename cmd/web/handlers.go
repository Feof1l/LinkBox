package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// Создается функция-обработчик "home", которая записывает байтовый слайс, содержащий
// обработчик главной страницы
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
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
		//app.errorLog.Println(err.Error())
		app.serverError(w, err)
		return
	}
	//мы используем метод Execute() для записи содержимого
	//шаблона в тело HTTP ответа. Последний параметр в Execute() предоставляет
	//возможность отправки динамических данных в шаблон.
	err = ts.Execute(w, nil)
	if err != nil {
		//app.errorLog.Println(err.Error())
		app.serverError(w, err)
	}
	//w.Write([]byte("Привет из LinkBox"))
}

// Обработчик для отображения содержимого заметки
func (app *application) showLink(w http.ResponseWriter, r *http.Request) {
	// извлекаем id из URL строки запроса используя метод r.URL.Query().Get().
	//Метод всегда будет возвращать значение параметра
	//в виде строки или пустую строку "", если нет совпадающего параметра;
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	// проверяем id,т к оно должно быть целым положительынм числом
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	fmt.Fprintf(w, "Отображение выбранной заметки с ID %d...", id)
	//w.Write([]byte("Отображение заметки ..."))
}

// Обработчик для создание новой заметки
func (app *application) createLink(w http.ResponseWriter, r *http.Request) {
	//Создание новой заметки в базе данных является не идемпотентным действием,
	//которое изменяет состояние нашего сервера. Поэтому мы должны
	// ограничить маршрут, чтобы он отвечал только на POST-запросы.
	if r.Method != http.MethodPost {
		// для каждого ответа 405 «метод запрещен»,
		//чтобы пользователь знал, какие HTTP-методы поддерживаются для определенного URL.
		// указываем доствупные методы
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	title := "История про улитку"
	content := "Улитка выползла из раковины,\nвытянула рожки,\nи опять подобрала их."
	expires := "7"
	id, err := app.links.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/link?id=%d", id), http.StatusSeeOther)
	//w.Write([]byte("Форма для создания новой заметки ..."))
}
