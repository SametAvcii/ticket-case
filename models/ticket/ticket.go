package models

import "gorm.io/gorm"

type Ticket struct {
	gorm.Model
	Name       string `json:"name"`
	Desc       string `json:"desc"`
	Allocation uint   `json:"allocation"`
}
