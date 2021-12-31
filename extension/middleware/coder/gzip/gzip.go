package gzip

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"

	"github.com/wubbalubbaaa/easyRpc"
	"github.com/wubbalubbaaa/easyRpc/extension/middleware/coder"
)

func gzipCompress(data []byte) []byte {
	var in bytes.Buffer
	w := gzip.NewWriter(&in)
	w.Write(data)
	w.Close()
	return in.Bytes()
}

func gzipUnCompress(data []byte) ([]byte, error) {
	b := bytes.NewReader(data)
	r, err := gzip.NewReader(b)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	undatas, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return undatas, nil
}

// Gzip represents a gzip coding middleware.
type Gzip int

// Encode implements easyRpc MessageCoder.
func (g *Gzip) Encode(client *easyRpc.Client, msg *easyRpc.Message) *easyRpc.Message {
	if len(msg.Buffer) > int(*g) && !msg.IsFlagBitSet(coder.FlagBitGZip) {
		buf := gzipCompress(msg.Buffer[easyRpc.HeaderIndexReserved+1:])
		total := len(buf) + easyRpc.HeaderIndexReserved + 1
		if total < len(msg.Buffer) {
			copy(msg.Buffer[easyRpc.HeaderIndexReserved+1:], buf)
			msg.Buffer = msg.Buffer[:total]
			msg.SetBodyLen(total - 16)
			msg.SetFlagBit(coder.FlagBitGZip, true)
		}
	}
	return msg
}

// Decode implements easyRpc MessageCoder.
func (g *Gzip) Decode(client *easyRpc.Client, msg *easyRpc.Message) *easyRpc.Message {
	if msg.IsFlagBitSet(coder.FlagBitGZip) {
		buf, err := gzipUnCompress(msg.Buffer[easyRpc.HeaderIndexReserved+1:])
		if err == nil {
			msg.Buffer = append(msg.Buffer[:easyRpc.HeaderIndexReserved+1], buf...)
			msg.SetFlagBit(coder.FlagBitGZip, false)
			msg.SetBodyLen(len(msg.Buffer) - 16)
		}
	}
	return msg
}

// New returns the gzip coding middleware.
func New(n int) *Gzip {
	var g = Gzip(n)
	return &g
}
