package main

import (
  "fmt"
  "github.com/juju/loggo"
)

func ConfigureLogger(debug bool) {

  baseLogLevel := "WARNING"
  if debug {
    baseLogLevel = "DEBUG"
  }
  loggo.ConfigureLoggers(fmt.Sprintf("<root>=%v", baseLogLevel))
}
