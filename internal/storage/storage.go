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
	url, _ := netUrl.Parse(fmt.Sprintf("%v/%v", s.url, "create"))
	headers := http.Header{}
	headers.Add(httpHead.ContentType, ContentTypeJson)
	data := CreateTaskRequest{
		Name:     name,
		Body:     body,
		Deadline: deadline,
	}
	bodyJson, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req := http.Request{
		Method: http.MethodPost,
		URL:    url,
		Header: headers,
		Body:   io.NopCloser(bytes.NewReader(bodyJson)),
	}

	client := http.Client{}
	resp, err := client.Do(&req)
	respApi := CreateTaskResponse{}
	toBytes, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(toBytes, &respApi)
	if err != nil {
		return "", err
	}

	return respApi.Uuid, nil
}

func (s *Storage) UpdateTask(uuid string, name, body, status *string, deadline *time.Time) error{
	url, _ := netUrl.Parse(fmt.Sprintf("%v/%v", s.url, "update"))
	headers := http.Header{}
	headers.Add(httpHead.ContentType, ContentTypeJson)
	data := UpdateTaskRequest{
		Uuid:     uuid,
		Name:     name,
		Body:     body,
		Status:   status,
		Deadline: deadline,
	}
	bodyJson, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req := http.Request{
		Method: http.MethodPatch,
		URL:    url,
		Header: headers,
		Body:   io.NopCloser(bytes.NewReader(bodyJson)),
	}

	client := http.Client{}
	resp, err := client.Do(&req)
	respApi := UpdateTaskResponse{}
	toBytes, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(toBytes, &respApi)
	if err != nil {
		return err
	}

	return nil
}
