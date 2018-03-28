package db

import "time"

type PoloniexOrderBook struct {
	ID uint 			`gorm:"primary_key"`
	MarketID uint
	Time time.Time		`gorm:"primary_key"`
}

type PoloniexOrder struct {
	ID uint 			`gorm:"primary_key"`
	OrderbookID uint
	Rate float64
	Quantity float64
	// 0: initial, 1: orderBookModify, 2: orderBookRemove
	Type int64
	Buy bool
	Time time.Time		`gorm:"primary_key"`
}

type PoloniexMarket struct {
	ID uint 			`gorm:"primary_key"`
	Ticker string		`gorm:"unique_index:idx_market"`
	Quote string		`gorm:"unique_index:idx_market"`
}

type PoloniexTicker struct {
	ID  uint 			`gorm:"primary_key"`
	MarketID uint
	Ask float64
	Bid float64
	Time time.Time		`gorm:"primary_key"`
}
