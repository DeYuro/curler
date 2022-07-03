package main

import (
	"context"
	"errors"
	"flag"
	"github.com/deyuro/curler/internal/config"
	"github.com/deyuro/curler/internal/service"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func app() error {
	appCtx, cancel := context.WithCancel(context.Background())

	errChan := make(chan error)
	go func() {
		handleSIGINT()
		cancel()
	}()

	go func() {
		errChan <- run(logrus.New())
	}()
	select {
	case err := <-errChan:
		return err
	case <-appCtx.Done():
		return nil
	}
}

func run(logger *logrus.Logger) error {
	wait := flag.Int(`wait`, 0, `sleep between command`)
	times := flag.Int(`times`, 0, `how many times execute dig list`)
	threads := flag.Int(`threads`, 1, `threads`)
	output := flag.Bool(`output`, false, `show output`)

	requestListPath := flag.String(`requestList`, `.data/requestList.json`, `path to requestList`)
	flag.Parse()

	if *threads < 1 {
		return errors.New("threads must be > 0")
	}

	requestList, err := config.Load(*requestListPath)
	svc := service.NewService(requestList, logger, *times, *output, time.Second*time.Duration(*wait), *threads)
	if err != nil {
		return err
	}

	svc.Run()
	return nil
}

func handleSIGUSR2(logger *logrus.Logger) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGUSR2)
	for range ch {
		level := logger.GetLevel()
		switch level {
		case logrus.DebugLevel:
			logger.Warn("switching log level to INFO")
			logger.SetLevel(logrus.InfoLevel)
		default:
			logger.Warn("switching log level to DEBUG")
			logger.SetLevel(logrus.DebugLevel)
		}
	}
}

func handleSIGINT() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	for range sigCh {
		signal.Stop(sigCh)
		return
	}
}
