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
	"github.com/getsentry/raven-go"
	"strconv"
	"CIP-exchange-consumer-poloniex/pkg/pushers"
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

func init(){
	useDotenv := true
	if os.Getenv("PRODUCTION") == "true"{
		useDotenv = false
	}

	// this loads all the constants stored in the .env file (not suitable for production)
	// set variables in supervisor then.
	if useDotenv {
		err := godotenv.Load()
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
	}
	raven.SetDSN(os.Getenv("RAVEN_DSN"))
}


func main() {

	localdb, err := gorm.Open(os.Getenv("DB"), os.Getenv("DB_URL"))
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
	}
	defer localdb.Close()

	localdb.AutoMigrate(&db.PoloniexTicker{}, &db.PoloniexMarket{}, &db.PoloniexOrder{}, &db.PoloniexOrderBook{})
	err = localdb.Exec("CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;").Error
	if err != nil{
		raven.CaptureErrorAndWait(err, nil)
	}
	err = localdb.Exec("SELECT create_hypertable('poloniex_orders', 'time', 'orderbook_id', if_not_exists => TRUE)").Error
	if err != nil{
		raven.CaptureErrorAndWait(err, nil)
	}
	err = localdb.Exec("SELECT create_hypertable('poloniex_tickers', 'time', 'market_id', if_not_exists => TRUE)").Error
	if err != nil{
		raven.CaptureErrorAndWait(err, nil)
	}
	err =localdb.Exec("SELECT create_hypertable('poloniex_order_books', 'time', 'market_id', if_not_exists => TRUE)").Error
	if err != nil{
		raven.CaptureErrorAndWait(err, nil)
	}
	localdb.DB().SetMaxOpenConns(1000)

	remotedb, err := gorm.Open(os.Getenv("R_DB"), os.Getenv("R_DB_URL"))
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
	}
	defer remotedb.Close()

	//start a replication worker
	limit,  err:= strconv.ParseInt(os.Getenv("REPLICATION_LIMIT"), 10, 64)
	replicator := pushers.Replicator{Local:*localdb, Remote:*remotedb, Limit:limit}
	go replicator.Start()

	restClient := poloniex.NewPublicOnly()
	// get the ticker so we know all available markets
	tickers, err := restClient.Ticker()
		if err != nil{
			raven.CaptureErrorAndWait(err, nil)
	}
	ws := poloniex.NewWithCredentials("", "")

	markets := make(map[string]db.PoloniexMarket)
	for key, _ := range tickers{
		// Market Ids are formatted with an underscore e.g. BTC_ETH
		s := strings.Split(key, "_")
		market := db.CreateGetMarket(*localdb, s[0], s[1])
		// get the initial orderbook state and generate a new orderbook row in the db
		book := init_orderbook(*restClient, market, *localdb)
		// the tickerhandler needs to be able to quickly assign foreign key to ticker values,
		// so we pass a map of market structs and market ids
		markets[key] = market
		ws.Subscribe(key)
		handler := handlers.OrderHandler{book, *localdb}
		ws.On(key, handler.Handle)
	}
	tickerhandler := handlers.TickerHandler{markets, *localdb}
	ws.Subscribe("ticker")
	ws.On("ticker", tickerhandler.Handle)
	ws.StartWS()
	// block the main routine
	select{}
}
