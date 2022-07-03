package config

import (
	"github.com/jinzhu/configor"
	"github.com/pkg/errors"
	"os"
)

type RequestList struct {
	List []Request
}

type Request struct {
	Verb string `json:"verb"`
	Uri  string `json:"uri"`
}

func Load(filename string) (RequestList, error) {
	var list RequestList
	path, _ := os.Getwd()
	println(path)
	if err := configor.Load(&list, filename); err != nil {
		return RequestList{}, errors.WithMessage(err, "failed to load requestList")
	}
	return list, nil
}
