package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// var countResult = struct {
// 	sync.RWMutex
// 	m map[string]int
// }{m: make(map[string]int)}

var countMap = sync.Map{}
var wg = sync.WaitGroup{}

func main() {
	taskQueue := make(chan string, 1000)
	go startWorker(taskQueue)
	readDir(".", taskQueue)
	wg.Wait()
	countMap.Range(func(k, v interface{}) bool {
		fmt.Printf("%v:%v\n", k, v)
		return true
	})
	// countResult.RLock()
	// for k, v := range countResult.m {
	// 	fmt.Printf("%v:%v\n", k, v)
	// }
	// countResult.RUnlock()
}

func readDir(dir string, taskQueue chan<- string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		filePath := filepath.Join(dir, file.Name())
		if file.IsDir() {
			readDir(filePath, taskQueue)
		} else {
			wg.Add(1)
			taskQueue <- filePath
		}
	}
}

func startWorker(taskQueue <-chan string) {
	for {
		select {
		case filePath := <-taskQueue:
			go lineCount(filePath)
		}
	}
}

func lineCount(filePath string) {
	defer wg.Done()
	file, err := os.Open(filePath)
	if err != nil {
		countMap.Store(filePath, 0)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 1024))

	var (
		part   []byte
		prefix bool
		lines  int
	)
	lines = 0
	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			lines++
			buffer.Reset()
		}
	}
	if err == io.EOF {
		err = nil
	}
	countMap.Store(filePath, lines)
	// countResult.Lock()
	// countResult.m[filePath] = 1
	// countResult.Unlock()
}
