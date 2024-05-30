package models

type Payment struct {
	Transaction  string
	RequestID    string
	Currency     string
	Provider     string
	Amount       int64
	PaymentDT    int64
	Bank         string
	DeliveryCost int64
	GoodsTotal   int64
	CustomFee    int64
}
