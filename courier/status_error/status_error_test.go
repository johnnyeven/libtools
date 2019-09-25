package status_error

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseString(t *testing.T) {
	tt := assert.New(t)
	errString := ReadFailed.StatusError().String()
	statusError := ParseString(errString)
	tt.Equal(ReadFailed.StatusError(), statusError)
}
