package memdb

import (
	"restapisrv/srv/storage"
	"strings"
)

// Store Хранилище данных.
type Store struct{}

// New Конструктор объекта хранилища.
func New() *Store {
	return new(Store)
}

func (s *Store) AddItem(item storage.LocationItem) error {
	return nil
}

func (s *Store) DeleteItem(item storage.LocationItem) error {
	return nil
}

func (s *Store) DeleteAllItem() error {
	return nil
}

func (s *Store) Items() ([]storage.LocationItem, error) {
	return locations, nil
}

func (s *Store) StringItems() ([]storage.StringLocationItem, error) {
	return stringLocations, nil
}

func (s *Store) SortedItems(item storage.LocationItem) ([]storage.LocationItem, error) {
	var sortedLocations []storage.LocationItem
	subStr := item.Title
	for _, l := range locations {
		if strings.Contains(strings.ToLower(l.Title), subStr) {
			sortedLocations = append(sortedLocations, l)
		}
	}
	return sortedLocations, nil
}

var posts = []storage.Post{
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

var locations = []storage.LocationItem{
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
		Title:     "SyD Apt",
		Content:   "Sydney",
		Link:      "https://yahoo.com",
		Latitude:  151.268865,
		Longitude: -33.885690,
	},
}

var stringLocations = []storage.StringLocationItem{
	{
		ID:        "1",
		Title:     "Msc Apt",
		Content:   "Moscow",
		Link:      "https://ya.ru",
		Latitude:  "55.751244",
		Longitude: "37.618423",
	},
	{
		ID:        "2",
		Title:     "NY Apt",
		Content:   "New York",
		Link:      "https://google.com",
		Latitude:  "40.650002",
		Longitude: "-73.949997",
	},
	{
		ID:        "3",
		Title:     "SyD Apt",
		Content:   "Sydney",
		Link:      "https://yahoo.com",
		Latitude:  "151.268865",
		Longitude: "-33.885690",
	},
}
