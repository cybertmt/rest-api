package storage

// Post - публикация.
type Post struct {
	ID      int    // номер записи
	Title   string // заголовок публикации
	Content string // содержание публикации
	PubTime int64  // время публикации
	Link    string // ссылка на источник
}

// LocationItem - balloon.
type LocationItem struct {
	ID        int     `json:"id"`      // номер записи
	Title     string  `json:"title"`   // заголовок публикации
	Content   string  `json:"content"` // содержание публикации
	Link      string  `json:"link"`    // ссылка на источник
	Latitude  float64 `json:"lat"`     // широта
	Longitude float64 `json:"lon"`     // долгота
}

// Interface задаёт контракт на работу с БД.
type Interface interface {
	AddPost(Post) error                   // создание новой публикации
	DeletePost(Post) error                // удаление публикации по ID
	GetAllItems() ([]LocationItem, error) // получение публикации по ID
}
