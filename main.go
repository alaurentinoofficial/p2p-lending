package main

import (
	"fmt"
	"p2p-lending/models"
	"time"
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
		Amount: float32(12000),
		HasIndexer: false,
		Indexer: "NONE",
		Yield: float32(2.5),
		Validate: time.Now().AddDate(0,1, 0).String(),
		Status: false,
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

	fmt.Println("\nUser 1 Balance: ", models.GetUserById(user1.ID).Balance)
	fmt.Println("User 2 Balance: ", models.GetUserById(user2.ID).Balance)
	fmt.Println("User 3 Balance: ", models.GetUserById(user3.ID).Balance)
	// ----------------------------------------

}