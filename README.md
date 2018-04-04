# CIP-exchange-consumer-poloniex

recommended go version: 1.10

required env variables:

- DB (e.g. postgres)
- DB_URL (see http://doc.gorm.io/database.html#connecting-to-a-database) 
- RAVEN_DSN (https://docs.sentry.io/clients/go/)
- PRODUCTION (false/true)
- R_DB=postgres
- R_DB_URL
- REPLICATION_LIMIT
- WORKER_SLEEP
- TICKER_SLEEP

Docker version: 18.03.0~ce-0~ubuntu

During prototyping we don't work with binary releases. Just create a binary using 

```ssh
    env GOOS=linux GOARCH=amd64 go build main.go
```

rename the binary to consumer.bin (for reasons of automatic deployment)