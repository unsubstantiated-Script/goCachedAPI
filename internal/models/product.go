package models

type Product struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}
