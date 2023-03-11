package api

import (
	"encoding/json"
	"fmt"
	"github.com/7Maliko7/to-do-list-api/internal/storage"
	"github.com/go-http-utils/headers"
	"log"
	"net/http"
)

const (
	ContentTypeJson = "application/json"
)

var Store *storage.Storage

type ErrResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func GetListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		list, err := Store.GetListTask()
		if err != nil {
			e := ErrResponse{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
			makeErrorResponse(w, e)
			log.Print(err)
			return
		}

		err = makeResponse(w, list)
		if err != nil {
			e := ErrResponse{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
			makeErrorResponse(w, e)
			log.Print(err)
			return
		}
		return
	}
	e := ErrResponse{
		Code:    http.StatusMethodNotAllowed,
		Message: fmt.Sprintf("expect method GET at /list, got %v", r.Method),
	}
	makeErrorResponse(w, e)
	return
}

func makeResponse(w http.ResponseWriter, data any) error {
	w.Header().Set(headers.ContentType, ContentTypeJson)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return err
	}
	return nil
}

func makeErrorResponse(w http.ResponseWriter, e ErrResponse) {
	w.Header().Set(headers.ContentType, ContentTypeJson)
	w.WriteHeader(e.Code)
	_ = json.NewEncoder(w).Encode(e)
}
