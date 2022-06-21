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
	api.router.HandleFunc("/items", api.AddItemHandler).Methods(http.MethodPost, http.MethodOptions)
	api.router.HandleFunc("/items", api.DeleteItemHandler).Methods(http.MethodDelete, http.MethodOptions)
	api.router.HandleFunc("/clear", api.DeleteAllItemHandler).Methods(http.MethodDelete, http.MethodOptions)
	api.router.HandleFunc("/sortitems", api.SortedItemsHandler).Methods(http.MethodGet, http.MethodOptions)
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

// AddItemHandler Добавление публикации.
func (api *API) AddItemHandler(w http.ResponseWriter, r *http.Request) {
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
	w.Write([]byte("Item Added\n"))
}

// DeleteItemHandler Удаление публикации по ID.
func (api *API) DeleteItemHandler(w http.ResponseWriter, r *http.Request) {
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
	w.Write([]byte("Item Deleted\n"))
}

// DeleteAllItemHandler Удаление всех публикаций.
func (api *API) DeleteAllItemHandler(w http.ResponseWriter, r *http.Request) {
	err := api.db.DeleteAllItem()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All Items Deleted\n"))
}

// SortedItemsHandler Получение всех публикаций по substring.
func (api *API) SortedItemsHandler(w http.ResponseWriter, r *http.Request) {
	var p storage.LocationItem
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	items, err := api.db.SortedItems(p)
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
