package test

import (
	"fmt"
	"testing"
)

type User struct {
	Name string
}

func TestCopy(t *testing.T) {
	users := []User{{Name: "an"}}

	for i, _ := range users {
		users[i].Name = "bo"
	}

	fmt.Println(users)

	testCopy()
}
