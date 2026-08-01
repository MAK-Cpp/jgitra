package main

import (
	"log"

	"jgitra/internal/config"
)

func main() {
	_, err := config.Load()
	if err != nil {
		log.Fatal("error loading config: ", err)
	}
}
