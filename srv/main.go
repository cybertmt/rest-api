package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// Post - публикация.
type Post struct {
	ID      int    // номер записи
	Title   string // заголовок публикации
	Content string // содержание публикации
	PubTime int64  // время публикации
	Link    string // ссылка на источник
}

// LocationItem - публикация.
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
	AddPost(Post) error                       // создание новой публикации
	DeletePost(Post) error                    // удаление публикации по ID
	PostsNItems(Post) ([]LocationItem, error) // получение публикации по ID
}

var posts = []Post{
	{
		ID:      1,
		Title:   "Статья 1",
		Content: "Содержание статьи 1",
		PubTime: 1,
		Link:    "http://http1",
	},
	{
		ID:      2,
		Title:   "Статья 2",
		Content: "Содержание статьи 2",
		PubTime: 2,
		Link:    "http://http2",
	},
	{
		ID:      3,
		Title:   "Статья 3",
		Content: "Содержание статьи 3",
		PubTime: 3,
		Link:    "http://http3",
	},
}

var locations = []LocationItem{
	{
		ID:        1,
		Title:     "Msc Apt",
		Content:   "Moscow",
		Link:      "https://ya.ru",
		Latitude:  55.751244,
		Longitude: 37.618423,
	},
	{
		ID:        2,
		Title:     "NY Apt",
		Content:   "New York",
		Link:      "https://google.com",
		Latitude:  40.650002,
		Longitude: -73.949997,
	},
	{
		ID:        3,
		Title:     "Syd Apt",
		Content:   "Sydney",
		Link:      "https://yahoo.com",
		Latitude:  151.268865,
		Longitude: -33.885690,
	},
}

// Store Хранилище данных.
type Store struct{}

// NewStore Конструктор объекта хранилища.
func NewStore() *Store {
	return new(Store)
}

func (s *Store) AddPost(Post) error {
	return nil
}

func (s *Store) DeletePost(Post) error {
	return nil
}

func (s *Store) PostsNItems(Post) ([]LocationItem, error) {
	return locations, nil
}

// API Программный интерфейс сервера GoNews
type API struct {
	db     Interface
	router *mux.Router
}

// NewAPI Конструктор объекта API
func NewAPI(db Interface) *API {
	api := API{
		db: db,
	}
	api.router = mux.NewRouter()
	api.endpoints()
	return &api
}

//  endpoints Регистрация обработчиков API.
func (api *API) endpoints() {
	api.router.HandleFunc("/items", api.postsHandlerNItems).Methods(http.MethodGet, http.MethodOptions)
	api.router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))
}

// Router Получение маршрутизатора запросов.
// Router Требуется для передачи маршрутизатора веб-серверу.
func (api *API) Router() *mux.Router {
	return api.router
}

// postsHandlerNItems Получение публикаций, отсортированных по времени создания, в количестве = n.
func (api *API) postsHandlerNItems(w http.ResponseWriter, r *http.Request) {
	//Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var p Post
	n, _ := strconv.Atoi(mux.Vars(r)["n"])
	p.ID = n
	post, err := api.db.PostsNItems(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bytes, err := json.Marshal(post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(bytes)
}

func main() {
	// канал для чтения новостей rss-каналов и записи в бд
	//PostChannel := make(chan Post)
	// канал сбора ошибок rss обработчика
	//ErrorChannel := make(chan error)

	// структура config.jason для rss обработчика
	//type rssConfig struct {
	//	Rss     []string `json:"rss"`
	//	RPeriod int64    `json:"request_period"`
	//}

	// Читаем файл конфигурации в последовательность байт
	//jsonFile, err := os.Open("cmd/config.json")
	//if err != nil {
	//	ErrorChannel <- err
	//}
	//defer jsonFile.Close()
	//byteValue, _ := ioutil.ReadAll(jsonFile)

	// Конвертируем байты в структуру rss конфигурации
	//var rssCfg rssConfig
	//_ = json.Unmarshal(byteValue, &rssCfg)

	// Сервер GotNews.
	type server struct {
		db  Interface
		api *API
	}
	// Создаём объект сервера.
	var srv server

	// используем переменные окружения для адреса бд, имени бд и пароля
	//host := os.Getenv("host")
	//bdName := os.Getenv("bdName")
	//pwd := os.Getenv("pwd")

	//Создаём объекты баз данных.
	//
	//БД в памяти.
	db1 := NewStore()

	// Реляционная БД PostgreSQL.
	//db2, err := postgres.New("postgres://cyber:" + pwd + "@" + host + "/" + bdName)
	//if err != nil {
	//	ErrorChannel <- err
	//}
	//Документная БД MongoDB.
	//db3, err := mongo.New("mongodb://" + host + ":27017/")
	//if err != nil {
	//	ErrorChannel <- err
	//}
	//_, _, _ = db1, db2, db3

	//Инициализируем хранилище сервера конкретной БД.
	srv.db = db1

	// Запускаем запись новостей из канал PostChannel в бд в отдельной горутине
	//go func() {
	//	for {
	//		err := srv.db.AddPost(<-PostChannel)
	//		if err != nil {
	//			ErrorChannel <- err
	//		}
	//	}
	//}()

	// Запускаем логирование ошибок из канала ErrorChannel
	//go func() {
	//	for {
	//		log.Println(<-ErrorChannel)
	//	}
	//}()

	// Создаём объект API и регистрируем обработчики.
	srv.api = NewAPI(srv.db)

	// Запускаем веб-сервер на порту 8080 на всех интерфейсах.
	// Предаём серверу маршрутизатор запросов,
	// поэтому сервер будет все запросы отправлять на маршрутизатор.
	// Маршрутизатор будет выбирать нужный обработчик.
	_ = http.ListenAndServe(":7531", srv.api.Router())
	//	if err != nil {
	//		ErrorChannel <- err
	//	}
}
