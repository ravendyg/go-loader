### Description
Console multi threaded(?) file loder supported resuming downloading.

Cannot dynamically reallocate threads (maybe implement later).

### Compile
`go build main.go`

### Run
`./main -url <url>`

### Debugging
`go run main.go -url http://localhost:3011`


### TODO
- Make possible to continue an interrupted download
