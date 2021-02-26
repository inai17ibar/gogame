package model

import (
	"time"
)

type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Password     string    `json:"password"` //いまは使わない
	Created_date time.Time `json:"created_date"`
}

type JWT struct {
	Token string `json:"token"`
}

type Error struct {
	Message string `json:"message"`
}
