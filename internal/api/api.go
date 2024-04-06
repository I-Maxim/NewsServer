package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"newsServer/internal/domain"
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

func writeJson(rw http.ResponseWriter, statusCode int, statusMsg string, data any, totalItems int) {
	rw.WriteHeader(statusCode)
	err := json.NewEncoder(rw).Encode(&Response{
		Status: statusMsg,
		Data:   data,
		Metadata: Metadata{
			CreatedAt:  time.Now(),
			TotalItems: totalItems,
		},
	})
	if err != nil {
		log.Println(err)
	}
}

func writeDbError(rw http.ResponseWriter) {
	writeJson(rw, http.StatusInternalServerError, statusError, "unable to communicate with database", 0)
}

func writeRequestError(rw http.ResponseWriter) {
	writeJson(rw, http.StatusBadRequest, statusFail, "id title is required", 0)
}

func writeSuccessResponse(rw http.ResponseWriter, data any, totalItems int) {
	writeJson(rw, http.StatusOK, statusSuccess, data, totalItems)
}

func writeSuccessResponseWithSingleData(rw http.ResponseWriter, data any) {
	writeJson(rw, http.StatusOK, statusSuccess, data, 0)
}

func (a *Api) ArticlesHandler(rw http.ResponseWriter, r *http.Request) {
	articles, err := a.repo.List(r.Context())
	if err != nil {
		log.Println(err)
		writeDbError(rw)
		return
	}

	responseData := make([]*domain.ArticleResponse, len(articles))
	for i, article := range articles {
		responseData[i] = domain.MapDBtoResponseModel(article)
	}

	writeSuccessResponse(rw, responseData, len(articles))
}

func (a *Api) ArticleHandler(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok {
		writeRequestError(rw)
		return
	}

	article, err := a.repo.Load(r.Context(), id)
	if err != nil {
		log.Println(err)
		writeDbError(rw)
		return
	}

	writeSuccessResponseWithSingleData(rw, article)
}
