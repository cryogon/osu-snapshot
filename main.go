package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: command [backup|restore]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "backup":
		Backup()
	case "restore":
		Restore()
	default:
		fmt.Println("Usage: command [backup|restore]")
		os.Exit(1)
	}
}
