package models

import (
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
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
	HasIndexer      bool    `json:"has_indexer"`
	Indexer         string  `json:"indexer"`
	Yield           float32 `json:"yield"`
}

func (lending *Lending) BeforeCreate(scope *gorm.Scope) error {
	uu, _ := uuid.NewV4()
	_ = scope.SetColumn("ID", uu.String())
	_ = scope.SetColumn("CreationDate", time.Now().UTC().String())
	return nil
}

//func (lending *Lending) BeforeSave(scope *gorm.Scope) error {
//	fmt.Println(lending.Amount == lending.AlreadyInvested)
//	if !lending.Status && lending.Amount == lending.AlreadyInvested {
//		lenders := GetLendersByLending(lending.ID)
//
//		totalAmount := float32(0)
//		for _, lender := range lenders {
//			// Check if the all lenders has balance
//			if UserCheckBalance(lender.User, lender.Amount) {
//				// Sum in the total
//				totalAmount += lender.Amount
//			} else {
//				// Delete lender
//				DeleteLender(lender.ID, lending.ID, lending)
//			}
//		}
//
//		if totalAmount == lending.Amount {
//
//			// Reduce balance from users
//			for _, lender := range lenders {
//				UserLend(lender.User, lender.Amount, lending)
//
//				lender.Status = true
//				lender.Save()
//			}
//
//			// Transfer to taker
//			UserTake(lending.Taker, lending)
//
//			// Save configurations
//			lending.Status = true
//			lending.TransactionDate = time.Now().UTC().String()
//		}
//	}
//
//	return nil
//}

func (lending *Lending) Create() {
	GetDB().Create(&lending)
}

func (lending *Lending) Save() {
	GetDB().Save(&lending)
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

	// Convert the validate date to time.Date
	//validate, _ := time.Parse("2006-01-02T15:04:05.000Z", lending.Validate)

	// Check all money received and check the date
	if totalAmount == lending.Amount {

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
	} else {
		return false
	}
}

func GetLendingById(id string) *Lending {
	lending := Lending{}
	GetDB().Table("lendings").First(&lending)

	return &lending
}
