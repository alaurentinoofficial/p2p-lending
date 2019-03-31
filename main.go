package main

import (
	"fmt"
	"p2p-lending/models"
)

func main() {
	// ---------------[ User ]-----------------
	user1 := models.User{
		Name: "Anderson Laurentino",
		Email: "alaurentino.br@gmail.com",
		Password: "1234567890n",
		Type: 0,
		CpfCnpj: "12345678901",
		Score: 967,
		Balance: float32(0),
		State: "PE",
		City: "Recife",
		Neighborhood: "Boa vista",
		ZipCode: "50100000",
		Number: 1000,
		Complement: "",
	}
	user1.Create()

	user2 := models.User{
		Name: "Fulano de tal",
		Email: "fulano.tal@gmail.com",
		Password: "1234567890n",
		Type: 0,
		CpfCnpj: "12345678901",
		Score: 967,
		Balance: float32(10000),
		State: "PE",
		City: "Recife",
		Neighborhood: "Boa vista",
		ZipCode: "50100000",
		Number: 1001,
		Complement: "",
	}
	user2.Create()

	user3 := models.User{
		Name: "Rick e Morty",
		Email: "rick.morty@gmail.com",
		Password: "1234567890n",
		Type: 0,
		CpfCnpj: "12345678901",
		Score: 967,
		Balance: float32(7000),
		State: "PE",
		City: "Recife",
		Neighborhood: "Boa vista",
		ZipCode: "50100000",
		Number: 1001,
		Complement: "",
	}
	user3.Create()
	// ----------------------------------------


	// -------------[The lending]--------------
	lending := models.Lending{
		Taker: user1.ID,
		AlreadyInvested: 0,
		Amount: 12000,
		HasIndex: false,
		Index: -1,
		PrefixedYield: 2.5,
		TimeMonth: 12,
	}
	lending.Create()
	// ----------------------------------------

	// -------------[ Lenders ]----------------
	lender1 := models.Lender{
		User: user2.ID,
		Amount: float32(5000),
		Status: false,
		Lending: lending.ID,
	}
	lender1.Create()

	lender2 := models.Lender{
		User: user3.ID,
		Amount: float32(7000),
		Status: false,
		Lending: lending.ID,
	}
	lender2.Create()
	// ----------------------------------------

	// ----------------------------------------
	fmt.Println("Lending Status: ", lending.Status)

	fmt.Println("\nUser 1 Balance: ", user1.Balance, "\t\t-> ", models.GetUserById(user1.ID).Balance)
	fmt.Println("User 2 Balance: ", user2.Balance, "\t-> ",  models.GetUserById(user2.ID).Balance)
	fmt.Println("User 3 Balance: ", user3.Balance, "\t-> ",  models.GetUserById(user3.ID).Balance)
	// ----------------------------------------

}