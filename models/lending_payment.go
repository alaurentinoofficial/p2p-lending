package models

import (
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

type LendingPayment struct {
	ID            string  `json:"id" gorm:"primary_key;"`
	Lending       string  `json:"lending"`
	Taker         float32 `json:"taker"`
	Validate      string  `json:"validate"`
	Value         float32 `json:"value"`
	Portion       int     `json:"Portion"`
	MonthlyDelays int     `json:"monthly_delays"`
	Status        bool    `json:"status"`
	PaymentDate   string  `json:"payment_day"`
}

func (payment *LendingPayment) BeforeCreate(scope *gorm.Scope) error {
	uu, _ := uuid.NewV4()
	_ = scope.SetColumn("ID", uu.String())
	_ = scope.SetColumn("PaymentDate", "")
	_ = scope.SetColumn("MonthlyDelays", 0)
	_ = scope.SetColumn("Status", false)
	return nil
}
