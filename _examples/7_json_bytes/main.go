package main

import (
	"encoding/json"
	"os"

	"github.com/kataras/tableprinter"
)

func getMyJSONBytes() []byte {
	data := struct {
		// json tags are optionally but if set they are being used for the headers on `PrintJSON`.
		Firstname string `json:"first name"`
		Lastname  string `json:"last name"`
	}{"Georgios", "Callas"}
	b, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		panic(err)
	}

	return b
}

func main() {
	b := getMyJSONBytes()
	tableprinter.PrintJSON(os.Stdout, b)
}
