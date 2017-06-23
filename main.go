package pginit

import (
	"os/exec"
	"time"

	"github.com/Unaxiom/ulogger"

	"fmt"

	"strings"

	pgx "gopkg.in/jackc/pgx.v2"
)

var db *pgx.ConnPool
var dbErr error
var log *ulogger.Logger

// Init needs to be called to set up the logger object
func Init(appName string, orgName string, production bool) {
	log = ulogger.New()
	log.ApplicationName = appName + " pginit"
	log.OrganizationName = orgName
	log.RemoteAvailable = production
	if production {
		ulogger.RemoteURL = "https://logs.unaxiom.com/newlog"
	}
	// Check if the DB is live

	dbStatus := CheckIfDBIsLive()
	if !dbStatus {
		log.Warningln("Since DB isn't live, will try after a while...")
		<-time.After(time.Second * time.Duration(2))
		Init(appName, orgName, production)
	} else {
		log.Infoln("PostgreSQL server is live!")
	}
}

// CheckIfDBIsLive returns a bool stating if the postgres server is live
func CheckIfDBIsLive() bool {
	command := fmt.Sprint("SELECT EXISTS(SELECT 1)")
	out, err := runSQL(command)
	if err != nil {
		log.Warningln("DB isn't live.")
		return false
	}
	exists := parseExistence(out)
	return exists
}

// CreateDB creates the database, if it doesn't exist
func CreateDB(dbName string) {
	// command := fmt.Sprint("psql -c \"CREATE DATABASE hey_there\"")
	exists := checkIfDBExists(dbName)
	if exists {
		log.Infoln("Db --> ", dbName, " already exists.")
		return
	}
	log.Infoln("DB --> ", dbName, " does not exist. Need to create it.")
	// Create the database here
	command := fmt.Sprintf("CREATE DATABASE %s", dbName)
	out, err := runSQL(command)
	if err != nil {
		log.Fatalln("Couldn't create Database --> ", dbName, ". Error is --> ", err.Error(), ". Output is: \n", out)
	}
	log.Infoln("Created database --> ", dbName, ". Output of the command is:\n\n", out)

}

// CreateUser creates the user, if it doesn't exist
func CreateUser(username string, password string, dbName string) {
	exists := checkIfUserExists(username)
	if exists {
		log.Infoln("User --> ", username, " already exists.")
		// Instead of returning straight away, grant permissons to the user
		// on the mentioned database as well. Because, there's a chance that
		// the database might be newly created, but the user may already exist.
		// Thus, the existing user needs to have permissions on the newly
		// created database.
		// return
	} else {
		log.Infoln("User --> ", username, " does not exist. Need to create it.")
		// Create the user here
		command := fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s'", username, password)
		out, err := runSQL(command)
		if err != nil {
			log.Fatalln("Couldn't create user --> ", username, ". Error is --> ", err.Error(), ". Output is: \n", out)
		}
		log.Infoln("Created user --> ", username, ". Output of the command is:\n\n", out)
	}

	grantAllPermissionsToUser(username, dbName)
}

// checkIfUserExists accepts the username and returns if the user exists
func checkIfUserExists(username string) bool {
	// SELECT EXISTS(SELECT * FROM pg_catalog.pg_user WHERE usename = 'abc')
	command := fmt.Sprintf("SELECT EXISTS(SELECT * FROM pg_catalog.pg_user WHERE usename = '%s')", username)
	out, err := runSQL(command)
	if err != nil {
		log.Fatalln("Couldn't check if user --> ", username, " exists. Error is --> ", err.Error(), ". Output is: \n", out)
	}
	exists := parseExistence(out)
	return exists
}

// checkIfDBExists accepts the database name and returns if it exists
func checkIfDBExists(dbName string) bool {
	// SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = 'play')
	command := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = '%s')", dbName)
	out, err := runSQL(command)
	if err != nil {
		log.Fatalln("Couldn't check if database --> ", dbName, " exists. Error is --> ", err.Error(), ". Output is: \n", out)
	}
	exists := parseExistence(out)
	return exists
}

// grantAllPermissionsToUser grants all permissions to the specific user on the specific database
func grantAllPermissionsToUser(username string, dbName string) {
	// Need to check if the user already has permissions on this table
	command := fmt.Sprintf("GRANT ALL ON DATABASE %s TO %s", dbName, username)
	out, err := runSQL(command)
	if err != nil {
		log.Fatalln("Couldn't grant permissions on database --> ", dbName, " to user --> ", username, ". Error is --> ", err.Error(), ". Output is: \n", out)
	}
	log.Infoln("Granted all permissions on database --> ", dbName, " to user --> ", username, ". Output of the command is:\n\n", out)
}

// runSQL accepts a psql command, runs it, and returns the combined output along with an error, if any
func runSQL(command string) (string, error) {
	commandToRun := fmt.Sprintf(`psql -c "%s"`, command)
	cmd := exec.Command("sudo", "-u", "postgres", "bash", "-c", commandToRun)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// parseExistence accepts a string and returns true if `t` is present or false if `f` is present
func parseExistence(str string) bool {
	strList := strings.Split(str, "\n")
	if len(strList) < 3 {
		return false
	}
	str = strings.TrimSpace(strList[2])
	if str == "t" {
		return true
	}
	return false
}
