package models

import (
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"time"
)

type Statement struct {
	ID            string  `json:"id" gorm:"primary_key;"`
	Title         string  `json:"Title"`
	User          string  `json:"user"`
	Type          int     `json:"types"`
	Amount        float32 `json:"amount"`
	OperationDate string  `json:"operation_date"`
}

func (statement *Statement) BeforeCreate(scope *gorm.Scope) error {
	uu, _ := uuid.NewV4()
	_ = scope.SetColumn("ID", uu.String())
	_ = scope.SetColumn("OperationDate", time.Now().UTC().String())
	return nil
}

func (statement *Statement) Create() {
	GetDB().Create(&statement)
}

func (statement *Statement) Save() {
	GetDB().Save(&statement)
}

func GetStatementByUser(userID string) []*Statement {
	statements := []*Statement{}
	GetDB().Table("statement").Where("user = ?", userID).Find(&statements)
	return statements
}
