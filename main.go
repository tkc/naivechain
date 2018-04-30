package main

import (
	"./app"
	"flag"
	"log"
	"net/http"
	"strings"
	"golang.org/x/net/websocket"
)

var (
	httpAddr     = flag.String("api", ":3001", "api server address.")
	p2pAddr      = flag.String("p2p", ":6001", "p2p server address.")
	initialPeers = flag.String("peers", "ws://localhost:6001", "initial peers")
)

func errFatal(msg string, err error) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

func main() {

	flag.Parse()
	app.ConnectToPeers(strings.Split(*initialPeers, ","))

	http.HandleFunc("/blocks", app.HandleBlocks)
	http.HandleFunc("/mine_block", app.HandleMineBlock)
	http.HandleFunc("/peers", app.HandlePeers)
	http.HandleFunc("/add_peer", app.HandleAddPeer)

	go func() {
		log.Println("Listen HTTP on", *httpAddr)
		errFatal("start api server", http.ListenAndServe(*httpAddr, nil))
	}()

	http.Handle("/", websocket.Handler(app.WsHandleP2P))
	log.Println("Listen P2P on ", *p2pAddr)
	errFatal("start p2p server", http.ListenAndServe(*p2pAddr, nil))
}
