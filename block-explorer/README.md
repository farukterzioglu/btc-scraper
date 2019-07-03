`$ go get -u github.com/go-swagger/go-swagger/cmd/swagger`  
`$ swagger generate spec -o ./swaggerui/swagger.json --scan-models`  
`go run .`  

Navigate to http://localhost:8000/swaggerui/#/  

`curl -X GET "http://localhost:8000/v1/btc/block" -H "accept: application/json"`  
`curl -X GET "http://localhost:8000/v1/btc/block/5" -H "accept: application/json" `  