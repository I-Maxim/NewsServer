package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"newsServer/internal/api"
	"newsServer/internal/config"
	"newsServer/internal/repo"
	"newsServer/internal/services"
	"os"
	"os/signal"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config file failed: %w", err)
	}

	mongoRepo, err := repo.NewMongoRepo(cfg.MongoURI, cfg.DbName, cfg.CollectionName)
	if err != nil {
		return fmt.Errorf("create repository failed: %w", err)
	}

	puller := services.NewPoller(cfg.PollPeriod, cfg.PollBatchSize, mongoRepo)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	go puller.Start(ctx)

	articlesApi := api.NewApi(mongoRepo)

	r := mux.NewRouter()
	r.HandleFunc("/articles", articlesApi.ArticlesHandler).Methods("GET")
	r.HandleFunc("/articles/{id}", articlesApi.ArticleHandler).Methods("GET")

	srv := &http.Server{
		Addr:         cfg.ServerAddr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}
	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		err = srv.Shutdown(ctx)
		if err != nil {
			log.Println(err)
		}
	}()
	return srv.ListenAndServe()
}
