package status_error

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestStatusErrorCodeMap_String(t *testing.T) {
	spew.Dump(StatusErrorCodes.String())
}
