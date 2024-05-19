

Start the Redis server with 
`redis-server`

Start the Redis CLI with 
`redis-cli`

` cd $ROOT/controller`
`go run .`

```
 curl http://127.0.0.1:8080/account --data '{"name":"Neel","balance":300}' -H "C
ontent-Type: application/json"
 curl http://127.0.0.1:8080/account?user=Neel
```