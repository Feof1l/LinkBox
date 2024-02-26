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
	stmt := `INSERT INTO links (title,content,created,expires)
	VALUES (?,?,UTC_TIMESTAMP(),DATE_ADD(UTC_TIMESTAMP(),INTERVAL ? DAY))`
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}
	// Используем метод LastInsertId(), чтобы получить последний ID созданной записи из таблицу links.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// Get - Метод для возвращения данных заметки по её идентификатору ID.
func (m *LinkModel) Get(id int) (*models.Link, error) {
	return nil, nil
}

// Latest - Метод возвращает 10 наиболее часто используемые заметки.
func (m *LinkModel) Latest() ([]*models.Link, error) {
	return nil, nil
}
