package utils



import (
	"os"
	"fmt"
)


func CheckIfIsDir(path string) {
	CheckIfFileExists(path)
	file, err := os.Lstat(path)
	if err != nil {
		fmt.Println("path does not exist")
		os.Exit(1)
	}
	if file.Mode().IsDir() == false {
		fmt.Println("path is not a directory")
		os.Exit(1)
	}
}

func CheckIfFileExists(path string) {
	_, err := os.Lstat(path)
	if err != nil {
		fmt.Println("no such directory")
		os.Exit(1)
	}
}
