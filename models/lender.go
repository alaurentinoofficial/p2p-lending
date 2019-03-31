package models

import (
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"time"
)

type Lender struct {
	ID            string  `json:"id"`
	Lending       string  `json:"lending"`
	User          string  `json:"user"`
	Amount        float32 `json:"amount"`
	OperationDate float32 `json:"operation_date"`
	Status        bool    `json:"status"`
}

func (lender *Lender) BeforeCreate(scope *gorm.Scope) error {
	lending := GetLendingById(lender.Lending)
	lending.AlreadyInvested += lender.Amount
	lending.Save()

	uu, _ := uuid.NewV4()
	_ = scope.SetColumn("ID", uu.String())
	_ = scope.SetColumn("OperationDate", time.Now().UTC().String())
	return nil
}

func (lender *Lender) AfterCreate(scope *gorm.Scope) error {
	lending := GetLendingById(lender.Lending)

	if lending.Amount == lending.AlreadyInvested {
		lending.Transfer()
	}

	return nil
}

func (lender *Lender) Create() {
	GetDB().Table("lenders").Create(&lender)
}

func (lender *Lender) Save() {
	GetDB().Table("lenders").Save(&lender)
}

func GetLenderById(id string) *Lender {
	lender := Lender{}
	GetDB().Table("lenders").Where("id = ?", id).First(&lender)

	return &lender
}

func GetLendersByLending(lendingID string) []*Lender {
	var lenders []*Lender
	GetDB().Table("lenders").Where("lending = ?", lendingID).Find(&lenders)

	return lenders
}

func DeleteLender(id string) {
	lender := GetLenderById(id)

	lending := GetLendingById(lender.Lending)
	lending.AlreadyInvested -= lender.Amount
	lending.Save()

	GetDB().Table("lenders").Where("id = ?", id).Delete(Lender{})
}
