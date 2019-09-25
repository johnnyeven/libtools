package gen

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johnnyeven/libtools/courier"
)

func ToMap(schema interface{}) (result interface{}) {
	bytes, _ := json.Marshal(schema)
	json.Unmarshal(bytes, &result)
	return
}

func TestParseSuccessMetadata(t *testing.T) {
	tt := assert.New(t)
	tt.Equal(ParseSuccessMetadata("@success content-type text/plain"), courier.Metadata{
		"content-type": []string{"text/plain"},
	})
}
