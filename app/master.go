package main

import (
	"fmt"
	"os"
)

const RDBMasterStateFilePath = "empty.rdb"

func RDBState() string {
	data, err := os.ReadFile(RDBMasterStateFilePath)
	if err != nil {
		panic("RDB file corrupted")
	}
	return formatRDB(string(data))
}

func formatRDB(data string) string {
	formattedStr := "$"
	formattedStr += fmt.Sprintf("%d\r\n", len(data))
	formattedStr += data
	return formattedStr
}
