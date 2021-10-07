package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/gothicfann/go-config-api/app/handlers"
)

func main() {
	l := log.New(os.Stdout, "configs-api ", log.LstdFlags)
	port := os.Getenv("SERVE_PORT")
	if port == "" {
		l.Fatalln("Environment variable SERVE_PORT not configured properly")
	}

	r := mux.NewRouter()

	ch := handlers.NewConfigs(l)

	getRouter := r.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/configs", ch.GetConfigs)
	getRouter.HandleFunc("/configs/{name}", ch.GetConfig)
	getRouter.HandleFunc("/search", ch.QueryConfigs)
	getRouter.HandleFunc("/health", ch.Health)

	addRouter := r.Methods(http.MethodPost).Subrouter()
	addRouter.HandleFunc("/configs", ch.AddConfig)
	addRouter.Use(ch.MiddlewareValidateConfig)

	putRouter := r.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/configs/{name}", ch.PutConfig)
	putRouter.Use(ch.MiddlewareValidateConfig)

	patchRouter := r.Methods(http.MethodPatch).Subrouter()
	patchRouter.HandleFunc("/configs/{name}", ch.PatchConfig)
	patchRouter.Use(ch.MiddlewareValidateConfig)

	deleteRouter := r.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/configs/{name}", ch.DeleteConfig)

	s := http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      r,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		l.Fatalln(s.ListenAndServe())
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	l.Println("Got signal:", <-c)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	l.Fatalln(s.Shutdown(ctx))
}
