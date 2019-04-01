package models

import (
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"p2p-lending/types"
	"time"
)

type Lender struct {
	ID            string  `json:"id" gorm:"primary_key;"`
	Lending       string  `json:"lending"`
	User          string  `json:"user"`
	Amount        float32 `json:"amount"`
	OperationDate float32 `json:"operation_date"`
	Status        bool    `json:"status"`
}

func (lender *Lender) BeforeCreate(scope *gorm.Scope) error {
	uu, _ := uuid.NewV4()
	_ = scope.SetColumn("ID", uu.String())
	_ = scope.SetColumn("OperationDate", time.Now().UTC().Format(time.RFC3339))
	return nil
}

func (lender *Lender) AfterCreate(scope *gorm.Scope) error {
	scope.CommitOrRollback()

	lending := GetLendingById(lender.Lending)
	lending.AlreadyInvested += lender.Amount
	lending.Save()

	lending.Transfer()

	return nil
}

func (lender *Lender) Verify() int {
	isvalid := true
	isvalid = isvalid && lender.Amount > float32(0)

	lending := GetLendingById(lender.Lending)
	user := GetUserById(lender.User)

	if lending.ID == "" || user.ID == "" || !isvalid {
		return types.Response.InvalidArguments
	}

	if lending.Amount - lending.AlreadyInvested >= lender.Amount {
		if user.Balance - lender.Amount >= 0 {
			return types.Response.Ok
		}

		return types.Response.InsufficientFunds
	}

	return types.Response.PaymentCeiling
}

func (lender *Lender) Create() int {
	result := lender.Verify()

	if result == types.Response.Ok {
		GetDB().Create(&lender)
	}

	return result
}

func (lender *Lender) Save() int {
	result := lender.Verify()

	if result == types.Response.Ok {
		GetDB().Save(&lender)
	}

	return result
}

func GetLenderById(id string) *Lender {
	lender := Lender{}
	GetDB().Table("lenders").Where("id = ?", id).First(&lender)

	return &lender
}

func GetLendersByUser(userID string) []*Lender {
	var lenders []*Lender
	GetDB().Where(&Lender{User: userID}).Find(&lenders)
	return lenders
}

func GetLendersByLending(lendingID string) []*Lender {
	var lenders []*Lender
	GetDB().Table("lenders").Where("lending = ?", lendingID).Find(&lenders)

	return lenders
}

func DeleteLender(id string, lenderingID string, lending *Lending) {

	if lending == nil {
		lending = GetLendingById(lenderingID)
	}

	lending.AlreadyInvested -= GetLenderById(id).Amount
	lending.Save()

	GetDB().Table("lenders").Where("id = ?", id).Delete(&Lender{})
}
