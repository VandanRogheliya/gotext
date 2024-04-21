package main

import (
	"io/fs"
	"os"
)

func ReadFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func WriteFile(path string, content string) {
	data := []byte(content)
	/*
		The first digit 0 is for special permissions and is often ignored in basic file permissions.
		The second digit 6 (in binary 110) gives the owner of the file read (4) and write (2) permissions.
		The third digit 4 (in binary 100) gives all other users read (4) permissions.
	*/
	base8UnixFilePermissions := 0644
	err := os.WriteFile(path, data, fs.FileMode(base8UnixFilePermissions))

	if err != nil {
		panic(err)
	}
}
