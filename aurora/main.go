package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"./models"
	"./routes"

	"github.com/gorilla/handlers"
)

func main() {
	logs := &models.Logger{}
	logs.InitLogging("House - API", os.Stdout, os.Stdout, os.Stdout, os.Stderr,
		os.Stderr, os.Stdout)

	router, con := routes.NewRouter() // create routes

	f, err := os.OpenFile("connections.log", os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0666)
	if err != nil {
		fmt.Println("Failed to initialize logger")
	}
	defer f.Close()
	con.Logger.Log.SetOutput(f)

	// These two lines are important in order to allow access from the
	//	front-end side to the methods
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST",
		"DELETE", "PUT"})

	cfg := &tls.Config{
		MinVersion: tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384,
			tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	srv := &http.Server{
		Addr:      ":" + con.Session.LiteConfig.Config["default"]["port"],
		Handler:   handlers.CORS(allowedOrigins, allowedMethods)(router),
		TLSConfig: cfg,
		TLSNextProto: make(map[string]func(
			*http.Server,
			*tls.Conn,
			http.Handler), 0),
	}

	// Launch server with CORS validations
	logs.Fatal.Println(srv.ListenAndServeTLS("certs/server.rsa.crt",
		"certs/server.rsa.key"))
}
