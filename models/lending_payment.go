package models

import (
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"time"
)

type LendingPayment struct {
	ID            string  `json:"id" gorm:"primary_key;"`
	Lending       string  `json:"lending"`
	Taker         string  `json:"taker"`
	Validate      string  `json:"validate"`
	Value         float32 `json:"value"`
	Total         float32 `json:"total"`
	Portion       int     `json:"Portion"`
	MonthlyDelays int     `json:"monthly_delays"`
	Status        bool    `json:"status"`
	PaymentDate   string  `json:"payment_day"`
	LastPortion   bool    `json:"last_portion"`
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
	isvalid = isvalid && len(payment.Taker) == len("dc5ccc85-c1ee-41b0-92a5-bd7bae46ad35")
	isvalid = isvalid && payment.Portion >= 1
	isvalid = isvalid && payment.MonthlyDelays >= 0
	isvalid = isvalid && payment.Value >= 0

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

func (payment *LendingPayment) Pay() {
	payment.PaymentDate = time.Now().UTC().Format(time.RFC3339)
	payment.Status = true
	payment.Save()

	total := payment.Value
	lenders := GetLendersByLending(payment.Lending)

	for _, lender := range lenders {
		UserInsertMoney(lender.User, float32(Round(float64(total*(lender.Amount/payment.Total)), .5, 2)), "Recebimento do empréstimo + Juros")
	}
}

func GetLendingPayment(id string) *LendingPayment {
	lendingPayment := LendingPayment{}
	GetDB().Table("lending_payments").Where("id = ?", id).First(&lendingPayment)
	return &lendingPayment
}

func GetLendingPaymentsByTaker(takerID string) []*LendingPayment {
	lendingPayments := []*LendingPayment{}
	GetDB().Table("lending_payments").Where("taker = ?", takerID).Find(&lendingPayments)
	return lendingPayments
}
