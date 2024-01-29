![Go Report Card](https://goreportcard.com/badge/github.com/ildomm/eskalationszeit?cache=v1)

# Eskalationszeit
![N|Solid](https://static.coindesk.com/wp-content/uploads/2014/05/coinstackrpricecharts.png)

Golang project for currency prices history persistence, targeting keep it using time ranges.

# Time ranges accepted, in minutes:
 1, 5, 15, 30, 60, 120, 720, 1440, 43800, 262800


# Components
### RESTApi server price generator
  - Answer random price for requested currency
  - No database connections
  - Optional, just need here to prove a point
  - location: /preisgenerator

### Backgroung worker history builder
  - Invokes server price generator to get values operate
  - Persists price bases on 1 minute window
  - Process sequencial time windows using first currencies entry
  - Persists all data into Redis database.
  - location: /zeitarbeiter

### RESTApi server price history provider
  - Answer history price for requested currency, according of time range
  - Reads prices data history from Redis database
  - location: /preisviewer

  
### Execution
```sh
$ cd eskalationszeit/
$ ./start.sh
```

### Development
Want to contribute? Great!


### Todos
 - Write MORE Tests

License
----

MIT


**This project is in full development and many things can change!**
