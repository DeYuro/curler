package config

import (
	"github.com/jinzhu/configor"
	"github.com/pkg/errors"
)

type RequestList struct {
	List []Request
}

type Request struct {
	Verb string `json:"verb"`
	Url  string `json:"url"`
	Urn  string `json:"Urn"`
}

func Load(filename string) (RequestList, error) {
	var list RequestList
	if err := configor.Load(&list, filename); err != nil {
		return RequestList{}, errors.WithMessage(err, "failed to load digList")
	}
	return list, nil
}
