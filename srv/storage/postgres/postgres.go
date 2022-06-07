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

// AddPost создает статью и проверяет, если статья с таким title уже существует
func (s *Storage) AddPost(p storage.Post) error {
	rows, err := s.db.Query(context.Background(), `
		INSERT INTO posts (title, content, pubTime, link)
       	SELECT $1, $2, $3, $4
       	WHERE NOT EXISTS (SELECT 1 FROM posts WHERE title=$1);
	`,
		p.Title, p.Content, p.PubTime, p.Link,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	// ВАЖНО не забыть проверить rows.Err()
	return rows.Err()
}

// DeletePost удаляет статью по id.
func (s *Storage) DeletePost(p storage.Post) error {
	rows, err := s.db.Query(context.Background(), `
		DELETE FROM posts
		WHERE posts.id = $1;
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

// GetAllItems возвращает статьи, отсортированные по времени создания, в количестве = n.
func (s *Storage) GetAllItems() ([]storage.LocationItem, error) {
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
