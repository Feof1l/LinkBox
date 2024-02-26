package mysql

import (
	"database/sql"

	"github.com/Feof1l/LinkBox/pkg/models"
)

// Определяем тип который обертывает пул подключения sql.DB
type LinkModel struct {
	DB *sql.DB
}

// Метод для создания новой заметки в базе дынных.
func (m *LinkModel) Insert(title, content, expires string) (int, error) {
	return 0, nil
}

// Get - Метод для возвращения данных заметки по её идентификатору ID.
func (m *LinkModel) Get(id int) (*models.Link, error) {
	return nil, nil
}

// Latest - Метод возвращает 10 наиболее часто используемые заметки.
func (m *LinkModel) Latest() ([]*models.Link, error) {
	return nil, nil
}
