package db

import "time"

type PoloniexOrderBook struct {
	ID uint64 			`gorm:"primary_key"`
	MarketID uint64
	Time time.Time		`gorm:"primary_key"`
}

type PoloniexOrder struct {
	ID uint64 			`gorm:"primary_key"`
	OrderbookID uint64
	Rate float64
	Quantity float64
	// 0: initial, 1: orderBookModify, 2: orderBookRemove
	Type int64
	Buy bool
	Time time.Time		`gorm:"primary_key"`
}

type PoloniexMarket struct {
	ID uint64 			`gorm:"primary_key"`
	Ticker string		`gorm:"unique_index:polo_idx_market"`
	Quote string		`gorm:"unique_index:polo_idx_market"`
}

type PoloniexTicker struct {
	ID  uint64 			`gorm:"primary_key"`
	MarketID uint64
	Ask float64
	Bid float64
	Time time.Time		`gorm:"primary_key"`
}

type PoloniexTrade struct {
	ID  uint64 			`gorm:"primary_key"`
	MarketID uint64
	Rate float64
	Amount float64
	Total float64
	Buy bool
	Time time.Time		`gorm:"primary_key"`
}
