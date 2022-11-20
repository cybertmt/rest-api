package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"restapisrv/api"
	c "restapisrv/constants" // uncomment and rename to constant.go file
	"restapisrv/storage"
	"restapisrv/storage/postgres"
	"time"

	"github.com/foomo/simplecert"
	"github.com/foomo/tlsconfig"
)

func main() {

	// сервер описание
	type server struct {
		db  storage.RestInterface
		api *api.API
	}
	// создаём объект сервера
	var srv server

	// реляционная БД PostgresSQL
	db1, errCon := postgres.New("postgres://" + c.User + ":" + c.Pwd + "@" + c.Host + ":" + c.Port + "/" + c.BdName)
	if errCon != nil {
		log.Println(errCon)
	}

	// инициализируем хранилище сервера конкретной БД
	srv.db = db1

	// создаём объект API и регистрируем обработчики
	srv.api = api.New(srv.db)

	// Запускаем сервис на порту 80.
	// Предаём серверу маршрутизатор запросов.
	go func() {
		log.Fatal(http.ListenAndServe(":80", srv.api.Router()))
	}()
	log.Println("HTTP server started @ :80")
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

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
	cfg.Domains = []string{c.Url}
	cfg.CacheDir = c.CertDir
	cfg.SSLEmail = c.Email

	// Запускаем веб-сервер на порту 80 на всех интерфейсах.
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
	// log.Println("starting HTTP Listener on Port 80")
	// go http.ListenAndServe(":80", http.HandlerFunc(simplecert.Redirect))

	// enable hot reload
	tlsConf.GetCertificate = certReloader.GetCertificateFunc()

	// start serving
	log.Println("will serve at: https://" + cfg.Domains[0])
	serve(ctx, srv2)
	<-signalCh
	log.Println("News HTTP server stopped")

	<-make(chan bool)
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
