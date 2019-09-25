package transform

import (
	"encoding/json"
	"encoding/xml"

	"github.com/johnnyeven/libtools/courier/httpx"
)

func init() {
	RegisterContentTransformer(&ContentTransformer{
		Key:         "json",
		ContentType: httpx.MIMEJSON,
		Marshal:     json.Marshal,
		Unmarshal:   json.Unmarshal,
	})
	RegisterContentTransformer(&ContentTransformer{
		Key:         "xml",
		ContentType: httpx.MIMEXML,
		Marshal:     xml.Marshal,
		Unmarshal:   xml.Unmarshal,
	})
}

var contentTransformers = map[string]*ContentTransformer{}

func RegisterContentTransformer(contentTransformer *ContentTransformer) {
	contentTransformers[contentTransformer.Key] = contentTransformer
	contentTransformers[contentTransformer.ContentType] = contentTransformer
}

func GetContentTransformer(keyOrContentType string) *ContentTransformer {
	return contentTransformers[keyOrContentType]
}

type ContentTransformer struct {
	Key         string
	ContentType string
	Marshal     func(v interface{}) ([]byte, error)
	Unmarshal   func(data []byte, v interface{}) error
}
