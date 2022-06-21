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

// StringLocationItem - balloon.
type StringLocationItem struct {
	ID        string `json:"id"`      // номер записи
	Title     string `json:"title"`   // заголовок публикации
	Content   string `json:"content"` // содержание публикации
	Link      string `json:"link"`    // ссылка на источник
	Latitude  string `json:"lat"`     // широта
	Longitude string `json:"lon"`     // долгота
}

// Interface задаёт контракт на работу с БД.
type Interface interface {
	Items() ([]LocationItem, error)             // получение всех публикаций
	StringItems() ([]StringLocationItem, error) // получение всех публикаций
	AddItem(item LocationItem) error            // создание новой публикации
	DeleteItem(item LocationItem) error         // удаление публикации по ID
	DeleteAllItem() error                       // удаление всех публикаций
	SortedItems(item LocationItem) ([]LocationItem, error)
	//PostById(Post) ([]Post, error)  // получение публикации по ID
}
