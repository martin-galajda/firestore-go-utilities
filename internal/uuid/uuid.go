package uuid

import (
	"github.com/google/uuid"
)

func MakeUUID() string {
	UUID, err := uuid.NewRandom()

	if err != nil {
		panic(err)
	}

	return UUID.String()
}
