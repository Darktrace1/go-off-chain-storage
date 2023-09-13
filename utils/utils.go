package utils

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func Log_init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: time.StampMilli,
	})
}

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ReadFile(path string) []byte {
	ReadInfo, err := os.ReadFile(path)
	CheckErr(err)

	return ReadInfo
}
