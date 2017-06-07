### pginit
Golang package that will create a database (if the database does not exist), create a user with a given password (if the user does not exists), and consequently grants all permissions to the user on that database.

### Usage

1. [Ulogger](https://github.com/Unaxiom/ulogger) needs to be setup. It needs to be created and passed as 
```
pginit.Init(appName, orgName, production) // appName and orgName are strings. production is bool.
```

2. Create the database
```
pginit.CreateDB(dbName)
```

3. Create the user
```
pginit.CreateUser("userA", "passwd", dbName)
```