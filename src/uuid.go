package main

import (
	"github.com/google/uuid"
)


func makeUUID() string {
	UUID, err := uuid.NewRandom()

	if err != nil {
		panic(err)
	}

	return UUID.String()
}
