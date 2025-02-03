package main

import (
	"fmt"

	"github.com/shun-ideguchi/golang-ddd/internal/domain/model/fullname"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/model/money"
)

func main() {
	// FullName 失敗Ver.

	// lastNameが10文字以上となり、lastNameのルールに違反しています
	// 不正な値を存在させないことになります
	username, err := fullname.NewFullName("山田1", "太郎345678901")
	if err != nil {
		errorText := fmt.Sprintf("FullName Error: %s", err.Error())
		fmt.Println(errorText)
	} else {
		fmt.Println(username)
	}

	// FullName 成功Ver.

	username2, err := fullname.NewFullName("山田1", "太郎1")
	if err != nil {
		errorText := fmt.Sprintf("FullName Error: %s", err.Error())
		fmt.Println(errorText)
	} else {
		fmt.Println(username2)
	}

	// Money 加算 失敗Ver.

	// Moneyはユーロを想定していないためcurrencyのルールに違反しています
	// 不正な値を存在させないことになります
	myMoney, err := money.NewMoney(1000, "EUR")
	if err != nil {
		errorText := fmt.Sprintf("Money Error: %s", err.Error())
		fmt.Println(errorText)
	} else {
		fmt.Println(myMoney)
	}

	// Money加算 成功Ver.

	myMoney2, err := money.NewMoney(1000, "JPY")
	if err != nil {
		errorText := fmt.Sprintf("Money Error: %s", err.Error())
		fmt.Println(errorText)
	} else {
		fmt.Println(myMoney)
	}
	myMoney3, err := money.NewMoney(1000, "JPY")
	if err != nil {
		errorText := fmt.Sprintf("Money Error: %s", err.Error())
		fmt.Println(errorText)
	} else {
		fmt.Println(myMoney)
	}
	// Moneyインスタンス同士を加算し、新しいオブジェクトとして返却
	// 値オブジェクトを交換可能としている
	result, err := myMoney2.Add(*myMoney3)
	if err != nil {
		errorText := fmt.Sprintf("Money Error: %s", err.Error())
		fmt.Println(errorText)
	} else {
		fmt.Println(result)
	}
}
