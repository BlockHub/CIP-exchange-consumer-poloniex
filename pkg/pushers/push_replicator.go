package pushers

import (
	"github.com/jinzhu/gorm"
	"log"
	"CIP-exchange-consumer-poloniex/internal/db"
	"strings"
	"fmt"
)



type Replicator struct {
	//Used for logging purposes
	Name string
	// local db
	Local gorm.DB

	//remote DB (the data warehouse)
	Remote gorm.DB
	DBlink string
	//schema related settings

	//replication related settings
	Limit int64	// max rows to be fetched from remote and inserted (should be as high as possible)

}
// send the initial Markets data to remote
func (r *Replicator) PushMarkets() {
	markets := []db.PoloniexMarket{}
	r.Local.Limit(r.Limit).Find(&markets)

	// we don't delete the local copies of the markets, as they are needed for FK relations
	// and don't take up much space
	for _, market := range markets {
		err := r.Remote.Create(&market).Error
		if err != nil {
			if ! strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				log.Panic(err)
			}
		}
	}
}
// Create a persistent dblink
func (r *Replicator) Link() {
	fmt.Println("name", r.Name)
	err := r.Remote.Exec(
		fmt.Sprintf(`SELECT dblink_connect('%s', '%s');`, r.Name, r.DBlink)).Error
	if err != nil{
		log.Panic(err)
	}
}

// close the persistent dblink
func (r *Replicator) Unlink(){
	err := r.Remote.Exec(
		fmt.Sprintf(`SELECT dblink_disconnect('%s');`, r.Name)).Error
	if err != nil{
		log.Panic(err)
	}
}

func (r *Replicator) SendOrders(){
	err := r.Remote.Exec(
		fmt.Sprintf(
			`INSERT INTO poloniex_orders (id, orderbook_id, rate, quantity, time, buy, type)
					SELECT *
					FROM dblink(
						'%s',
						' DELETE FROM poloniex_orders WHERE id in (SELECT id FROM poloniex_orders ORDER BY time ASC LIMIT %d) RETURNING id, orderbook_id, rate, quantity, time, buy, type;'
					) AS deleted (id INT, orderbook_id INT, rate NUMERIC, quantity NUMERIC, time TIMESTAMP, buy BOOLEAN, type INT)`, r.Name, r.Limit)).Error
	if err != nil{
		log.Panic(err)
	}

}

func (r *Replicator) SendTickers(){
	err := r.Remote.Exec(
		fmt.Sprintf(
			`INSERT INTO poloniex_tickers (id, market_id, bid, ask, time)
					SELECT *
					FROM dblink(
						'%s',
						' DELETE FROM poloniex_tickers WHERE id in (SELECT id FROM poloniex_tickers ORDER BY time ASC LIMIT %d) RETURNING id, market_id, ask, bid, time;'
					) AS deleted (id INT, market_id INT, ask NUMERIC, bid NUMERIC, time TIMESTAMP)`, r.Name, r.Limit)).Error
	if err != nil{
		log.Panic(err)
	}
}

func (r *Replicator) Start() {
	// an out interface to store lots of Order objects
	for true {
		r.SendTickers()
		r.SendOrders()
	}
}