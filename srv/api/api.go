package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"restapisrv/srv/storage"
	"strconv"
)

// API Программный интерфейс сервера GoNews
type API struct {
	db     storage.Interface
	router *mux.Router
}

// New Конструктор объекта API
func New(db storage.Interface) *API {
	api := API{
		db: db,
	}
	api.router = mux.NewRouter()
	api.endpoints()
	return &api
}

// endpoints Регистрация обработчиков API.
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
	var p storage.Post
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
