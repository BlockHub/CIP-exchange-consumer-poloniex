package handlers

import (
	"github.com/pharrisee/poloniex-api"
	"CIP-exchange-consumer-poloniex/internal/db"
	"github.com/jinzhu/gorm"
)


type TickerHandler struct {
	// because we can only use a single handler to deal with the tickers,
	// we preload all the markets in the tickerhandler, to lessen the strain
	// on the db.
	Markets map[string]db.PoloniexMarket
	Gorm gorm.DB
}
	func (t TickerHandler) Handle(ticker poloniex.WSTicker, ){
	db.AddTicker(t.Gorm, ticker.Ask, ticker.Bid, t.Markets[ticker.Pair])
	}


type OrderHandler struct {
	Book db.PoloniexOrderBook
	Market db.PoloniexMarket
	Gorm gorm.DB
}
	func (oh OrderHandler)Handle(order poloniex.WSOrderbook){
		var buy bool

		if order.Event == "trade"{
			switch event := order.Type; event {
				case "sell": 	buy = false
				default:		buy = true
			}
			db.AddTrade(oh.Gorm, oh.Market, order.Rate, order.Amount, order.Total, buy)
			}

		var TypeId int64
		if order.Event == "modify"{
			TypeId = 1
		}
		if order.Event == "remove"{
			TypeId = 2
		}

		if order.Type == "bid"{
			buy = true
		}
		if order.Type == "ask"{
			buy=false
		}
		db.AddOrder(oh.Gorm, oh.Book, order.Rate, order.Amount, TypeId, buy)
	}
