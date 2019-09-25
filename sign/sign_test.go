package sign

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	tt := assert.New(t)

	secret := Secret("23a1d8b6-1394-4953-9133-bddde9f2a3dd-ad74926f-8013-414e-832e-66b184f3c202")

	body, _ := json.Marshal(map[string]interface{}{
		"amount": 20,
		"extID":  "21232423asfa",
	})

	query := url.Values{}
	query.Add("randString", "123456")

	sign, origin := secret.Encode(query, body)
	t.Log(string(origin))
	tt.Equal("ff134c81864c1fbdc38c82335ba56a26", string(sign))
}
