package service

import (
	"github.com/deyuro/curler/internal/config"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type service struct {
	requestList config.RequestList
	logger      *logrus.Logger
	times       int
	wait        time.Duration
	output      bool
	threads     int
	client      Doer
}

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

func NewService(digList config.RequestList, logger *logrus.Logger, times int, output bool, wait time.Duration, thread int) *service {
	return &service{requestList: digList,
		logger:  logger,
		times:   times,
		wait:    wait,
		output:  output,
		threads: thread,
		client:  &http.Client{},
	}
}

func (s service) Run() {
	var wg sync.WaitGroup

	for i := 0; i < s.threads; i++ {
		wg.Add(1)
		go func(num int) {
			if s.times == 0 {
				s.infinity()
			} else {
				s.repeatedly()
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func (s *service) infinity() {
	for {
		for _, v := range s.requestList.List {
			s.curl(v.Verb, v.Uri, "")
		}
	}
}

func (s *service) repeatedly() {
	for i := 0; i <= s.times; i++ {
		for _, v := range s.requestList.List {
			s.curl(v.Verb, v.Uri, "")
		}
	}
}

// TODO parse post params
func (s *service) curl(verb, uri, body string) {
	time.Sleep(s.wait)
	s.cmd(verb, uri, body)
}

func (s *service) cmd(verb, uri, body string) {

	req, err := http.NewRequest(verb, uri, strings.NewReader(body))
	if err != nil {
		s.logger.WithError(err).Error("Unable to make request")
		return
	}
	resp, err := s.client.Do(req)

	if err != nil {
		s.logger.WithError(err).Error("Unable to do request")
		return
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.WithError(err).Error("Unable to read body")
		return
	}

	if s.output {
		s.logger.Info(string(b))
	}
}
