package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/Feof1l/LinkBox/pkg/models"
)

// Создается функция-обработчик "home", которая записывает байтовый слайс, содержащий
// обработчик главной страницы
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	s, err := app.links.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	// отображаем шаблон
	app.render(w, r, "home.page.tmpl", &templateData{
		Links: s,
	})
	/*
		// Создаем экземпляр структуры templateData, содержащий срез с заметками.
		data := &templateData{Links: s}

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
		err = ts.Execute(w, data)
		if err != nil {
			//app.errorLog.Println(err.Error())
			app.serverError(w, err)
		}
			for _, link := range s {
				fmt.Fprintf(w, "%v\n", link)
			}

	*/

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

	s, err := app.links.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.render(w, r, "show.page.tmpl", &templateData{
		Link: s,
	})

	/*
		// Создаем экземпляр структуры templateData, содержащей данные заметки.
		data := &templateData{Link: s}

		files := []string{
			"./ui/html/show.page.tmpl",
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
		err = ts.Execute(w, data)
		if err != nil {
			//app.errorLog.Println(err.Error())
			app.serverError(w, err)
		}
	*/
	//fmt.Fprintf(w, "%v", s)
	//w.Write([]byte("Отображение заметки ..."))
}

type FormData struct {
	Title   string "title"
	Content string "content"
	Option  string "option"
}

// Обработчик для создание новой заметки
func (app *application) createLink(w http.ResponseWriter, r *http.Request) {

	// Инициализируем срез содержащий пути к двум файлам. Обратите внимание, что
	// файл home.page.tmpl должен быть *первым* файлом в срезе.
	files := []string{
		"./ui/html/create.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	// Используем функцию template.ParseFiles() для чтения файлов шаблона.
	// Если возникла ошибка, мы запишем детальное сообщение ошибки и
	// используя функцию http.Error() мы отправим пользователю
	// ответ: 500 Internal Server Error (Внутренняя ошибка на сервере)
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// Затем мы используем метод Execute() для записи содержимого
	// шаблона в тело HTTP ответа. Последний параметр в Execute() предоставляет
	// возможность отправки динамических данных в шаблон.
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
	//http.Redirect(w, r, "/new-page", http.StatusSeeOther)

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
	var formData FormData
	err = json.NewDecoder(r.Body).Decode(&formData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := app.links.Insert("INSERT INTO mytable (title, content) VALUES (?, ?)", formData.Title, formData.Content)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Ошибка при записи данных в базу данных", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/link?id=%d", id), http.StatusSeeOther)
	//w.Write([]byte("Форма для создания новой заметки ..."))

}
