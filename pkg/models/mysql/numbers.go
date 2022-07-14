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

func (m *NumberModel) Delete(id int) error {
	stmt := `DELETE FROM numbers WHERE id = ?`

	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *NumberModel) Edit(name, phone string, id int) error {
	stmt := `UPDATE numbers set name = ?, phone = ? WHERE id = ?`

	_, err := m.DB.Exec(stmt, name, phone, id)
	if err != nil {
		return err
	}

	return nil
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

func (m *NumberModel) Search(str1, str2 string) ([]*models.Number, error) {
	stmt := `SELECT id, name, phone, created FROM numbers WHERE name = ? OR phone = ?`

	rows, err := m.DB.Query(stmt, str1, str2)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var numbers []*models.Number

	for rows.Next() {
		s := &models.Number{}
		err = rows.Scan(&s.ID, &s.Name, &s.Phone, &s.Created)
		if err != nil {
			return nil, err
		}
		numbers = append(numbers, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return numbers, nil
}

func (m *NumberModel) AllRecords() ([]*models.Number, error) {
	stmt := `SELECT id, name, phone, created FROM numbers
	ORDER BY name ASC`

	// Используем метод Query() для выполнения SQL запроса.
	// В ответ получаем sql.Rows, который содержит результат нашего запроса.
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// Откладываем вызов rows.Close(), чтобы быть уверенным, что набор результатов из sql.Rows
	// правильно закроется перед вызовом метода. Этот оператор откладывания
	// должен выполнится *после* проверки на наличие ошибки в методе Query().
	// В противном случае, если Query() вернет ошибку, это приведет к панике
	// так как он попытается закрыть набор результатов у которого значение: nil.
	defer rows.Close()

	// Инициализируем пустой срез для хранения объектов models.Numbers.
	var numbers []*models.Number

	// Используем rows.Next() для перебора результата. Этот метод предоставляем
	// первый а затем каждую следующею запись из базы данных для обработки
	// методом rows.Scan().
	for rows.Next() {
		// Создаем указатель на новую структуру Number
		s := &models.Number{}
		// Используем rows.Scan(), чтобы скопировать значения полей в структуру.
		// Аргументы предоставленные в row.Scan()
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
