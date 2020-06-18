package main

import (
	log "github.com/sirupsen/logrus"
)

func CheckErr(err error) {
	if err != nil {
		log.Error(err)
	}
}
