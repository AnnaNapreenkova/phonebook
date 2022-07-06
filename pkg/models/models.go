package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

type Number struct {
	ID      int
	Name    string
	Phone   string
	Created time.Time
}
