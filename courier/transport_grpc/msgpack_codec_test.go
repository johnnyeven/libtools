package transport_grpc

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

type Body struct {
}

type WithID struct {
	ID string `json:"id" in:"path"`
}

type Item2 struct {
	WithID
}

type Item struct {
	WithID
	Ping string `json:"ping" in:"query"`
	Path string
}

type ItemReq struct {
	ID   string `json:"id" in:"path"`
	Ping []byte `json:"ping" in:"query"`
}

func TestMsg(t *testing.T) {
	tt := assert.New(t)

	codec := MsgPackCodec{}
	item := ItemReq{}
	item.ID = "13123.123"
	item.Ping = []byte("ping")

	bytes, err := codec.Marshal(item)
	spew.Dump(bytes)
	tt.Nil(err)

	{
		withID := WithID{}
		err := codec.Unmarshal(bytes, &withID)
		tt.Nil(err)
		spew.Dump(withID)
	}

	{
		item := &Item{}
		err := codec.Unmarshal(bytes, &item)
		tt.Nil(err)
		spew.Dump(item)
	}
}
