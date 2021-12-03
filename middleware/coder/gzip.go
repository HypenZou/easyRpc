package coder

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"

	"github.com/wubbalubbaaa/easyRpc"
)

const GZipFlagBit = 0

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

type Gzip struct {
	critical int
	flagMask byte
}

func (c *Gzip) Encode(client *easyRpc.Client, msg *easyRpc.Message) *easyRpc.Message {
	if len(msg.Buffer) > c.critical && !msg.IsFlagBitSet(GZipFlagBit) {
		buf := gzipCompress(msg.Buffer[easyRpc.HeaderIndexReserved+1:])
		total := len(buf) + easyRpc.HeaderIndexReserved + 1
		if total < len(msg.Buffer) {
			copy(msg.Buffer[easyRpc.HeaderIndexReserved+1:], buf)
			msg.Buffer[easyRpc.HeaderIndexReserved] |= c.flagMask
			msg.Buffer = msg.Buffer[:total]
			msg.SetBodyLen(total - 16)
			msg.SetFlagBit(GZipFlagBit, true)
		}
	}
	return msg
}

func (c *Gzip) Decode(client *easyRpc.Client, msg *easyRpc.Message) *easyRpc.Message {
	if msg.IsFlagBitSet(GZipFlagBit) {
		buf, err := gzipUnCompress(msg.Buffer[easyRpc.HeaderIndexReserved+1:])
		if err == nil {
			msg.Buffer = append(msg.Buffer[:easyRpc.HeaderIndexReserved+1], buf...)
			msg.SetFlagBit(GZipFlagBit, false)
		}
	}
	return msg
}

func NewGzip() *Gzip {
	return &Gzip{critical: 256, flagMask: 0x1}
}
