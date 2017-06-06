### pginit
Golang package that will create a database (if the database does not exist), create a user with a given password (if the user does not exists), and consequently grants all permissions to the user on that database.

### Usage
1. Install the package  
```
go get -u github.com/Unaxiom/pginit
```

2. [Ulogger](https://github.com/Unaxiom/ulogger) needs to be setup. It needs to be created and passed as 
```
log := ulogger.New()
pginit.Init(log)
```

3. Create the database
```
pginit.CreateDB(dbName)
```

4. Create the user
```
pginit.CreateUser("userA", "passwd", dbName)
```