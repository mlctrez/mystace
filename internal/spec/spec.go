package spec

import (
	"encoding/json"
	"os"
)

type File struct {
	Overview string `json:"overview"`
	Tests    []Test `json:"tests"`
}

type Test struct {
	Name        string            `json:"name"`
	Description string            `json:"desc"`
	Data        interface{}       `json:"data"`
	Partials    map[string]string `json:"partials"`
	Template    string            `json:"template"`
	Expected    string            `json:"expected"`
}

func Read(path string) (f File, err error) {
	var file *os.File
	file, err = os.Open(path)
	if err != nil {
		return
	}
	err = json.NewDecoder(file).Decode(&f)
	return
}
