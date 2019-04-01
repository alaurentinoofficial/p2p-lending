package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"math"
	"time"
)

type Lending struct {
	ID                  string  `json:"id" gorm:"primary_key;"`
	Taker               string  `json:"taker"`
	Amount              float32 `json:"amount"`
	AlreadyInvested     float32 `json:"already_invested"`
	Status              bool    `json:"status"`
	CreationDate        string  `json:"creation_date"`
	Validate            string  `json:"validate"`
	TransactionDate     string  `json:"transaction_date"`
	PrefixedYield       float32 `json:"prefixed_yield"`
	MonthlyInterestRate float32 `json:"monthly_interest_rate"`
	PortionAmount       float32 `json:"portion_amount"`
	PortionAlreadyPayed int     `json:"portion_already_payed"`
	PaymentTimeMonth    int     `json:"payment_time_month"`
	//HasIndex            bool    `json:"has_index"`
	//Index               int     `json:"index"`
	//IndexYield          float32 `json:"index_yield"`
}

func (lending *Lending) BeforeCreate(scope *gorm.Scope) error {
	total := math.Round(float64(lending.Amount * (1 + lending.PrefixedYield/100)))
	monthlyPayment := Round(total/float64(lending.PaymentTimeMonth), .5, 2)

	uu, _ := uuid.NewV4()
	_ = scope.SetColumn("ID", uu.String())
	_ = scope.SetColumn("PortionAmount", float32(monthlyPayment))
	_ = scope.SetColumn("Status", false)
	_ = scope.SetColumn("AlreadyInvested", 0)
	_ = scope.SetColumn("PortionAlreadyPayed", 0)
	_ = scope.SetColumn("CreationDate", time.Now().UTC().Format(time.RFC3339))
	_ = scope.SetColumn("Validate", time.Now().UTC().AddDate(0, 1, 0).Format(time.RFC3339))
	return nil
}

func (lending *Lending) Verify() bool {
	isvalid := true

	isvalid = isvalid && len(lending.Taker) == len("dc5ccc85-c1ee-41b0-92a5-bd7bae46ad35")
	isvalid = isvalid && lending.Amount > 1000

	isvalid = isvalid && lending.PrefixedYield > 0
	//if lending.HasIndex {
	//	isvalid = isvalid && types.Index.Check(lending.Index)
	//	isvalid = isvalid && lending.IndexYield > 0
	//} else {
	//	isvalid = isvalid && lending.PrefixedYield > 0
	//}

	isvalid = isvalid && lending.PaymentTimeMonth > 1

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
				UserRemoveMoney(lender.User, lender.Amount, "Transferência de empréstimo")

				lender.Status = true
				lender.Save()
			}

			// Transfer to taker
			UserInsertMoney(lending.Taker, lending.Amount, "Transferência de empréstimo")

			// Save configurations
			lending.Status = true
			lending.TransactionDate = time.Now().UTC().String()
			lending.Save()

			for i := 1; i <= lending.PaymentTimeMonth; i++ {
				payment := LendingPayment{
					Lending:             lending.ID,
					Portion:             i,
					LastPortion:         i == lending.PaymentTimeMonth,
					Validate:            time.Now().UTC().AddDate(0, i, 0).Format(time.RFC3339),
				}
				payment.Create()
			}

			return true
		} else if time.Now().UTC().After(validate) {
			DeleteLending(lending.ID)
		}
	}

	return false
}

func (lending *Lending) CalculatePrice(validate time.Time) float32 {
	months := monthsCountSince(validate)

	if months == 0 {
		// Normal Value
		return lending.PortionAmount
	} else {
		return lending.PortionAmount * float32(math.Pow(float64(lending.MonthlyInterestRate/100+1.0), float64(months)))
	}
}

func GetLendingById(id string) *Lending {
	lending := Lending{}
	GetDB().Table("lendings").First(&lending)

	return &lending
}

func GetLendings() []*Lending {
	lendings := []*Lending{}
	GetDB().Table("lendings").Find(&lendings)

	return lendings
}

func DeleteLending(id string) {
	fmt.Println("[*] Lending Expired!")
	GetDB().Table("lendings").Where("id=?", id).Delete(&Lending{})
}

func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}
