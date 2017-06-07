package main

import (
	"github.com/Unaxiom/pginit"
)

func main() {
	pginit.Init("App", "Org", false)
	dbName := "play5"
	pginit.CreateDB(dbName)
	pginit.CreateUser("userA", "passwd", dbName)
}
