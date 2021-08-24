package main

import "time"

type BlocketAd struct {
	ID      string    `gorm:"primarykey" json:"Id,omitempty"`
	Created time.Time `json:"CreatedAt,omitempty"`
	Subject string    `json:"Subject"`
	Body    string    `json:"Body"`
	Email   string    `json:"Email"`
	Price   *float64  `json:"Price,omitempty"`
}
