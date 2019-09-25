package transform

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTagValueAndFlags(t *testing.T) {
	tt := assert.New(t)

	cases := []struct {
		tag   string
		v     string
		flags TagFlags
	}{
		{
			tag:   "json",
			v:     "json",
			flags: TagFlags{},
		},
		{
			tag: "json,omitempty",
			v:   "json",
			flags: TagFlags{
				"omitempty": true,
			},
		},
	}

	for _, tc := range cases {
		v, flags := GetTagValueAndFlags(tc.tag)
		tt.Equal(tc.v, v)
		tt.Equal(tc.flags, flags)
	}
}

type TestDataForLocateJSONPath struct {
	Data struct {
		Bool        bool `json:"bool"`
		StructSlice []struct {
			Name string `json:"name"`
		} `json:"structSlice"`
		StringSlice []string `json:"stringSlice"`
		NestedSlice []struct {
			Names []string `json:"names"`
		} `json:"nestedSlice"`
	} `json:"data"`
}

func TestLocateJSONPath(t *testing.T) {
	tt := assert.New(t)

	data := TestDataForLocateJSONPath{}

	{
		d := []byte(`
{
 	"data": {
		"bool":   ""
	}
}
`)
		typeError := json.Unmarshal(d, &data).(*json.UnmarshalTypeError)
		tt.Equal("data.bool", LocateJSONPath(d, typeError.Offset))
	}

	{
		d := []byte(`
{
		"data": {
			"structSlice": [
				{"name":"{"},
				{"name":"1"},
				{"name": { "test": 1 }},
				{"name":"1"}
			]
		}
}
`)
		typeError := json.Unmarshal(d, &data).(*json.UnmarshalTypeError)
		tt.Equal("data.structSlice[2].name", LocateJSONPath(d, typeError.Offset))
	}

	{
		d := []byte(`
	{
		"data": {
			"stringSlice":["1","2",3]
		}
	}
	`)
		typeError := json.Unmarshal(d, &data).(*json.UnmarshalTypeError)
		tt.Equal("data.stringSlice[2]", LocateJSONPath(d, typeError.Offset))
	}

	{
		d := []byte(`
	{
		"data": {
			"bool": true,
			"nestedSlice": [
				{ "names": ["1","2","3"] },
	            { "names": ["1","\"2", 3] }
			]
		}
	}
	`)
		typeError := json.Unmarshal(d, &data).(*json.UnmarshalTypeError)
		tt.Equal("data.nestedSlice[1].names[2]", LocateJSONPath(d, typeError.Offset))
	}
}
