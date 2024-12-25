package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Start of golang-fault-tolerant")
	_, err := readPossibleNonExistingFile()
	if err != nil {
		fmt.Println("handled fault, error reading file:", err)
	}

	channel := make(chan []byte)
	go recoverableAsyncReadPossibleNonExistingFile(channel)
	<-channel

	fmt.Println("End of golang-fault-tolerant")
}

func readPossibleNonExistingFile() ([]byte, error) {
	file, err := os.Open("nonexistentfile.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// read file
	data := make([]byte, 100)
	_, err = file.Read(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// asyncReadPossibleNonExistingFile assume the function is definedin 3rd party module that panics
func asyncReadPossibleNonExistingFile(channel chan []byte) {
	data, err := readPossibleNonExistingFile()
	if err != nil {
		panic(err)
	}
	channel <- data
	close(channel)
}

func recoverableAsyncReadPossibleNonExistingFile(channel chan []byte) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("handled fault, panic occurred:", r)
		}
		close(channel)
	}()
	asyncReadPossibleNonExistingFile(channel)
}
