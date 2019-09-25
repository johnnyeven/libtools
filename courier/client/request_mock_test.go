package client

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestUnMarshalRequestID(t *testing.T) {
	tt := assert.New(t)
	mock := Mock("service-").For("b.some", MockData{
		Data: []byte(`{"a":1}`),
	})
	parsedMock, err := ParseMockID("service-", mock.RequestID())
	tt.Nil(err)
	spew.Dump(parsedMock)
	tt.Equal(mock, parsedMock)
}

func TestUnMarshalRequestID2(t *testing.T) {
	tt := assert.New(t)
	_, err := ParseMockID("service-x-open", "service-x-open:pay.FetchChannelAccountByID:eyJkYXRhIjoiZXlKaVlXeGhibU5sUVcxdmRXNTBJam94TENKamFHRnVibVZzUVdOamIzVnVkRWxFSWpvaU1USXpJaXdpWTJoaGJtNWxiRWxFSWpvaVVFRkNYMWRKVkU1RlUxTWlMQ0pqZFhKeVpXNWplU0k2SWxKTlFpSXNJbWxrSWpvaU1USXpJaXdpYkdGemRFSmhiR0Z1WTJWQmJXOTFiblFpT2pFc0lteGhjM1JDWVd4aGJtTmxRMmhoYm1kbFZHbHRaU0k2SWpJd01UY3RNRGN0TXpCVU1qSTZNRGM2TURRdU5URTRXaUlzSW5OMFlYUjFjeUk2SWtaU1QxcEZUaUo5In0=")
	tt.Nil(err)
}
