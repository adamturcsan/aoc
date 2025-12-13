package main

import (
	"bufio"
	"fmt"
	"os"
)

func readFileLineByLine(filePath string, lineCh chan string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error while opening file")
	}
	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		lineCh <- string(line)
	}
	close(lineCh)
}

func broadcastChannel[T any](channel chan T, numberOfLines int) []chan T {
	lines := make([]chan T, numberOfLines)
	for i := 0; i < numberOfLines; i++ {
		lines[i] = make(chan T)
	}
	go func() {
		for data := range channel {
			for _, line := range lines {
				line <- data
			}
		}
		for _, line := range lines {
			close(line)
		}
	}()
	return lines
}
