package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"p2p-lending/types"
	"time"
)

type Lending struct {
	ID              string  `json:"id" gorm:"primary_key;"`
	Taker           string  `json:"taker"`
	Amount          float32 `json:"amount"`
	AlreadyInvested float32 `json:"already_invested"`
	Status          bool    `json:"status"`
	CreationDate    string  `json:"creation_date"`
	Validate        string  `json:"validate"`
	TransactionDate string  `json:"transaction_date"`
	HasIndex        bool    `json:"has_index"`
	Index           int     `json:"index"`
	IndexYield      float32 `json:"index_yield"`
	PrefixedYield   float32 `json:"prefixed_yield"`
	TimeMonth       int     `json:"time_month"`
}

func (lending *Lending) BeforeCreate(scope *gorm.Scope) error {
	uu, _ := uuid.NewV4()
	_ = scope.SetColumn("ID", uu.String())
	_ = scope.SetColumn("Status", false)
	_ = scope.SetColumn("CreationDate", time.Now().UTC().Format(time.RFC3339))
	_ = scope.SetColumn("Validate", time.Now().UTC().AddDate(0,1, 0).Format(time.RFC3339))
	return nil
}

func (lending *Lending) Verify() bool {
	isvalid := true

	isvalid = isvalid && len(lending.Taker) == len("dc5ccc85-c1ee-41b0-92a5-bd7bae46ad35")
	isvalid = isvalid && lending.Amount > 1000

	if lending.HasIndex {
		isvalid = isvalid && types.Index.Check(lending.Index)
		isvalid = isvalid && lending.IndexYield > 0
	} else {
		isvalid = isvalid && lending.PrefixedYield > 0
	}

	isvalid = isvalid && lending.TimeMonth > 1

	return isvalid
}

func (lending *Lending) Create() bool {
	if lending.Verify() {
		GetDB().Create(&lending)
		return true
	} else {
		return false
	}
}

func (lending *Lending) Save() bool {
	if lending.Verify() {
		GetDB().Save(&lending)
		return true
	} else {
		return false
	}
}

func (lending *Lending) Transfer() bool {
	// Get all Lenders from database
	lenders := GetLendersByLending(lending.ID)

	totalAmount := float32(0)
	for _, lender := range lenders {
		// Check if the all lenders has balance
		if UserCheckBalance(lender.User, lender.Amount) {
			// Sum in the total
			totalAmount += lender.Amount
		} else {
			// Delete lender
			DeleteLender(lender.ID, lending.ID, lending)
			return false
		}
	}

	// Check all money received and check the date
	if totalAmount == lending.Amount {
		// Convert the validate date to time.Date
		validate, _ := time.Parse(time.RFC3339, lending.Validate)

		if time.Now().UTC().Before(validate) {
			// Reduce balance from users
			for _, lender := range lenders {
				UserLend(lender.User, lender.Amount, lending)

				lender.Status = true
				lender.Save()
			}

			// Transfer to taker
			UserTake(lending.Taker, lending)

			// Save configurations
			lending.Status = true
			lending.TransactionDate = time.Now().UTC().String()
			lending.Save()
			return true
		} else if time.Now().UTC().After(validate) {
			DeleteLending(lending.ID)
		}
	}

	return false
}

func GetLendingById(id string) *Lending {
	lending := Lending{}
	GetDB().Table("lendings").First(&lending)

	return &lending
}

func DeleteLending(id string) {
	fmt.Println("[*] Lending Expired!")
	GetDB().Table("lendings").Where("id=?", id).Delete(&Lending{})
}
