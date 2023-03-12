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

func GetUndoneListHandler(w http.ResponseWriter, r *http.Request) {
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
		undoneList := make([]storage.Task, 0, len(list.List))
		for _, v := range list.List {
			if !v.Status {
				undoneList = append(undoneList, v)
			}

		}
		err = makeResponse(w, undoneList)
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

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		task := storage.Task{}
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			e := ErrResponse{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}
			makeErrorResponse(w, e)
			log.Print(err)
			return
		}
		uuid, err := Store.CreateTask(task.Name, task.Body, task.Deadline)
		if err != nil {
			e := ErrResponse{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}
			makeErrorResponse(w, e)
			log.Print(err)
			return
		}
		res := storage.CreateTaskResponse{Uuid: uuid}

		err = makeResponse(w, res)
		if err != nil {
			log.Print(err)
			return
		}
		return
	}
	e := ErrResponse{
		Code:    http.StatusMethodNotAllowed,
		Message: fmt.Sprintf("expect method POST at /create, got %v", r.Method),
	}
	makeErrorResponse(w, e)
	return

}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		task := storage.GetTaskRequest{}
		_ = json.NewDecoder(r.Body).Decode(&task)
		OneTask, err := Store.GetTask(task.Uuid)
		if err != nil {
			e := ErrResponse{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}

			makeErrorResponse(w, e)
			log.Print(err)
			return
		}
		if OneTask.Uuid == "" {
			e := ErrResponse{
				Code:    http.StatusNotFound,
				Message: "Task not found",
			}
			makeErrorResponse(w, e)
			return
		}

		err = makeResponse(w, task)
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
		Message: fmt.Sprintf("expect method GET at /get, got %v", r.Method),
	}
	makeErrorResponse(w, e)
	return
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		task := storage.DeleteTaskRequest{}
		_ = json.NewDecoder(r.Body).Decode(&task)
		if !isExist(task.Uuid) {
			e := ErrResponse{
				Code:    http.StatusNotFound,
				Message: "Task not found",
			}
			makeErrorResponse(w, e)
			return
		}
		taskList, err := Store.GetListTask()
		if err != nil {
			e := ErrResponse{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
			makeErrorResponse(w, e)
			log.Print(err)
			return
		}
		if taskList.List == nil {
			e := ErrResponse{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
			makeErrorResponse(w, e)
			log.Print(err)
			return
		}
		err = Store.DeleteTask(task.Uuid)
		if err != nil {
			e := ErrResponse{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
			makeErrorResponse(w, e)
			log.Print(err)
			return
		}
		err = makeResponse(w, storage.DeleteTaskResponse{})
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
		Message: fmt.Sprintf("expect method DELETE at /delete, got %v", r.Method),
	}
	makeErrorResponse(w, e)
	return
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPatch {
		task := storage.UpdateTaskRequest{}
		_ = json.NewDecoder(r.Body).Decode(&task)
		if !isExist(task.Uuid) {
			e := ErrResponse{
				Code:    http.StatusNotFound,
				Message: "Task not found",
			}
			makeErrorResponse(w, e)
			return
		}
		err := Store.UpdateTask(task.Uuid, task.Name, task.Body, task.Status, task.Deadline)
		if err != nil {
			e := ErrResponse{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
			makeErrorResponse(w, e)
			log.Print(err)
			return
		}
		err = makeResponse(w, storage.UpdateTaskResponse{Code: 200})
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
		Message: fmt.Sprintf("expect method PATCH at /patch, got %v", r.Method),
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
func isExist(uuid string) bool {
	OneTask, err := Store.GetTask(uuid)
	if err != nil {
		log.Print(err)
		return false
	}
	if OneTask.Uuid == "" {
		return false
	}
	return true
}
