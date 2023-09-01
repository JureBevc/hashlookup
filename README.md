# hashlookup

# How to run

Create tables in database
```
go run .\main.go create-tables
```

Drop tables from database
```
go run .\main.go delete-tables
```

Create lookup table
```
go run .\main.go create-lookup -h [algorithm] -f [filepath]
```

Start API
```
go run .\main.go start-api
```

Build executable
```
go build
```

# Data sources

```
1. https://github.com/danielmiessler/SecLists/blob/master/Passwords/Common-Credentials/10-million-password-list-top-1000000.txt
2. https://github.com/danielmiessler/SecLists/blob/master/Passwords/xato-net-10-million-passwords.txt
3. https://xato.net/today-i-am-releasing-ten-million-passwords-b6278bbe7495
4. https://github.com/piotrcki/wordlist/releases
```