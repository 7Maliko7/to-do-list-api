package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	httpHead "github.com/go-http-utils/headers"
	"io"
	"net/http"
	netUrl "net/url"
	"time"
)

const (
	ContentTypeJson = "application/json"
)

func New(url string) Storage {
	client := http.Client{}
	return Storage{client: client, url: url}
}

type Storage struct {
	client http.Client
	url    string
}

func (s *Storage) CreateTask(name, body string, deadline time.Time) (string, error) {
	data := CreateTaskRequest{
		Name:     name,
		Body:     body,
		Deadline: deadline,
	}
	resp := CreateTaskResponse{}
	err := s.makeRequest("create", http.MethodPost, data, &resp)
	if err != nil {
		return "", nil
	}

	return resp.Uuid, nil
}

func (s *Storage) UpdateTask(uuid string, name, body, status *string, deadline *time.Time) error {
	data := UpdateTaskRequest{
		Uuid:     uuid,
		Name:     name,
		Body:     body,
		Status:   status,
		Deadline: deadline,
	}
	resp := UpdateTaskResponse{}
	err := s.makeRequest("update", http.MethodPatch, data, &resp)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetTask(uuid string) (Task, error) {
	data := GetTaskRequest{
		Uuid: uuid,
	}
	resp := GetTaskResponse{}
	err := s.makeRequest("get", http.MethodGet, data, &resp)
	if err != nil {
		return Task{}, nil
	}
	return Task(resp), nil
}

func (s *Storage) GetListTask() (ListTaskResponse, error) {
	data := ListTaskRequest{}
	resp := ListTaskResponse{}
	err := s.makeRequest("list", http.MethodGet, data, &resp)
	if err != nil {
		return ListTaskResponse{}, nil
	}
	return resp, nil
}

func (s *Storage) DeleteTask(uuid string) error {
	data := DeleteTaskRequest{
		Uuid: uuid,
	}
	resp := DeleteTaskResponse{}
	err := s.makeRequest("delete", http.MethodDelete, data, &resp)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) makeRequest(action, method string, data, response interface{}) error {
	url, _ := netUrl.Parse(fmt.Sprintf("%v/%v", s.url, action))
	headers := http.Header{}
	headers.Add(httpHead.ContentType, ContentTypeJson)
	bodyJson, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req := http.Request{
		Method: method,
		URL:    url,
		Header: headers,
		Body:   io.NopCloser(bytes.NewReader(bodyJson)),
	}
	client := http.Client{}
	resp, err := client.Do(&req)
	toBytes, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(toBytes, response)
	if err != nil {
		return err
	}
	return nil
}
