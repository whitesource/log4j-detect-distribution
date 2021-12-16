package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

// fp: file path of file to Decode
// obj: struct pointer to Decode the file content into
// return: error in case failed to open or Decode the file
func Decode(fp string, obj interface{}) error {
	opened, err := os.Open(fp)
	if err != nil {
		// TODO: change to debug log file
		fmt.Println("Decode - Failed to open file ", fp, err)
		return err
	}

	decoder := json.NewDecoder(opened)
	if err = decoder.Decode(obj); err != nil {
		// TODO: change to debug log file
		fmt.Println("Decode - Failed to Decode file ", fp, err)
		return err
	}

	return nil
}
