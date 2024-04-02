package main

import (
	"context"
	"flag"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
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
	var configPath string
	flag.StringVar(&configPath, "config", "config/config.yaml", "path to config yaml file")
	flag.Parse()

	configData, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	cfg := config.Config{}
	err = yaml.Unmarshal(configData, &cfg)
	if err != nil {
		panic(err)
	}

	mongoRepo, err := repo.NewMongoRepo(cfg.MongoURI, cfg.DbName, cfg.CollectionName)
	if err != nil {
		panic(err)
	}

	puller := services.NewPoller(cfg.PollPeriod, cfg.PollBatchSize, mongoRepo)
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
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
	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
}
