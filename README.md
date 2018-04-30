# naivechain

a simple Blockchain inspired by https://github.com/kofj/naivechain

## First Node

```bash
go run main.go 
```

## Second Node

```bash
go run main.go -api :3002 -p2p :6002 -peers ws://localhost:6001
```

```bash
$ curl http://localhost:3001/blocks
```

## mine block

```bash
$ curl -H "Content-type:application/json" --data '{"data" : "Some data to the first block"}' http://localhost:3001/mine_block
```

## add peer

```bash
$ curl -H "Content-type:application/json" --data '{"peer" : "ws://localhost:6002"}' http://localhost:3001/add_peer
```

## query peers

```bash
$ curl http://localhost:3001/peers
```
