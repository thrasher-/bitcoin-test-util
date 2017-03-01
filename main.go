package main

import (
	"log"
	"encoding/json"
	"io/ioutil"
	"os/exec"
	"strings"
)

type BitcoinTests []struct {
	Exec string `json:"exec"`
	Args []string `json:"args"`
	OutputCmp string `json:"output_cmp,omitempty"`
	Description string `json:"description"`
	Input string `json:"input,omitempty"`
	ReturnCode int `json:"return_code,omitempty"`
}

func WriteFile(path string, data []byte) error {
	if strings.Contains(path, ".hex") {
		strings.Replace(string(data), "\n","",-1)
	}
	err := ioutil.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func OpenFile(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if strings.Contains(path, ".hex") {
		strings.Replace(string(data), "\n","",-1)
	}
	return data, nil
}

func ExecuteLitecoinTX(args []string, exePath, testPath, input, output string) error {
	var cmdOut []byte
	var err error

	cmd := exec.Command(exePath, args...)

	if input != "" {
		data, err := OpenFile(testPath + input)
		if err != nil {
			return err
		}
		cmd.Stdin = strings.NewReader(string(data))
	}

	cmdOut, err = cmd.Output()
	if err != nil {
		return err
	}

	if output != "" {
		outputCmp, err := OpenFile(testPath + output)
		if err != nil {
			return err
		}
		//log.Printf("Output file contents %s", outputCmp)

		if string(cmdOut) == string(outputCmp) {
			log.Println("Input matches expected output")
		} else {
			log.Println("Updating output cmp file")
			err = WriteFile(testPath + output, cmdOut)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func main() {
	testPath := "/home/thrasher/Desktop/dev/litecoin/src/test/data/"
	jsonFile := testPath + "bitcoin-util-test.json"
	exePath := "/home/thrasher/Desktop/dev/litecoin/src/litecoin-tx"

	data, err := OpenFile(jsonFile)
	if err != nil {
		log.Fatal("Failed to read file. Err: %s", err)
	}

	tests := BitcoinTests{}
	json.Unmarshal(data, &tests)

	for _, x := range tests {
		err := ExecuteLitecoinTX(x.Args, exePath, testPath, x.Input, x.OutputCmp)
		if err != nil {
			if strings.Contains(x.Description, "Expected to fail") {
				continue
			}
			log.Printf("Failure on test: %s. Error: %s", x.Description, err)
			log.Printf("Args %s\n", strings.Join(x.Args, " "))
		}
	}
}


