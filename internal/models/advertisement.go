package models

import "time"

// Advertisement представляет структуру parsed объявления
type Advertisement struct {
	Title           string
	Link            string
	Description     string
	Type            string
	ViewsCount      int
	City            string
	PublicationDate time.Time
	UserName        string
	UserPhone       string
	UserOnlineDate  time.Time
	UserID          string
}
