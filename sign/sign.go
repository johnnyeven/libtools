package sign

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"github.com/johnnyeven/libtools/courier/status_error"
	"github.com/johnnyeven/libtools/courier/transport_http/transform"
)

var Sign = "sign"
var RandString = "randString"
var AccessKey = "AccessKey"

type SecretExchanger func(key string) (string, error)

type SignParams struct {
	// AccessKey
	AccessKey string `json:"AccessKey" in:"header"`
	// 签名
	Sign string `json:"sign" validate:"@string[1,32]" in:"query"`
	// 随机字符串
	RandString string `json:"randString" validate:"@string[1,32]" in:"query"`
}

func getSign(req *http.Request, query url.Values, secretExchanger SecretExchanger) (sign []byte, origin []byte, err error) {
	accessKey := req.Header.Get(AccessKey)
	randString := query.Get(RandString)

	if accessKey != "" && randString != "" {
		secret, errForExchange := secretExchanger(accessKey)
		if errForExchange != nil {
			err = status_error.InvalidSecret.StatusError().WithDesc(errForExchange.Error())
			return
		}
		bodyBytes := make([]byte, 0)
		if req.Body != nil {
			bodyBytes, err = transform.CloneRequestBody(req)
			if err != nil {
				return
			}
		}
		sign, origin = Secret(secret).Encode(query, bodyBytes)
	}
	return
}

type Secret string

func (secret Secret) Encode(query url.Values, body []byte) (sign []byte, origin []byte) {
	keyList := make([]string, 0)
	for key := range query {
		keyList = append(keyList, key)
	}
	sort.Strings(keyList)

	rawSignStr := &bytes.Buffer{}
	for _, key := range keyList {
		values := query[key]
		if len(values) == 0 || key == Sign {
			continue
		}
		for _, v := range values {
			rawSignStr.WriteString(key)
			rawSignStr.WriteString("=")
			rawSignStr.WriteString(v)
			rawSignStr.WriteString("&")
		}
	}

	if len(body) > 0 {
		rawSignStr.WriteString("body")
		rawSignStr.WriteString("=")
		rawSignStr.Write(genMd5(body))
		rawSignStr.WriteString("&")
	}

	rawSignStr.WriteString("secret")
	rawSignStr.WriteString("=")
	rawSignStr.WriteString(string(secret))

	origin = rawSignStr.Bytes()
	sign = genMd5(origin)
	return
}

func genMd5(src []byte) (dst []byte) {
	hasher := md5.New()
	hasher.Write(src)
	sum := hasher.Sum(nil)

	dst = make([]byte, hex.EncodedLen(len(sum)))
	hex.Encode(dst, sum)
	return
}

func GenSign(secret string, queryData map[string]interface{}, body []byte) string {
	keyList := []string{}
	for key, _ := range queryData {
		keyList = append(keyList, key)
	}
	sort.Strings(keyList)
	var rawSignStr string
	for _, key := range keyList {
		rawSignStr += fmt.Sprintf("%v=%v%s", key, queryData[key], "&")
	}
	if len(body) > 0 {
		rawSignStr += fmt.Sprintf("body=%s%s", string(genMd5(body)), "&")
	}
	rawSignStr += fmt.Sprintf("secret=%s", secret)

	return string(genMd5([]byte(rawSignStr)))
}
