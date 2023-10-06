package main

import (
	"fmt"
	"os"
)

func main() {
	// make sure we're using shred correctly
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run shred.go <file-path>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	err := Shred(filePath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func Shred(filePath string) error {
	// get some info on the file and check that it exists
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", filePath)
	} else if err != nil {
		return fmt.Errorf("error getting file info: %v", err)
	}

	fileSize := fileInfo.Size()

	// open the file in read-write mode and handle errors
	file, err := os.OpenFile(filePath, os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close() // Ensure file is closed after operation

	return nil
}
