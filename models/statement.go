package models

import (
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"p2p-lending/types"
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

func (statement *Statement) Verify() bool {
	isvalid := true

	isvalid = isvalid && len(statement.Title) > 0
	isvalid = isvalid && types.Statement.Check(statement.Type)
	isvalid = isvalid && statement.Amount > 0

	return isvalid
}

func (statement *Statement) Create() bool {
	if statement.Verify() {
		GetDB().Create(&statement)
		return true
	} else {
		return false
	}
}

func (statement *Statement) Save() bool {
	if statement.Verify() {
		GetDB().Save(&statement)
		return true
	} else {
		return false
	}
}

func GetStatementByUser(userID string) []*Statement {
	statements := []*Statement{}
	GetDB().Table("statement").Where("user = ?", userID).Find(&statements)
	return statements
}
