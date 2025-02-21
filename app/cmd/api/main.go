package main

import (
	"fmt"

	"github.com/shun-ideguchi/golang-ddd/internal/domain/model/user"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/service"
)

func main() {
	msg := "success"
	userService := service.NewUserService()

	user, err := user.NewUser("uuid", "test")
	if err != nil {
		msg = "failed to initialize user model"
	}

	isExist := userService.IsExists(user)
	if isExist {
		msg = "duplicate user"
	}

	fmt.Println(msg)
}
