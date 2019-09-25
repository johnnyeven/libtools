package transport_grpc

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johnnyeven/libtools/courier"
	"github.com/johnnyeven/libtools/courier/transport_grpc/testify"
)

type Op struct {
	ID   string `json:"id" in:"path"`
	Ping string `json:"ping" in:"query"`
	courier.EmptyOperator
}

func TestCreateStreamHandler(t *testing.T) {
	tt := assert.New(t)

	stream := testify.NewStreamMock(&MsgPackCodec{})

	stream.SendMsg(Op{
		ID:   "id",
		Ping: "ping",
	})

	op := Op{}

	err := MarshalOperator(stream, &op)
	tt.Nil(err)
	tt.Equal(op.Ping, "ping")
	tt.Equal(op.ID, "id")
}
