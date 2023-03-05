package storage

import (
	"fmt"
	"net/http"
	"net/url"
)

func New(url string)Storage{
	client := http.Client{}
	return Storage{client: client, url: url}
}

type Storage struct{
	client http.Client
	url string
}

func(s *Storage) CreateTask(){
	url,_:= url.Parse(fmt.Sprintf("%v/%v", s.url, "create"))
	req:= http.Request{
		Method:           http.MethodPost,
		URL:              url,
		Proto:            "",
		ProtoMajor:       0,
		ProtoMinor:       0,
		Header:           nil,
		Body:             nil,
		GetBody:          nil,
		ContentLength:    0,
		TransferEncoding: nil,
		Close:            false,
		Host:             "",
		Form:             nil,
		PostForm:         nil,
		MultipartForm:    nil,
		Trailer:          nil,
		RemoteAddr:       "",
		RequestURI:       "",
		TLS:              nil,
		Cancel:           nil,
		Response:         nil,
	}

}