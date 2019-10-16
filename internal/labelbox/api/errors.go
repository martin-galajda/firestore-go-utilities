package api

import "fmt"

type NotFound error

type notFoundError struct {
	msg string
}

func (err notFoundError) Error() string {
	return fmt.Sprintf("Error - Not Found: %v", err.msg)
}
