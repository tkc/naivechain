package app

import (
	"io"
	"log"
	"sort"
	"encoding/json"
	"golang.org/x/net/websocket"
)

const (
	queryLatest = iota
	queryAll = iota
	responseBlockChain = iota
)

var (
	sockets  []*websocket.Conn
	blockchain = []*Block{genesisBlock}
)

type ResponseBlockchain struct {
	Type int    `json:"type"`
	Data string `json:"data"`
}

func ConnectToPeers(peersAddr []string) {
	for _, peer := range peersAddr {
		if peer == "" {
			continue
		}
		ws, err := websocket.Dial(peer, "", peer)
		if err != nil {
			continue
		}
		initConnection(ws)
	}
}

func initConnection(ws *websocket.Conn) {
	go WsHandleP2P(ws)
	ws.Write(queryLatestMsg())
}

func WsHandleP2P(ws *websocket.Conn) {

	var (
		v    = &ResponseBlockchain{}
		peer = ws.LocalAddr().String()
	)

	sockets = append(sockets, ws)

	for {

		var msg []byte
		err := websocket.Message.Receive(ws, &msg)

		if err == io.EOF {
			break
		}

		if err != nil {
			break
		}

		log.Printf("Received[from %s]: %s.\n", peer, msg)
		err = json.Unmarshal(msg, v)
		log.Fatalln("invalid p2p msg", err)

		switch v.Type {
		case queryLatest:

			bs := responseLatestMsg()
			log.Printf("responseLatestMsg: %s\n", bs)
			ws.Write(bs)

		case queryAll:

			d, _ := json.Marshal(blockchain)
			v.Type = responseBlockChain
			v.Data = string(d)
			bs, _ := json.Marshal(v)
			log.Printf("responseChainMsg: %s\n", bs)
			ws.Write(bs)

		case responseBlockChain:
			handleBlockChainResponse([]byte(v.Data))
		}
	}
}

func broadcast(msg []byte) {
	for n, socket := range sockets {
		_, err := socket.Write(msg)
		if err != nil {
			log.Printf("peer [%s] disconnected.", socket.RemoteAddr().String())
			sockets = append(sockets[0:n], sockets[n+1:]...)
		}
	}
}

func handleBlockChainResponse(msg []byte) {
	var receivedBlocks = []*Block{}

	err := json.Unmarshal(msg, &receivedBlocks)
	log.Fatalln("invalid blockchain", err)
	sort.Sort(ByIndex(receivedBlocks))

	latestBlockReceived := receivedBlocks[len(receivedBlocks)-1]
	latestBlockHeld := getLatestBlock()

	if latestBlockReceived.Index > latestBlockHeld.Index {
		log.Printf("blockchain possibly behind. We got: %d Peer got: %d", latestBlockHeld.Index, latestBlockReceived.Index)
		if latestBlockHeld.Hash == latestBlockReceived.PreviousHash {
			blockchain = append(blockchain, latestBlockReceived)
		} else if len(receivedBlocks) == 1 {
			broadcast(queryAllMsg())
		} else {
			replaceChain(receivedBlocks)
		}

	} else {
		log.Println("received blockchain is not longer than current blockchain. Do nothing.")
	}

}
