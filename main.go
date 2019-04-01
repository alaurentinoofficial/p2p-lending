package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"p2p-lending/controllers"
	"p2p-lending/middlewares"
	"p2p-lending/models"
	"p2p-lending/utils"
)

var port = ":8080"

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/register", controllers.AddUser).Methods("POST")
	r.HandleFunc("/login", controllers.Login).Methods("POST")
	r.HandleFunc("/user", controllers.GetUser).Methods("GET")
	r.HandleFunc("/user/{id}", controllers.GetUserById).Methods("GET")
	r.HandleFunc("/pay/{paymentID}", controllers.PayLending).Methods("POST")

	r.HandleFunc("/lendings", controllers.GetLendings).Methods("GET")
	r.HandleFunc("/lendings", controllers.AddLending).Methods("POST")
	r.HandleFunc("/lendings/{id}", controllers.GetLendingById).Methods("GET")

	r.HandleFunc("/lenders", controllers.GetLenders).Methods("GET")
	r.HandleFunc("/lenders", controllers.AddLender).Methods("POST")

	r.Use(middlewares.JwtAuthentication)

	if os.Getenv("PORT") != "" {
		port = ":" + os.Getenv("PORT")
	}

	simulate()

	fmt.Println("[*] Listening in ", port)
	_ = http.ListenAndServe(port, r)
}

func simulate() {
	// ---------------[ User ]-----------------
	user1 := models.User{
		Name: "Anderson Laurentino",
		Email: "alaurentino.br@gmail.com",
		Password: "1234567890n",
		Type: 0,
		CpfCnpj: "12345678901",
		Salary: 5000,
		Score: 967,
		Balance: float32(10000),
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
		Salary: 0,
		Score: 967,
		Balance: float32(3000),
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
		Salary: 10000,
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
		Amount: 10000,
		//HasIndex: false,
		//Index: 0,
		//IndexYield: 1,
		PrefixedYield: 10,
		PaymentTimeMonth: 12,
		MonthlyInterestRate: 1.7,
	}
	lending.Create()
	// ----------------------------------------


	// -------------[ Lenders ]----------------
	lender1 := models.Lender{
		User: user2.ID,
		Amount: float32(3000),
		Lending: lending.ID,
	}
	lender1.Create()

	lender2 := models.Lender{
		User: user3.ID,
		Amount: float32(7000),
		Lending: lending.ID,
	}
	lender2.Create()
	// ----------------------------------------


	// ----------------------------------------
	fmt.Println("\nUser 1 Balance: ", user1.Balance, "\t-> ", models.GetUserById(user1.ID).Balance)
	fmt.Println("User 2 Balance: ", user2.Balance, "\t-> ",  models.GetUserById(user2.ID).Balance)
	fmt.Println("User 3 Balance: ", user3.Balance, "\t-> ",  models.GetUserById(user3.ID).Balance)
	// ----------------------------------------


	// --------------[ Payment ]---------------
	payments := models.GetLendingPaymentsByTaker(user1.ID)
	taker := models.GetUserById(user1.ID)

	fmt.Println()
	for i, payment := range payments {
		if i != 12 {
			fmt.Println("[*] Payment month: ", i+1, " -> ", utils.ResponseMap[taker.Pay(payment.ID)])
		}
		//fmt.Println("[*] Payment month: ", i)
		//time.Sleep(time.Millisecond * 300)
	}
	// ----------------------------------------


	// ----------------------------------------
	fmt.Println("\nUser 1 Balance: ", user1.Balance, "\t-> ", models.GetUserById(user1.ID).Balance)
	fmt.Println("User 2 Balance: ", user2.Balance, "\t-> ",  models.GetUserById(user2.ID).Balance)
	fmt.Println("User 3 Balance: ", user3.Balance, "\t-> ",  models.GetUserById(user3.ID).Balance)
	// ----------------------------------------

	for _, statement := range models.GetStatementsByUser(user2.ID) {
		fmt.Println(statement)
	}
}