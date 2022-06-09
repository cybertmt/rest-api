package postgres

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"restapisrv/srv/storage"
)

// Storage Хранилище данных.
type Storage struct {
	db *pgxpool.Pool
}

// New Конструктор, принимает строку подключения к БД.
func New(constr string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := Storage{
		db: db,
	}
	return &s, nil
}

// AddItem создает статью и проверяет, если статья с таким title уже существует
func (s *Storage) AddItem(p storage.LocationItem) error {
	rows, err := s.db.Query(context.Background(), `
		INSERT INTO locations (title, content, link, latitude, longitude)
       	SELECT $1, $2, $3, $4, $5
       	WHERE NOT EXISTS (SELECT 1 FROM locations WHERE title=$1);
	`,
		p.Title, p.Content, p.Link, p.Latitude, p.Longitude,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	// ВАЖНО не забыть проверить rows.Err()
	return rows.Err()
}

// DeleteItem удаляет статью по id.
func (s *Storage) DeleteItem(p storage.LocationItem) error {
	rows, err := s.db.Query(context.Background(), `
		DELETE FROM locations
		WHERE locations.id = $1;
	`,
		p.ID,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	// ВАЖНО не забыть проверить rows.Err()
	return rows.Err()
}

// Items возвращает статьи, отсортированные по времени создания, в количестве = n.
func (s *Storage) Items() ([]storage.LocationItem, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			*
		FROM locations
		ORDER BY id ASC;
	`,
	)
	if err != nil {
		return nil, err
	}
	var locations []storage.LocationItem
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var t storage.LocationItem
		err = rows.Scan(
			&t.ID,
			&t.Title,
			&t.Content,
			&t.Link,
			&t.Latitude,
			&t.Longitude,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		locations = append(locations, t)

	}
	// ВАЖНО не забыть проверить rows.Err()
	return locations, rows.Err()
}
