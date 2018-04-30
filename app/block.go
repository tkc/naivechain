package app

import (
	"fmt"
	"log"
	"time"
	"crypto/sha256"
)

var genesisBlock = &Block{
	Index:        0,
	PreviousHash: "0",
	Timestamp:    1465154705,
	Data:         "genesis block",
	Hash:         "816534932c2b7154836da6afc367695e6337db8a921823784c14378abed4f7d7",
}

type Block struct {
	Index        int64  `json:"index"`
	PreviousHash string `json:"previousHash"`
	Timestamp    int64  `json:"timestamp"`
	Data         string `json:"data"`
	Hash         string `json:"hash"`
}

func (b *Block) String() string {
	return fmt.Sprintf("index: %d,previousHash:%s,timestamp:%d,data:%s,hash:%s", b.Index, b.PreviousHash, b.Timestamp, b.Data, b.Hash)
}

type ByIndex []*Block

func (b ByIndex) Len() int           { return len(b) }
func (b ByIndex) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByIndex) Less(i, j int) bool { return b[i].Index < b[j].Index }

func addBlock(b *Block) {
	if isValidNewBlock(b, getLatestBlock()) {
		blockchain = append(blockchain, b)
	}
}

func isValidNewBlock(nb, pb *Block) bool {
	return  nb.Hash == calculateHashForBlock(nb) && pb.Index+1 == nb.Index && pb.Hash == nb.PreviousHash
}

func isValidChain(bc []*Block) bool {
	if bc[0].String() != genesisBlock.String() {
		log.Println("No same GenesisBlock.", bc[0].String())
		return false
	}

	var temp = []*Block{bc[0]}
	for i := 1; i < len(bc); i++ {
		if isValidNewBlock(bc[i], temp[i-1]) {
			temp = append(temp, bc[i])
		} else {
			return false
		}

	}
	return true
}

func replaceChain(bc []*Block) {
	if isValidChain(bc) && len(bc) > len(blockchain) {
		blockchain = bc
		broadcast(responseLatestMsg())
	}
}

func getLatestBlock() (block *Block) { return blockchain[len(blockchain)-1] }

func generateNextBlock(data string) (nb *Block) {
	var previousBlock = getLatestBlock()
	nb = &Block{
		Data:         data,
		PreviousHash: previousBlock.Hash,
		Index:        previousBlock.Index + 1,
		Timestamp:    time.Now().Unix(),
	}
	nb.Hash = calculateHashForBlock(nb)
	return nb
}

func calculateHashForBlock(b *Block) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%d%s%d%s", b.Index, b.PreviousHash, b.Timestamp, b.Data))))
}