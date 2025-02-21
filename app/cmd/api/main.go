package main

import (
	"fmt"

	"github.com/shun-ideguchi/golang-ddd/internal/domain/model/user"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/service"
	"github.com/shun-ideguchi/golang-ddd/internal/infrastructure/persistence"
)

func main() {
	msg := "success"
	userRepository := persistence.NewUserPersistence()
	userService := service.NewUserService(userRepository)

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
