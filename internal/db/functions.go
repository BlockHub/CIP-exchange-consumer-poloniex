package db

import (
	"github.com/jinzhu/gorm"
	"time"
	"strings"
)

func CreateGetMarket(gorm gorm.DB, quote string, ticker string) PoloniexMarket {
	market := PoloniexMarket{0, ticker, quote}
	err := gorm.Create(&market).Error
	if err != nil{
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			gorm.Where(map[string]interface{}{"ticker": ticker, "quote": quote}).Find(&market)
		}
	}
	return market
}

func AddTicker(gorm gorm.DB, Ask float64, Bid float64, market PoloniexMarket) PoloniexTicker {
	ticker := PoloniexTicker{0, market.ID, Ask, Bid, int64(time.Now().Unix())}
	err := gorm.Create(&ticker).Error
	if err != nil{
		panic(err)
	}
	return ticker
}

func AddOrderBook(gorm gorm.DB, market PoloniexMarket) PoloniexOrderBook{
	book := PoloniexOrderBook{0, market.ID, int64(time.Now().Unix())}
	err := gorm.Create(&book).Error
	if err != nil{
		panic(err)
	}
	return book
}

func AddOrder(Gorm gorm.DB, Book PoloniexOrderBook, Rate float64, Quantity float64, Type int64, Buy bool){
	order := PoloniexOrder{0, Book.ID, Rate, Quantity, Type, Buy, int64(time.Now().Unix())}
	err := Gorm.Create(&order).Error
	if err != nil{
		panic(err)
	}
}