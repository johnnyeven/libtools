package bank_gateway

import (
	"fmt"
	"testing"
)

func TestReceptReserve(t *testing.T) {
	type arg struct {
		bankSendTransID string
		want            string
	}

	testCase := make([]arg, 0)

	for i := 0; i < 100; i++ {
		bankSendTransID := fmt.Sprintf("%d", i)
		testCase = append(testCase, arg{
			bankSendTransID: bankSendTransID,
			want:            RECEPT_RESERVE_PRE + bankSendTransID,
		})
	}
	for _, tc := range testCase {
		outPut := GenerateReceptReserve(tc.bankSendTransID)
		if outPut != tc.want {
			t.Errorf("TestGenerateReceptReserve err")
		}

		m, _ := GetBankSendTransIDFromReceptReserve(outPut)
		if m != tc.bankSendTransID {
			t.Errorf("GetBankSendTransIDFromReceptReserve err")
		}
	}
}
