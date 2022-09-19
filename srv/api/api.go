package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	c "restapisrv/constants" // uncomment and rename constant.go file
	"restapisrv/sendemail"
	"restapisrv/storage"

	web "restapisrv/webpages"
	"strings"
	"time"

	"github.com/gorilla/mux"
	iuliia "github.com/mehanizm/iuliia-go"
	"golang.org/x/crypto/bcrypt"
)

//regexp
var reg = regexp.MustCompile(`[^A-Za-z0-9]`)

// API программный интерфейс сервера.
type API struct {
	db     storage.RestInterface
	router *mux.Router
}

// New конструктор объекта API.
func New(db storage.RestInterface) *API {
	api := API{
		db: db,
	}
	api.router = mux.NewRouter()
	api.endpoints()
	return &api
}

// endpoints регистрация обработчиков API.
func (api *API) endpoints() {

	api.router.HandleFunc("/products", api.ProductsHandler).Methods(http.MethodGet, http.MethodOptions)
	api.router.HandleFunc("/products", api.AddProductHandler).Methods(http.MethodPost, http.MethodOptions)
	api.router.HandleFunc("/products", api.DeleteProductHandler).Methods(http.MethodDelete, http.MethodOptions)
	api.router.HandleFunc("/clearproducts", api.DeleteAllProductsHandler).Methods(http.MethodDelete, http.MethodOptions)
	api.router.HandleFunc("/sortproducts", api.SearchSortedProductsHandler).Methods(http.MethodPost, http.MethodOptions)
	api.router.HandleFunc("/stores", api.StoresHandler).Methods(http.MethodGet, http.MethodOptions)
	api.router.HandleFunc("/stores", api.AddStoreHandler).Methods(http.MethodPost, http.MethodOptions)
	api.router.HandleFunc("/stores", api.DeleteStoreHandler).Methods(http.MethodDelete, http.MethodOptions)
	api.router.HandleFunc("/clearstores", api.DeleteAllStoresHandler).Methods(http.MethodDelete, http.MethodOptions)
	api.router.HandleFunc("/prices", api.AddUpdatePriceHandler).Methods(http.MethodPost, http.MethodOptions)
	api.router.HandleFunc("/prices", api.DeletePriceHandler).Methods(http.MethodDelete, http.MethodOptions)
	api.router.HandleFunc("/clearprices", api.DeleteAllPricesHandler).Methods(http.MethodDelete, http.MethodOptions)
	api.router.HandleFunc("/pricelist", api.PriceListHandler).Methods(http.MethodGet, http.MethodOptions)
	api.router.HandleFunc("/productprice", api.ProductPriceHandler).Methods(http.MethodPost, http.MethodOptions)
	api.router.HandleFunc("/signup", api.SignUpHandler).Methods(http.MethodPost, http.MethodOptions)
	api.router.HandleFunc("/signin", api.SignInHandler).Methods(http.MethodPost, http.MethodOptions)
	api.router.HandleFunc("/emailconfirm", api.EmailConfirmHandler).Methods(http.MethodGet, http.MethodOptions)

	//api.router.HandleFunc("/", api.PageHandler).Methods(http.MethodGet, http.MethodOptions)
	//api.router.Use(api.Logger)
	fs := http.FileServer(http.Dir("./webpages"))
	//cssHandler := http.FileServer(http.Dir("./webpages/css"))
	//imagesHandler := http.FileServer(http.Dir("./webpages/images"))

	api.router.Handle("/", fs)
	//api.router.PathPrefix("/emailconfirm/").Handler(http.StripPrefix("/emailconfirm/", fs))
	api.router.PathPrefix("/").Handler(fs)
	//api.router.PathPrefix("/images/").Handler(http.StripPrefix("/images/", fs))
	//api.router.PathPrefix("/css/").Handler(fs)
	//api.router.PathPrefix("/emailconfirm/css/").Handler(http.StripPrefix("/emailconfirm/", fs))
}

// Router получение маршрутизатора запросов.
// Router требуется для передачи маршрутизатора веб-серверу.
func (api *API) Router() *mux.Router {
	return api.router
}

// Logger логирование запросов в файл.
func (api *API) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := os.OpenFile(c.Logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
		if err != nil {
			http.Error(w, fmt.Sprintf("os.OpenFile error: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		defer file.Close()
		var p storage.ProductItem
		if r.Method != "GET" && r.RequestURI != "/clearproducts" && r.RequestURI != "/clearstores" && r.RequestURI != "/clearprices" {
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

// Продукты.
// ProductsHandler получение всех продуктов.
func (api *API) ProductsHandler(w http.ResponseWriter, r *http.Request) {
	products, err := api.db.Products()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bytes, err := json.Marshal(products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}

// AddProductHandler добавление продукта.
func (api *API) AddProductHandler(w http.ResponseWriter, r *http.Request) {
	var p storage.ProductItem
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(p.Prod_tr_name) == 0 {
		p.Prod_tr_name = iuliia.Gost_7034.Translate(p.Prod_name)
	}
	err = api.db.AddProduct(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product Added\n"))
}

// DeleteProductHandler удаление продукта по ID.
func (api *API) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	var p storage.ProductItem
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = api.db.DeleteProduct(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product Deleted\n"))
}

// DeleteAllProductsHandler удаление всех продуктов.
func (api *API) DeleteAllProductsHandler(w http.ResponseWriter, r *http.Request) {
	err := api.db.DeleteAllProducts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All Products Deleted\n"))
}

// SearchSortedProductsHandler получение всех продуктов по substring.
func (api *API) SearchSortedProductsHandler(w http.ResponseWriter, r *http.Request) {
	var p storage.SearchItem
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	products, err := api.db.SearchSortedProducts(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bytes, err := json.Marshal(products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}

// Магазины.
// StoresHandler получение всех магазинов.
func (api *API) StoresHandler(w http.ResponseWriter, r *http.Request) {
	stores, err := api.db.Stores()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bytes, err := json.Marshal(stores)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}

// AddStoreHandler добавление магазина.
func (api *API) AddStoreHandler(w http.ResponseWriter, r *http.Request) {
	var s storage.StoreItem
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = api.db.AddStore(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Store Added\n"))
}

// DeleteStoreHandler удаление магазина по ID.
func (api *API) DeleteStoreHandler(w http.ResponseWriter, r *http.Request) {
	var s storage.StoreItem
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = api.db.DeleteStore(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Store Deleted\n"))
}

// DeleteAllStoresHandler удаление всех магазинов.
func (api *API) DeleteAllStoresHandler(w http.ResponseWriter, r *http.Request) {
	err := api.db.DeleteAllStores()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All Stores Deleted\n"))
}

// Цены.
// AddPriceHandler добавление цены на продукт в магазине.
func (api *API) AddUpdatePriceHandler(w http.ResponseWriter, r *http.Request) {
	var pr storage.PriceItem
	err := json.NewDecoder(r.Body).Decode(&pr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = api.db.AddUpdatePrice(pr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Price Added\n"))
}

// DeleteAllPricesHandler удаление всех цены.
func (api *API) DeleteAllPricesHandler(w http.ResponseWriter, r *http.Request) {
	err := api.db.DeleteAllPrices()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All Prices Deleted\n"))
}

// DeletePriceHandler удаление магазина по store_id + product_id.
func (api *API) DeletePriceHandler(w http.ResponseWriter, r *http.Request) {
	var pr storage.PriceItem
	err := json.NewDecoder(r.Body).Decode(&pr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = api.db.DeletePrice(pr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Price Deleted\n"))
}

// PriceList.
func (api *API) PriceListHandler(w http.ResponseWriter, r *http.Request) {
	prices, err := api.db.PriceList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bytes, err := json.Marshal(prices)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}

// ProductPriceHandler получение всех цен по названию прордукта prod_name.
func (api *API) ProductPriceHandler(w http.ResponseWriter, r *http.Request) {
	var pr storage.PriceListItem
	err := json.NewDecoder(r.Body).Decode(&pr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	prices, err := api.db.ProductPrice(pr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bytes, err := json.Marshal(prices)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}

// Users.
// SignUpHandler добавление нового пользователя.
func (api *API) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var user storage.Credentials
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password+user.Useremail), 8)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	err = api.db.SignUp(user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			w.WriteHeader(http.StatusConflict)
			return
		}
		// If there is an issue with the database, return a 500 error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	timeNow := time.Now().String()
	confirmLink, err := bcrypt.GenerateFromPassword([]byte(user.Useremail+timeNow), 8)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var userConfirm storage.CredentialsConfirm
	strConfirmlink := string(confirmLink)
	strConfirmlink = reg.ReplaceAllString(strConfirmlink, "x")
	log.Print("ConfirmString is ", strConfirmlink)
	userConfirm.Useremail = user.Useremail
	userConfirm.Confirmstring = strConfirmlink
	err = api.db.SetConfirmString(userConfirm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	f := c.EmailFrom
	pass := c.EmailPassword
	t := []string{user.Useremail}
	h := c.EmailHost
	p := c.EmailPort
	s := "Subject: " + "Активация email адреса\n"
	b := "Здравствуйте!\n\nДля получения уведомлений об активации, " +
		"остановке и необходимости продления услуг, необходимо подтвердить адрес электронной почты. " +
		"Для подтверждения адреса, пожалуйста, перейдите по ссылке: \n" +
		c.FullURL + "/emailconfirm?k=" + string(strConfirmlink)

	err = sendemail.Sendemail(f, pass, t, b, h, p, s)
	if err != nil {
		log.Print(err)
	}

	w.WriteHeader(http.StatusOK)
}

// SignInHandler вход пользователя.
func (api *API) SignInHandler(w http.ResponseWriter, r *http.Request) {
	var user storage.CredentialsShort
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := api.db.SignIn(user)
	var fixedResult storage.CredentialsFixed

	if err != nil {
		if err == storage.ErrUserNotFound {
			fixedResult.Useremail = user.Useremail + " not found"
			bytes, err := json.Marshal(fixedResult)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Print("!!Not found")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(bytes)
			return
		}
		// If there is an issue with the database, return a 500 error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// if result.Userstatus < 1 {
	// 	fixedResult.Useremail = user.Useremail + " not confirmed"
	// 	bytes, err := json.Marshal(fixedResult)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// 	log.Print("!!Not confirmed")
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	w.Write(bytes)
	// 	return
	// }

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password+user.Useremail))
	// If the two passwords don't match, return a 401 status
	if err != nil {
		fixedResult.Useremail = user.Useremail + " is illegal"
		bytes, err := json.Marshal(fixedResult)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Print("!!Illigal")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(bytes)
		return
	}
	fixedResult.Useremail = user.Useremail
	fixedResult.Usernickname = result.Usernickname
	fixedResult.Userstatus = result.Userstatus
	bytes, err := json.Marshal(fixedResult)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Print("!! loggedIned", user)
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

// EmailConfirmHandler подтверждение почты.
func (api *API) EmailConfirmHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	var user storage.CredentialsConfirm
	user.Confirmstring = params.Get("k")

	result, err := api.db.ConfirmStringAndStatus(user)
	data := map[string]string{
		"useremail": result.Useremail,
		"reply":     "already confirmed",
	}
	dir := "./webpages/emailconfirm.html"
	if err != nil {
		if err == storage.ErrAlreadyConfirmed {
			w.WriteHeader(http.StatusOK)
			err = web.SendWebPage(w, dir, data)
			if err != nil {
				log.Print(err)
				return
			}

			return
		}
		if err == storage.ErrConfirmStringNotFound {
			data["reply"] = "got invalid confirm string"
			w.WriteHeader(http.StatusUnauthorized)
			err = web.SendWebPage(w, dir, data)
			if err != nil {
				log.Print(err)
				return
			}
			return
		}
		// If there is an issue with the database, return a 500 error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data["reply"] = "confirmed"

	w.WriteHeader(http.StatusOK)
	err = web.SendWebPage(w, dir, data)
	if err != nil {
		log.Print(err)
		return
	}
	//w.Write([]byte("Email " + result.Useremail + " confirmed"))
}

// // EmailConfirmHandler подтверждение почты.
// func (api *API) PageHandler(w http.ResponseWriter, r *http.Request) {

// 		// If there is an issue with the database, return a 500 error
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// }
