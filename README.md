Related article (tr) : [Bitcoin İle Konuşan Uygulamalar](https://link.medium.com/pUHIIt0poY)  
  
Follow the instruction from : [Bitcoin'i Kodlamak (Coding Bitcoin, Turkish)](https://medium.com/@farukterzioglu/bitcoini-kodlamak-1-golang-ile-5e7833c0dc19) (or from [setup_notes.md](setup_notes.md))
  

Copy rpc.cert  
`sudo cp ~/.btcd/rpc.cert /`
  
Send rpc request 
`curl --cacert /rpc.cert --user myuser:SomeDecentp4ssw0rd --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "getinfo", "params": [] }' -H 'content-type: text/plain;' https://127.0.0.1:18556/`
  
Result;  
`{"result":{"version":120000,"protocolversion":70002,"blocks":1,"timeoffset":0,"connections":0,"proxy":"","difficulty":1,"testnet":false,"relayfee":0.00001,"errors":""},"error":null,"id":"curltest"}`

(For other requests, check -> [rpc.http](rpc.http))

### Resources
https://en.bitcoin.it/wiki/API_reference_(JSON-RPC)  
https://bitcoin.org/en/developer-reference#remote-procedure-calls-rpcs  
https://github.com/btcsuite/btcd/blob/master/docs/json_rpc_api.md#ExampleGetBlockCount
