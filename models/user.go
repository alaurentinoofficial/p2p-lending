package models

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"p2p-lending/types"
	"strings"
	"time"
)

type Token struct {
	UserId string
	jwt.StandardClaims
}

type TokenResponse struct {
	Token string `json:"token"`
}

type User struct {
	ID           string  `json:"id" gorm:"primary_key;"`
	Name         string  `json:"name"`
	Email        string  `json:"email"`
	Password     string  `json:"password"`
	CreationDate string  `json:"creation_date"`
	Type         int     `json:"type"`
	Salary       float32 `json:"salary"`
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
	_ = scope.SetColumn("CreationDate", time.Now().UTC().Format(time.RFC3339))
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
	isvalid = isvalid && user.Salary >= 0

	isvalid = isvalid && len(user.State) > 0
	isvalid = isvalid && len(user.City) > 0
	isvalid = isvalid && len(user.Neighborhood) > 0
	isvalid = isvalid && user.Number > 0
	isvalid = isvalid && len(user.ZipCode) == 8

	return isvalid
}

func (user *User) Pay(paymentID string) int {
	payment := GetLendingPayment(paymentID)
	lending := GetLendingById(payment.Lending)
	price := payment.CalculatePrice()

	if lending.PortionAlreadyPayed+1 == payment.Portion {
		if user.Balance-price >= 0 {
			user.Balance = float32(Round(float64(user.Balance-price), .5, 2))
			user.Save()

			lending.PortionAlreadyPayed += 1
			lending.Save()

			payment.Pay()
			return types.Response.Ok
		}

		return types.Response.InsufficientFunds
	}

	return types.Response.PayPreviousPortions
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

func UserRemoveMoney(userID string, amount float32, title string) {
	// Find user in database
	user := GetUserById(userID)

	// Remove money
	user.Balance -= amount
	user.Save()

	// Create a statement
	statement := Statement{Title: title, User: userID, Amount: amount, Type: types.Statement.Out}
	statement.Create()
}

func UserInsertMoney(userID string, amount float32, title string) {
	// Find user in database
	user := GetUserById(userID)

	// Add money
	user.Balance += amount
	user.Save()

	// Create a statement
	statement := Statement{Title: title, User: userID, Amount: amount, Type: types.Statement.In}
	statement.Create()
}
