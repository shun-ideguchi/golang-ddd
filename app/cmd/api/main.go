package main

import (
	"fmt"

	"github.com/shun-ideguchi/golang-ddd/internal/application/usecase/user"
	"github.com/shun-ideguchi/golang-ddd/internal/application/usecase/user/command"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/model/factory"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/service"
	"github.com/shun-ideguchi/golang-ddd/internal/infrastructure/persistence"
)

func main() {
	msg := "success"
	userRepository := persistence.NewUserPersistence()
	userService := service.NewUserService(userRepository)
	userFactory := factory.NewUserFactory()
	userCreateUsecase := user.NewCreateUsecase(userFactory, userRepository, *userService)

	// 新規作成
	createCommand := command.NewCreateCommand("test name", "test@test.com")
	if err := userCreateUsecase.Execute(createCommand); err != nil {
		msg = err.Error()
	}

	fmt.Println(msg)
}
