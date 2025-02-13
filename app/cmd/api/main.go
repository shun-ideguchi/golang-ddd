package main

import (
	"fmt"

	"github.com/shun-ideguchi/golang-ddd/internal/domain/user"
)

func main() {
	user1, err := user.NewUser("uuid1234", "山田 太郎")
	if err != nil {
		fmt.Println("failed")
	}

	user2, err := user.NewUser("uuid123", "山田 太郎")
	if err != nil {
		fmt.Println("failed")
	}

	if user1.Equals(user2) {
		fmt.Println("same object")
	} else {
		fmt.Println("not same object")
	}
}
