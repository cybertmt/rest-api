package storage

import "errors"

var ErrUserNotFound = errors.New("user not found")
var ErrConfirmStringNotFound = errors.New("invalid confirmstring")
var ErrAlreadyConfirmed = errors.New("already confirmed")

// ProductItem - продукт.
type ProductItem struct {
	Prod_id      int    `json:"prod_id"`
	Prod_name    string `json:"prod_name"`
	Prod_tr_name string `json:"prod_tr_name"`
	Prod_desc1   string `json:"prod_desc1"`
	Prod_desc2   string `json:"prod_desc2"`
	Prod_desc3   string `json:"prod_desc3"`
	Prod_logo    string `json:"prod_logo"`
}

// StoreItem - магазин.
type StoreItem struct {
	Store_id        int     `json:"store_id"`
	Store_name      string  `json:"store_name"`
	Store_tr_name   string  `json:"store_tr_name"`
	Store_address   string  `json:"store_address"`
	Store_phone     string  `json:"store_phone"`
	Store_email     string  `json:"store_email"`
	Store_logo      string  `json:"store_logo"`
	Store_latitude  float64 `json:"store_lat"`
	Store_longitude float64 `json:"store_lon"`
}

// PriceItem - цена на товар в магазине.
type PriceItem struct {
	Store_id int     `json:"store_id"`
	Prod_id  int     `json:"prod_id"`
	Price    float64 `json:"price"`
}

// PriceList - цена на товар в магазине.
type PriceListItem struct {
	Prod_name       string  `json:"prod_name"`
	Prod_logo       string  `json:"prod_logo"`
	Store_name      string  `json:"store_name"`
	Store_address   string  `json:"store_address"`
	Store_phone     string  `json:"store_phone"`
	Store_email     string  `json:"store_email"`
	Store_logo      string  `json:"store_logo"`
	Store_latitude  float64 `json:"store_lat"`
	Store_longitude float64 `json:"store_lon"`
	Price           float64 `json:"price"`
}

// SearchList - цена на товар в магазине.
type SearchItem struct {
	Prod_name string  `json:"prod_name"`
	Price     float64 `json:"price"`
}

// Credentials - учетная запись пользователя.
type Credentials struct {
	Useremail       string `json:"useremail"`
	Password        string `json:"password"`
	Userstatus      int    `json:"status"`
	Confirmstring   string `json:"confirmedstring"`
	Usernickname    string `json:"usernickname"`
	Lastlogindate   string `json:"lastlogindate"`
	Lastlogindevice string `json:"lastlogindevice"`
}

// CredentialsFixed - учетная запись пользователя на SignIn.
type CredentialsFixed struct {
	Useremail    string `json:"useremail"`
	Userstatus   int    `json:"status"`
	Usernickname string `json:"usernickname"`
}

type CredentialsShort struct {
	Useremail string `json:"useremail"`
	Password  string `json:"password"`
}

type CredentialsConfirm struct {
	Useremail     string `json:"useremail"`
	Confirmstring string `json:"confirmstring"`
	Userstatus    int    `json:"status"`
}

type CredentialsUserEmailStatus struct {
	Useremail  string `json:"useremail"`
	Userstatus int    `json:"status"`
}

// RestInterface задаёт новый контракт на работу с БД Products Stores.
type RestInterface interface {
	Products() ([]ProductItem, error)                                           // получение всех продуктов
	AddProduct(prod ProductItem) error                                          // создание новой записи продукта
	DeleteProduct(prod ProductItem) error                                       // удаление продукта по ID
	DeleteAllProducts() error                                                   // удаление всех продуктов, очистка таблицы
	SearchSortedProducts(sr SearchItem) ([]SearchItem, error)                   // выдача продукта для поиска
	Stores() ([]StoreItem, error)                                               // получение всех магазинов
	AddStore(store StoreItem) error                                             // создание новой записи магазина
	DeleteStore(store StoreItem) error                                          // удаление магазина по ID
	DeleteAllStores() error                                                     // удаление всех магазинов, очистка таблицы
	AddUpdatePrice(price PriceItem) error                                       // добавление или обновление цены
	DeletePrice(price PriceItem) error                                          // удаление цены по ID магазина и продукта
	DeleteAllPrices() error                                                     // удаление всех цен, очистка таблицы
	PriceList() ([]PriceListItem, error)                                        // получение всех цен
	ProductPrice(price PriceListItem) ([]PriceListItem, error)                  // получение всех цен по названию продукта
	SignUp(user Credentials) error                                              // добавление нового пользователя
	SignIn(user CredentialsShort) (Credentials, error)                          // вход пользователя
	UserExistEmailStatus(user CredentialsUserEmailStatus) error                 // существует ли пользователь и вернуть статус
	SetConfirmString(user CredentialsConfirm) error                             // добавление строки подтверждения почты
	ConfirmStringAndStatus(user CredentialsConfirm) (CredentialsConfirm, error) // получение строки подтверждения почты и статуса
	ShortPriceList() ([]SearchItem, error)                                      // получение всех цен короткая форма
	//SetUserStatus(user CredentialsConfirm) error                              // изменение статуса пользователя по useremail
}
