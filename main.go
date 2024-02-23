package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
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
// текст "Привет из  LinkBox" как тело ответа.

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
	w.Write([]byte("Отображение заметки ..."))
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
