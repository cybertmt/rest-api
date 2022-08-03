package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"restapisrv/srv/api"
	"restapisrv/srv/storage"
	"restapisrv/srv/storage/postgres"
	"time"

	"github.com/foomo/simplecert"
	"github.com/foomo/tlsconfig"
)

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
		db  storage.Interface
		api *api.API
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
	//db1 := memdb.New()

	// Реляционная БД PostgresSQL.
	db2, _ := postgres.New("postgres://" + user + ":" + pwd + "@" + host + ":" + port + "/" + bdName)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//Документная БД MongoDB.
	//db3, err := mongo.New("mongodb://" + host + ":27017/")
	//if err != nil {
	//	ErrorChannel <- err
	//}
	//_, _, _ = db1, db2, db3

	//Инициализируем хранилище сервера конкретной БД.
	srv.db = db2

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
	srv.api = api.New(srv.db)

	var (
		// the structure that handles reloading the certificate
		certReloader *simplecert.CertReloader
		err          error
		numRenews    int
		ctx, cancel  = context.WithCancel(context.Background())

		// init strict tlsConfig (this will enforce the use of modern TLS configurations)
		// you could use a less strict configuration if you have a customer facing web application that has visitors with old browsers
		tlsConf = tlsconfig.NewServerTLSConfig(tlsconfig.TLSModeServerStrict)

		makeServer = func() *http.Server {
			return &http.Server{
				Addr:      ":443",
				Handler:   srv.api.Router(),
				TLSConfig: tlsConf,
			}
		}
		srv2 = makeServer()

		// init simplecert configuration
		cfg = simplecert.Default
	)

	// configure
	cfg.Domains = []string{url}
	cfg.CacheDir = certDir
	cfg.SSLEmail = email

	// Запускаем веб-сервер на порту 8080 на всех интерфейсах.
	// Предаём серверу маршрутизатор запросов,
	// поэтому сервер будет все запросы отправлять на маршрутизатор.
	// Маршрутизатор будет выбирать нужный обработчик.
	// disable HTTP challenges - we will only use the TLS challenge for this example.
	cfg.HTTPAddress = ""

	// this function will be called just before certificate renewal starts and is used to gracefully stop the service
	// (we need to temporarily free port 443 in order to complete the TLS challenge)
	cfg.WillRenewCertificate = func() {
		// stop server
		cancel()
	}

	cfg.DidRenewCertificate = func() {

		numRenews++

		// restart server: both context and server instance need to be recreated!
		ctx, cancel = context.WithCancel(context.Background())
		srv2 = makeServer()

		// force reload the updated cert from disk
		certReloader.ReloadNow()

		// here we go again
		go serve(ctx, srv2)
	}

	certReloader, err = simplecert.Init(cfg, func() {
		os.Exit(0)
	})
	if err != nil {
		log.Fatal("simplecert init failed: ", err)
	}

	// redirect HTTP to HTTPS
	log.Println("starting HTTP Listener on Port 80")
	go http.ListenAndServe(":80", http.HandlerFunc(simplecert.Redirect))

	// enable hot reload
	tlsConf.GetCertificate = certReloader.GetCertificateFunc()

	// start serving
	log.Println("will serve at: https://" + cfg.Domains[0])
	serve(ctx, srv2)

	//fmt.Println("waiting forever")
	<-make(chan bool)
	//
	//go func() {
	//	_ = http.ListenAndServe(":80", srv.api.Router())
	//}()
	//_ = http.ListenAndServe(":7531", srv.api.Router())
	//	if err != nil {
	//		ErrorChannel <- err
	//	}
}

func serve(ctx context.Context, srv *http.Server) {

	// lets go
	go func() {
		if err := srv.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %+s\n", err)
		}
	}()

	log.Printf("server started")
	<-ctx.Done()
	log.Printf("server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	err := srv.Shutdown(ctxShutDown)
	if err == http.ErrServerClosed {
		log.Printf("server exited properly")
	} else if err != nil {
		log.Printf("server encountered an error on exit: %+s\n", err)
	}
}
