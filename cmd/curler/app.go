package main

import (
	"github.com/sirupsen/logrus"
)

func main() {
	if err := app(); err != nil {
		logrus.WithError(err).Fatal("application failed with error")
	}
}
