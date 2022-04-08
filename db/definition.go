package db

import (
	"fmt"
	"go/types"
	"gorm.io/gorm"
	"log"
)

type DB[MODEL Sample] struct {
	db *gorm.DB
	Provider[MODEL]
}

type Sample interface {
	*gorm.DB | types.Nil
}

type Provider[MODEL Sample] interface {
	initial() error
	instance(key string) MODEL
}

func (db DB[Sample]) Initial() {
	err := db.initial()
	if err != nil {
		log.Fatal(fmt.Sprintf("db err occured:%s", err.Error()))
	}
}

func (db DB[Sample]) Instance(key string) (s Sample) {
	return db.Provider.instance(key)
}
