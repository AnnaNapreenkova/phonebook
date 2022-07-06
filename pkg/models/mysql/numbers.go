package mysql

import (
	"database/sql"
	"errors"

	"tutorial-go.com/phonebook/pkg/models"
)

type NumberModel struct {
	DB *sql.DB
}

func (m *NumberModel) Insert(name, phone string) (int, error) {
	stmt := `INSERT INTO numbers (name, phone, created)
	VALUES(?, ?, UTC_TIMESTAMP())`

	result, err := m.DB.Exec(stmt, name, phone)
	if err != nil {
		return 0, nil
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *NumberModel) Get(id int) (*models.Number, error) {
	stmt := `SELECT id, name, phone, created FROM numbers
	WHERE id = ?`

	row := m.DB.QueryRow(stmt, id)
	s := &models.Number{}

	err := row.Scan(&s.ID, &s.Name, &s.Phone, &s.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *NumberModel) Latest() ([]*models.Number, error) {
	stmt := `SELECT id, name, phone, created FROM numbers
	ORDER BY name ASC`

	// Используем метод Query() для выполнения нашего SQL запроса.
	// В ответ мы получим sql.Rows, который содержит результат нашего запроса.
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// Откладываем вызов rows.Close(), чтобы быть уверенным, что набор результатов из sql.Rows
	// правильно закроется перед вызовом метода Latest(). Этот оператор откладывания
	// должен выполнится *после* проверки на наличие ошибки в методе Query().
	// В противном случае, если Query() вернет ошибку, это приведет к панике
	// так как он попытается закрыть набор результатов у которого значение: nil.
	defer rows.Close()

	// Инициализируем пустой срез для хранения объектов models.Snippets.
	var numbers []*models.Number

	// Используем rows.Next() для перебора результата. Этот метод предоставляем
	// первый а затем каждую следующею запись из базы данных для обработки
	// методом rows.Scan().
	for rows.Next() {
		// Создаем указатель на новую структуру Snippet
		s := &models.Number{}
		// Используем rows.Scan(), чтобы скопировать значения полей в структуру.
		// Опять же, аргументы предоставленные в row.Scan()
		// должны быть указателями на место, куда требуется скопировать данные и
		// количество аргументов должно быть точно таким же, как количество
		// столбцов из таблицы базы данных, возвращаемых вашим SQL запросом.
		err = rows.Scan(&s.ID, &s.Name, &s.Phone, &s.Created)
		if err != nil {
			return nil, err
		}
		// Добавляем структуру в срез.
		numbers = append(numbers, s)
	}

	// Когда цикл rows.Next() завершается, вызываем метод rows.Err(), чтобы узнать
	// если в ходе работы у нас не возникла какая либо ошибка.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Если все в порядке, возвращаем срез с данными.
	return numbers, nil
}
