package models

import "time"

type Order struct {
	UID               string
	TrackNumber       string
	Entry             string
	Delivery          Delivery
	Payment           Payment
	Items             []Item
	Locale            string
	InternalSignature string
	CustomerID        string
	DeliveryService   string
	Shardkey          string
	SmID              int64
	DateCreated       time.Time
	OofShard          string
}
