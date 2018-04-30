package app

import (
	"fmt"
	"encoding/json"
)

func queryLatestMsg() []byte { return []byte(fmt.Sprintf("{\"type\": %d}", queryLatest)) }

func queryAllMsg() []byte { return []byte(fmt.Sprintf("{\"type\": %d}", queryAll)) }

func responseLatestMsg() (bs []byte) {
	var v = &ResponseBlockchain{Type: responseBlockChain}
	d, _ := json.Marshal(blockchain[len(blockchain)-1:])
	v.Data = string(d)
	bs, _ = json.Marshal(v)
	return bs
}
