package models

import (
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"os"
	"strconv"
	"time"
)

type LendingPayment struct {
	ID                  string  `json:"id" gorm:"primary_key;"`
	Lending             string  `json:"lending"`
	Validate            string  `json:"validate"`
	Portion             int     `json:"Portion"`
	Status              bool    `json:"status"`
	PaymentDate         string  `json:"payment_day"`
	LastPortion         bool    `json:"last_portion"`
}

func (payment *LendingPayment) BeforeCreate(scope *gorm.Scope) error {
	uu, _ := uuid.NewV4()
	_ = scope.SetColumn("ID", uu.String())
	_ = scope.SetColumn("PaymentDate", "")
	_ = scope.SetColumn("MonthlyDelays", 0)
	_ = scope.SetColumn("Status", false)
	return nil
}

func (payment *LendingPayment) Verify() bool {
	isvalid := true

	isvalid = isvalid && len(payment.Lending) == len("dc5ccc85-c1ee-41b0-92a5-bd7bae46ad35")
	isvalid = isvalid && payment.Portion >= 1

	return isvalid
}

func (payment *LendingPayment) Create() bool {
	if payment.Verify() {
		GetDB().Create(payment)
		return true
	} else {
		return false
	}
}

func (payment *LendingPayment) Save() bool {
	if payment.Verify() {
		GetDB().Save(payment)
		return true
	} else {
		return false
	}
}

func (payment *LendingPayment) Pay(lending Lending) {
	payment.PaymentDate = time.Now().UTC().Format(time.RFC3339)
	payment.Status = true
	payment.Save()

	if payment.LastPortion {
		admTaxe, _ := strconv.ParseFloat(os.Getenv("ADM_TAXE"), 32)
		total := lending.PortionAmount * float32(payment.Portion)
		total = lending.Amount + ((total - lending.Amount) * float32(admTaxe))
		lenders := GetLendersByLending(payment.Lending)

		for _, lender := range lenders {
			UserInsertMoney(lender.User, float32(Round(float64(total*(lender.Amount/lending.Amount)), .5, 2)), "Recebimento do empréstimo + Juros")
		}
	}
}

func GetLendingPayment(id string) *LendingPayment {
	lendingPayment := LendingPayment{}
	GetDB().Table("lending_payments").Where("id = ?", id).First(&lendingPayment)
	return &lendingPayment
}

func GetLendingPaymentsByLending(lendingID string) []*LendingPayment {
	lendingPayments := []*LendingPayment{}
	GetDB().Table("lending_payments").Where("lending = ?", lendingID).Find(&lendingPayments)
	return lendingPayments
}

func monthsCountSince(createdAtTime time.Time) int {
	now := time.Now().UTC()
	months := 0
	month := createdAtTime.Month()
	days := 0

	for createdAtTime.Before(now) {
		createdAtTime = createdAtTime.Add(time.Hour * 24)
		nextMonth := createdAtTime.Month()
		if nextMonth != month {
			months++
			days = 0
		} else {
			days++
		}
		month = nextMonth
	}

	if months == 0 && days > 0 {
		return 1
	}

	return months
}
