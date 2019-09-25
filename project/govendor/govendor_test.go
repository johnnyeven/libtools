package govendor

import (
	"fmt"
	"testing"
)

func TestLoadGoVendorJSON(t *testing.T) {
	vendor, _ := LoadGoVendorJSON("../../vendor/vendor.json")
	fmt.Printf("%#v", vendor.ListImportPath())
}
