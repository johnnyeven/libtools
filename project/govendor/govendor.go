package govendor

import (
	"encoding/json"
	"io/ioutil"
)

func LoadGoVendorJSON(p string) (Vendor, error) {
	jsonBytes, err := ioutil.ReadFile(p)
	vendor := Vendor{}
	if err != nil {
		return vendor, err
	}
	json.Unmarshal(jsonBytes, &vendor)
	return vendor, nil
}

type VendorPackage struct {
	Path         string `json:"path"`
	Revision     string `json:"revision"`
	RevisionTime string `json:"revisionTime"`
	ChecksumSHA1 string `json:"checksumSHA1"`
}

type Vendor struct {
	Comment  string          `json:"comment"`
	Ignore   string          `json:"ignore"`
	RootPath string          `json:"rootPath"`
	Package  []VendorPackage `json:"package"`
}

func (v *Vendor) ListImportPath() []string {
	slice := make([]string, 0)
	for _, p := range v.Package {
		path := p.Path
		slice = append(slice, path)
	}
	return slice
}
