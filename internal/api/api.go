package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"newsServer/internal/repo"
	"time"
)

type Response struct {
	Status   string   `json:"status"`
	Data     any      `json:"data"`
	Metadata Metadata `json:"metadata"`
}

type Metadata struct {
	CreatedAt  time.Time `json:"createdAt"`
	TotalItems int       `json:"totalItems,omitempty"`
	//Sort       string    `json:"sort"`
}

const (
	statusError   = "error"
	statusFail    = "fail"
	statusSuccess = "success"
)

type Api struct {
	repo repo.Repository
}

func NewApi(repository repo.Repository) *Api {
	return &Api{
		repo: repository,
	}
}

func (a *Api) ArticlesHandler(rw http.ResponseWriter, r *http.Request) {
	articles, err := a.repo.List(r.Context())
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(rw).Encode(&Response{
			Status: statusError,
			Data:   "unable to communicate with database",
			Metadata: Metadata{
				CreatedAt: time.Now(),
			},
		})
		if err != nil {
			log.Println(err)
		}
		return
	}

	rw.WriteHeader(http.StatusOK)
	err = json.NewEncoder(rw).Encode(&Response{
		Status: statusSuccess,
		Data:   articles,
		Metadata: Metadata{
			CreatedAt:  time.Now(),
			TotalItems: len(articles),
		},
	})
	if err != nil {
		log.Println(err)
	}
}

func (a *Api) ArticleHandler(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok {
		rw.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(rw).Encode(&Response{
			Status: statusFail,
			Data:   "id title is required",
			Metadata: Metadata{
				CreatedAt: time.Now(),
			},
		})
		if err != nil {
			log.Println(err)
		}
		return
	}

	article, err := a.repo.Load(r.Context(), id)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(rw).Encode(&Response{
			Status: statusError,
			Data:   "unable to communicate with database",
			Metadata: Metadata{
				CreatedAt: time.Now(),
			},
		})
		if err != nil {
			log.Println(err)
		}
		return
	}

	rw.WriteHeader(http.StatusOK)
	err = json.NewEncoder(rw).Encode(&Response{
		Status: statusSuccess,
		Data:   article,
		Metadata: Metadata{
			CreatedAt: time.Now(),
		},
	})
	if err != nil {
		log.Println(err)
	}
}
