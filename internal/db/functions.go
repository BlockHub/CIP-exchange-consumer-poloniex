package db

import (
	"github.com/jinzhu/gorm"
	"time"
	"strings"
	"log"
)

func CreateGetMarket(gorm gorm.DB, quote string, ticker string) PoloniexMarket {
	market := PoloniexMarket{0, ticker, quote}
	err := gorm.Create(&market).Error
	if err != nil{
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			gorm.Where(map[string]interface{}{"ticker": ticker, "quote": quote}).Find(&market)
		} else {
			log.Panic(err)
		}
	}
	return market
}

func AddTicker(gorm gorm.DB, Ask float64, Bid float64, market PoloniexMarket) PoloniexTicker {
	ticker := PoloniexTicker{0, market.ID, Ask, Bid, time.Now()}
	err := gorm.Create(&ticker).Error
	if err != nil{
		log.Panic(err)
	}
	return ticker
}

func AddOrderBook(gorm gorm.DB, market PoloniexMarket) PoloniexOrderBook{
	book := PoloniexOrderBook{0, market.ID, time.Now()}
	err := gorm.Create(&book).Error
	if err != nil{
		log.Panic(err)
	}
	return book
}

func AddOrder(Gorm gorm.DB, Book PoloniexOrderBook, Rate float64, Quantity float64, Type int64, Buy bool){
	order := PoloniexOrder{0, Book.ID, Rate, Quantity, Type, Buy, time.Now()}
	err := Gorm.Create(&order).Error
	if err != nil{
		log.Panic(err)
	}
}