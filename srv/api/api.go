package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"restapisrv/srv/storage"
	"time"
)

const logfile = "/var/log/restapi.log"

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
	api.router.HandleFunc("/sortitems", api.SortedItemsHandler).Methods(http.MethodPost, http.MethodOptions)
	api.router.Use(api.Logger)
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

func (api *API) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
		if err != nil {
			http.Error(w, fmt.Sprintf("os.OpenFile error: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		defer file.Close()
		var p storage.LocationItem
		if r.Method != "GET" && r.RequestURI != "/clear" {
			buf, _ := io.ReadAll(r.Body)
			rdr1 := io.NopCloser(bytes.NewBuffer(buf))
			rdr2 := io.NopCloser(bytes.NewBuffer(buf))
			err = json.NewDecoder(rdr1).Decode(&p)
			r.Body = rdr2 // OK since rdr2 implements the io.ReadCloser interface
		}

		rec := httptest.NewRecorder()
		next.ServeHTTP(rec, r)
		for k, v := range rec.Result().Header {
			w.Header()[k] = v
		}
		w.WriteHeader(rec.Code)
		if err != nil {
			return
		}
		rec.Body.WriteTo(w)

		fmt.Fprintf(file, "Time: %s\n", time.Now().Format(time.RFC1123))
		fmt.Fprintf(file, "Remote IP: %s\n", r.RemoteAddr)
		fmt.Fprintf(file, "Method: %s\n", r.Method)
		fmt.Fprintf(file, "Proto: %s\n", r.Proto)
		fmt.Fprintf(file, "URL: %s\n", r.RequestURI)
		fmt.Fprintf(file, "Options: %+v\n", p)
		fmt.Fprintf(file, "HTTP Status: %d\n", rec.Result().StatusCode)
		fmt.Fprintln(file)
	})
}
