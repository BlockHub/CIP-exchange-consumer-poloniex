package main

import (
	"github.com/pharrisee/poloniex-api"
	//"github.com/k0kubun/pp"
	//"time"
	"github.com/joho/godotenv"
	"log"
	"CIP-exchange-consumer-poloniex/pkg/handlers"
	"CIP-exchange-consumer-poloniex/internal/db"
	"github.com/jinzhu/gorm"
	"os"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"strings"
	"fmt"
)

// uses rest api to create new orderbook and passes orderbook instances
func init_orderbook(client poloniex.Poloniex , market db.PoloniexMarket, gorm gorm.DB) db.PoloniexOrderBook{
	MarketId := market.Quote + "_" + market.Ticker
	orders, err := client.OrderBook(MarketId)
	if err != nil{
		panic(err)
	}
	book := db.AddOrderBook(gorm, market)
	for _, order := range orders.Asks {
		db.AddOrder(gorm, book, order.Amount, order.Rate, 0, false)
	}
	for _, order := range orders.Bids {
		db.AddOrder(gorm, book, order.Amount, order.Rate, 0, true)
	}
	return book
}


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	gormdb, err := gorm.Open(os.Getenv("DB"), os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}
	defer gormdb.Close()

	gormdb.AutoMigrate(&db.PoloniexTicker{}, &db.PoloniexMarket{}, &db.PoloniexOrder{}, &db.PoloniexOrderBook{})
	err = gormdb.Exec("CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;").Error
	if err != nil{
		panic(err)
	}
	err = gormdb.Exec("SELECT create_hypertable('poloniex_orders', 'time', 'orderbook_id', if_not_exists => TRUE)").Error
	if err != nil{
		panic(err)
	}
	err = gormdb.Exec("SELECT create_hypertable('poloniex_tickers', 'time', 'market_id', if_not_exists => TRUE)").Error
	if err != nil{
		panic(err)
	}
	err =gormdb.Exec("SELECT create_hypertable('poloniex_order_books', 'time', 'market_id', if_not_exists => TRUE)").Error
	if err != nil{
		panic(err)
	}
	gormdb.DB().SetMaxOpenConns(1000)


	restClient := poloniex.NewPublicOnly()
	// get the ticker so we know all available markets
	tickers, err := restClient.Ticker()
		if err != nil{
		panic(err)
	}
	ws := poloniex.NewWithCredentials("", "")

	markets := make(map[string]db.PoloniexMarket)
	for key, _ := range tickers{
		// Market Ids are formatted with an underscore e.g. BTC_ETH
		s := strings.Split(key, "_")
		market := db.CreateGetMarket(*gormdb, s[0], s[1])
		// get the initial orderbook state and generate a new orderbook row in the db
		book := init_orderbook(*restClient, market, *gormdb)
		// the tickerhandler needs to be able to quickly assign foreign key to ticker values,
		// so we pass a map of market structs and market ids
		markets[key] = market
		ws.Subscribe(key)
		handler := handlers.OrderHandler{book, *gormdb}
		ws.On(key, handler.Handle)
	}
	tickerhandler := handlers.TickerHandler{markets, *gormdb}
	ws.Subscribe("ticker")
	ws.On("ticker", tickerhandler.Handle)
	fmt.Println("starting websocket")
	ws.StartWS()
	// block the main routine
	select{}
}
