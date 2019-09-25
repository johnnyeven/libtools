package httplib

import (
	"fmt"
	"strconv"
	"strings"

	"encoding/json"
)

type Uint64List []uint64

func (list Uint64List) MarshalJSON() ([]byte, error) {
	if len(list) == 0 {
		return []byte(`[]`), nil
	}
	strValues := make([]string, 0)
	for _, v := range list {
		strValues = append(strValues, fmt.Sprintf(`"%d"`, v))
	}
	return []byte(`[` + strings.Join(strValues, ",") + `]`), nil
}

func (list *Uint64List) UnmarshalJSON(data []byte) (err error) {
	strValues := make([]string, 0)
	err = json.Unmarshal(data, &strValues)
	if err != nil {
		return err
	}
	finalList := Uint64List{}
	for i, strValue := range strValues {
		v, parseErr := strconv.ParseUint(strValue, 10, 64)
		if parseErr != nil {
			err = fmt.Errorf(`[%d] cannot unmarshal string into value of type uint64`, i)
			return
		}
		finalList = append(finalList, v)
	}
	*list = finalList
	return
}
