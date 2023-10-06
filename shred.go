package main

func main() {
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
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", filePath)
	} else if err != nil {
		return fmt.Errorf("error getting file info: %v", err)
	}