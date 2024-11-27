# How to run 

After cloning the repository just run the command:
```
git clone https://github.com/biganashvili/tx-parser.git
cd ./tx-parser
go run cmd/api/main.go
```
This entry point exposes HTTP API and runs the blockchain parser simultaneously

## API endpoints

subscribe

```
curl --location --request POST 'localhost:8080/subscribe?address=0x28C6c06298d514Db089934071355E5743bf21d60'
```

currentBlock
```
curl --location 'localhost:8080/currentBlock'
```
transactions
```
curl --location 'localhost:8080/transactions?address=0x28C6c06298d514Db089934071355E5743bf21d60'
```

## Disadvantages
- The parser is slow because it uses one goroutine.
- Tests not implemented
- Doesn't read variables from env
- Logs needs to be improved

