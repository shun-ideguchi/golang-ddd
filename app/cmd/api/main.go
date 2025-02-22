package main

import (
	"fmt"

	"github.com/google/uuid"
	usecase "github.com/shun-ideguchi/golang-ddd/internal/application/usecase/user"
	"github.com/shun-ideguchi/golang-ddd/internal/application/usecase/user/command"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/service"
	"github.com/shun-ideguchi/golang-ddd/internal/infrastructure/persistence"
)

func main() {
	msg := "success"
	userRepository := persistence.NewUserPersistence()
	userService := service.NewUserService(userRepository)
	userCreateUsecase := usecase.NewCreateUsecase(userRepository, *userService)

	// 新規作成
	createCommand := command.NewCreateCommand("test name", "test@test.com")
	if err := userCreateUsecase.Execute(createCommand); err != nil {
		msg = err.Error()
	}

	getUsecase := usecase.NewGetUsecase(userRepository)
	// 取得
	dto, err := getUsecase.Execute(uuid.NewString())
	if err != nil {
		msg = err.Error()
	} else {
		fmt.Println(dto)
	}

	userUpdateUsecase := usecase.NewUpdateUsecase(userRepository, *userService)
	// 更新
	updateCommand := command.NewUpdateCommand(uuid.NewString(), command.WithName("update name")) // WithEmailを実行しないことでEmailの更新を制御する
	if err := userUpdateUsecase.Execute(updateCommand); err != nil {
		msg = err.Error()
	}

	userDeleteUsecase := usecase.NewDeleteUsecase(userRepository)
	// 削除
	deleteCommand := command.NewDeleteCommand(uuid.NewString())
	if err := userDeleteUsecase.Execute(deleteCommand); err != nil {
		msg = err.Error()
	}

	fmt.Println(msg)
}
