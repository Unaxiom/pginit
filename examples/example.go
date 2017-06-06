package main

import (
	"github.com/Unaxiom/pginit"
	"github.com/Unaxiom/ulogger"
)

func main() {
	log := ulogger.New()
	pginit.Init(log)
	dbName := "play5"
	pginit.CreateDB(dbName)
	pginit.CreateUser("userA", "passwd", dbName)
}
