package bank_gateway

import (
	"fmt"
	"strings"
)

const (
	RECEPT_RESERVE_PRE = "交易市场流水号:"
)

func GenerateReceptReserve(bankSendTransID string) string {
	return fmt.Sprintf("%s%s", RECEPT_RESERVE_PRE, bankSendTransID)
}

func GetBankSendTransIDFromReceptReserve(reseve string) (bankSendTransID string, err error) {
	a := strings.Split(reseve, RECEPT_RESERVE_PRE)
	if len(a) != 2 {
		err = fmt.Errorf("GetBankSendTransIDFromReceptReserve err reseve:[%s]", reseve)
		return
	}
	bankSendTransID = a[1]
	return
}
