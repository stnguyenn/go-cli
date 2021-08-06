package main

import (
	"github.com/sirupsen/logrus"
)

func main() {
	greeting := "world"
	logrus.Debugf("hello %v!", greeting)
}
