# hashlookup

A tool for creating reverse hash lookup tables written in Go. 

## How to use

### Database management
- Create tables
```
go run hashlookup.go create-tables
```

- Drop tables
```
go run hashlookup.go delete-tables
```

- Reset/clear tables (equivalent to calling drop and create tables)
```
go run hashlookup.go reset-tables
```

### Populating tables / generating data
- Create simple reverse lookup table
```
go run hashlookup.go create-lookup -alg [algorithm] -file [filepath]
```

- Create rainbow table
```
go run hashlookup.go create-rainbow -alg [algorithm] -file [filepath]
```

### Query the data

- Query the simple lookup table
```
go run hashlookup.go check-lookup -alg [algorithm] -hash [hash]
```

- Query the rainbow table
```
go run hashlookup.go check-rainbow -alg [algorithm] -hash [hash]
```

- You can also start an API to query via GET requests
```
go run hashlookup.go start-api
```

### Build
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