package mysql

import (
	"database/sql"
	"errors"

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
	stmt := `SELECT id,title,content,created, expires FROM links
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// Используем метод QueryRow() для выполнения SQL запроса,
	// передавая ненадежную переменную id в качестве значения для плейсхолдера
	// Возвращается указатель на объект sql.Row, который содержит данные записи.
	row := m.DB.QueryRow(stmt, id)

	// Инициализируем указатель на новую структуру Link
	s := &models.Link{}
	// Используйте row.Scan(), чтобы скопировать значения из каждого поля от sql.Row в
	// соответствующее поле в структуре Link.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

// Latest - Метод возвращает 10 наиболее часто используемые заметки.
func (m *LinkModel) Latest() ([]*models.Link, error) {
	stmt := `SELECT id,title,content,created,expires FROM links
	WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// Откладываем вызов rows.Close(), чтобы быть уверенным, что набор результатов из sql.Rows
	// правильно закроется перед вызовом метода Latest(). Этот оператор откладывания
	// должен выполнится *после* проверки на наличие ошибки в методе Query().
	// В противном случае, если Query() вернет ошибку, это приведет к панике
	// так как он попытается закрыть набор результатов у которого значение: n
	defer rows.Close()

	// Инициализируем пустой срез для хранения объектов models.Links
	var links []*models.Link

	// Используем rows.Next() для перебора результата. Этот метод предоставляем
	// первый а затем каждую следующею запись из базы данных для обработки
	// методом rows.Scan().
	for rows.Next() {
		s := &models.Link{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		links = append(links, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return links, nil
}
