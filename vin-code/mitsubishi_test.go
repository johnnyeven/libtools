package vincode

import (
	"testing"

	"github.com/johnnyeven/libtools/vin-code/mfrs/msbs"
)

func TestMSBS1(t *testing.T) {
	vin := "JA32A3HU6HH123456"
	mb := MSBSVINCode(vin)
	wmi, err := mb.ParseWMI()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("wmi:%+v", wmi)

	vds, err := mb.ParseVDS()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("vds:%+v", vds)

	vis, err := mb.ParseVIS()
	if err != nil {
		t.Error(err)
		re := msbs.GetVISRune(vin)
		t.Errorf("%+v", re)
		return
	}
	t.Logf("vid:%+v", vis)
}
