package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
)

// API Программный интерфейс приложения
type API struct {
	r *mux.Router
}

// New Конструктор объекта API
func New() *API {
	api := API{}
	api.r = mux.NewRouter()
	api.endpoints()
	return &api
}

// endpoints Регистрация обработчиков API.
func (api *API) endpoints() {

	// вывод списка новостей
	api.r.HandleFunc("/news/latest", api.latest).Methods(http.MethodGet)
	// вывод отфильтрованных новостей
	api.r.HandleFunc("/news/filter", api.filter).Methods(http.MethodGet)
	// детальное получение новости
	api.r.HandleFunc("/news/detailed", api.detailed).Methods(http.MethodGet)
	// добавление комментария
	api.r.HandleFunc("/comments/store", api.storeComment).Methods(http.MethodPost)
	// сквозная идентификация запросов
	api.r.Use(api.reqId)
}

func main() {
	// Структура сервера API Gateway
	type server struct {
		api *api.API
	}

	var srv server
	srv.api = api.New()

	// Запускаем сервис на порту 8080, интерфейс localhost.
	// Предаём серверу маршрутизатор запросов.
	go func() {
		log.Fatal(http.ListenAndServe("localhost:8080", srv.api.Router()))
	}()
	log.Println("API Gateway HTTP server started @ localhost:8080")
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	<-signalCh
	log.Println("API Gateway HTTP server stopped")
}
