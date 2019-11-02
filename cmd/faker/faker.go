package main

import (
	"log"

	"github.com/uknth/faker"
)

func main() {
	faker, err := faker.NewFaker()
	if err != nil {
		log.Fatal(err)
	}

	err = faker.Open()
	if err != nil {
		log.Fatal(err)
	}
}
