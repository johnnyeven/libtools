package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"

	"github.com/johnnyeven/libtools/format"
)

func LoadFiles(dir string, filter func(filename string) bool) (filenames []string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		filename := path.Join(dir, file.Name())
		if file.IsDir() {
			filenames = append(filenames, LoadFiles(filename, filter)...)
		} else {
			if filter(filename) {
				filenames = append(filenames, filename)
			}
		}
	}
	return
}

var cmdFormat = &cobra.Command{
	Use:   "fmt",
	Short: "format",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, _ := os.Getwd()
		files := LoadFiles(cwd, func(filename string) bool {
			return path.Ext(filename) == ".go" && !strings.Contains(filename, "/vendor/")
		})

		for _, filename := range files {
			fileInfo, _ := os.Stat(filename)
			bytes, _ := ioutil.ReadFile(filename)
			nextBytes, err := format.Process(filename, bytes)
			if err != nil {
				panic(fmt.Errorf("errors %s in %s", filename, err.Error()))
			}
			if string(nextBytes) != string(bytes) {
				fmt.Printf("reformatted %s\n", filename)
				err := ioutil.WriteFile(filename, nextBytes, fileInfo.Mode())
				if err != nil {
					panic(err)
				}
			}
		}
	},
}

func init() {
	cmdRoot.AddCommand(cmdFormat)
}
