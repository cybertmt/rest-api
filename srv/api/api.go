package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"restapisrv/srv/storage"
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
	api.router.HandleFunc("/items", api.ItemsHandler).Methods(http.MethodGet, http.MethodOptions)
	api.router.HandleFunc("/stringitems", api.StringItemsHandler).Methods(http.MethodGet, http.MethodOptions)
	api.router.HandleFunc("/items", api.addItemHandler).Methods(http.MethodPost, http.MethodOptions)
	api.router.HandleFunc("/items", api.deleteItemHandler).Methods(http.MethodDelete, http.MethodOptions)
}

// Router Получение маршрутизатора запросов.
// Router Требуется для передачи маршрутизатора веб-серверу.
func (api *API) Router() *mux.Router {
	return api.router
}

// ItemsHandler Получение всех публикаций.
func (api *API) ItemsHandler(w http.ResponseWriter, r *http.Request) {
	items, err := api.db.Items()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bytes, err := json.Marshal(items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}

// StringItemsHandler Получение всех публикаций.
func (api *API) StringItemsHandler(w http.ResponseWriter, r *http.Request) {
	items, err := api.db.StringItems()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bytes, err := json.Marshal(items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}

// addItemHandler Добавление публикации.
func (api *API) addItemHandler(w http.ResponseWriter, r *http.Request) {
	var p storage.LocationItem
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = api.db.AddItem(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// deleteItemHandler Удаление публикации.
func (api *API) deleteItemHandler(w http.ResponseWriter, r *http.Request) {
	var p storage.LocationItem
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = api.db.DeleteItem(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
