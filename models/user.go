package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"p2p-lending/types"
	"time"
)

type User struct {
	ID           string  `json:"id" gorm:"primary_key;"`
	Name         string  `json:"name"`
	Email        string  `json:"email"`
	Password     string  `json:"password"`
	CreationDate string  `json:"creation_date"`
	Type         string  `json:"types"`
	Score        int     `json:"score"`
	CpfCnpj      string  `json:"cpf_cpnj"`
	Balance      float32 `json:"balance"`
	Country      string  `json:"country"`
	State        string  `json:"state"`
	City         string  `json:"city"`
	Neighborhood string  `json:"neighborhood"`
	Number       int     `json:"number"`
	ZipCode      string  `json:"zipcode"`
	Complement   string  `json:"complement"`
}

func (user *User) BeforeCreate(scope *gorm.Scope) error {
	uu, _ := uuid.NewV4()
	_ = scope.SetColumn("ID", uu.String())
	_ = scope.SetColumn("CreationDate", time.Now().UTC().String())
	return nil
}

func (user *User) Create() {
	GetDB().Create(&user)
}

func (user *User) Save() {
	GetDB().Save(&user)
}

func GetUserById(id string) *User {
	user := User{}
	GetDB().Table("users").Where("id = ?", id).First(&user)

	return &user
}

func UserCheckBalance(userID string, amount float32) bool {
	user := GetUserById(userID)
	return user.Balance - amount >= 0
}

func UserLend(userID string, amount float32, lending *Lending) {
	// Find user in database
	user := GetUserById(userID)

	// Remove money
	fmt.Println("User: ", user.ID)
	fmt.Println("Before: ", user.Balance)

	user.Balance -= amount
	user.Save()

	fmt.Println("After: ", user.Balance, "\n")

	// Create a statement
	statement := Statement{Title: "Transferência de empréstimo", User: userID, Amount: amount, Type: types.Statement.Out}
	statement.Create()
}

func UserTake(userID string, lending *Lending) {
	// Find user in database
	user := GetUserById(userID)

	fmt.Println("User: ", user.ID)
	fmt.Println("Before: ", user.Balance)

	// Remove money
	user.Balance += lending.Amount
	user.Save()

	fmt.Println("After: ", user.Balance, "\n")

	// Create a statement
	statement := Statement{Title: "Transferência de empréstimo", User: userID, Amount: lending.Amount, Type: types.Statement.In}
	statement.Create()
}
