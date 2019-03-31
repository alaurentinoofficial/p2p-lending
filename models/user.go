package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"p2p-lending/types"
	"strings"
	"time"
)

type User struct {
	ID           string  `json:"id" gorm:"primary_key;"`
	Name         string  `json:"name"`
	Email        string  `json:"email"`
	Password     string  `json:"password"`
	CreationDate string  `json:"creation_date"`
	Type         int     `json:"types"`
	Score        int     `json:"score"`
	CpfCnpj      string  `json:"cpf_cpnj"`
	Balance      float32 `json:"balance"`
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

func (user *User) Create() bool {
	if user.Verify() {
		GetDB().Create(&user)
		return true
	} else {
		return false
	}
}

func (user *User) Save() bool {
	if user.Verify() {
		GetDB().Save(&user)
		return true
	} else {
		return false
	}
}

func (user *User) Verify() bool {
	isvalid := true

	isvalid = isvalid && strings.Contains(user.Email, "@")
	isvalid = isvalid && len(user.Password) > 6

	// Check if Physical or Legal person
	if user.Type == types.User.Physical {
		isvalid = isvalid && len(user.CpfCnpj) == 11
	} else {
		isvalid = isvalid && len(user.CpfCnpj) == 14
	}

	isvalid = isvalid && user.Score >= 0 && user.Score <= 1000

	isvalid = isvalid && len(user.State) > 0
	isvalid = isvalid && len(user.City) > 0
	isvalid = isvalid && len(user.Neighborhood) > 0
	isvalid = isvalid && user.Number > 0
	isvalid = isvalid && len(user.ZipCode) == 8

	return isvalid
}

func GetUserById(id string) *User {
	user := User{}
	GetDB().Table("users").Where("id = ?", id).First(&user)

	return &user
}

func UserCheckBalance(userID string, amount float32) bool {
	user := GetUserById(userID)
	return user.Balance-amount >= 0
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
