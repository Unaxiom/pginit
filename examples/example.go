package main

import (
	"git.unaxiom.com/pginit"
)

func main() {
	pginit.Init("App", "Org", false)
	dbName := "play5"
	pginit.CreateDB(dbName)
	pginit.CreateUser("usera", "passwd", dbName)
}
