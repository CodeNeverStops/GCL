package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// var countResult = struct {
// 	sync.RWMutex
// 	m map[string]int
// }{m: make(map[string]int)}

var Version = "0.0.1"

var countMap = sync.Map{}
var wg = sync.WaitGroup{}

var (
	help        = flag.Bool("help", false, "show usage help and quit")
	showVersion = flag.Bool("version", false, "show version and quit")
	fileType    = flag.String("filetype", "", "specify the file type to count line")
	topNum      = flag.Int("top", 0, "list top N files")
	sortAsc     = flag.Bool("sortasc", false, "sort files by line count asc")
)

func main() {
	flag.Usage = usage
	flag.Parse()
	if *help {
		usage()
		return
	}
	if *showVersion {
		fmt.Println(Version)
		return
	}
	dirList := flag.Args()

	taskQueue := make(chan string, 10000)
	go startWorker(taskQueue)

	var fileTypeList []string
	if len(*fileType) > 0 {
		fileTypeList = strings.Split(*fileType, "|")
	}

	if len(dirList) == 0 {
		dirList = []string{"."}
	}
	for _, dir := range dirList {
		readDir(dir, taskQueue, fileTypeList)
	}

	wg.Wait()

	total := 0
	type kv struct {
		Key   string
		Value int
	}
	var countResult []kv
	countMap.Range(func(k, v interface{}) bool {
		total += v.(int)
		countResult = append(countResult, kv{k.(string), v.(int)})
		// fmt.Printf("%v:%v\n", k, v)
		return true
	})
	fmt.Printf("%s: %d\n", "total", total)
	if *topNum == 0 {
		return
	}
	if *sortAsc {
		sort.Slice(countResult, func(i, j int) bool {
			return countResult[i].Value < countResult[j].Value
		})
	} else {
		sort.Slice(countResult, func(i, j int) bool {
			return countResult[i].Value > countResult[j].Value
		})
	}
	fmt.Println("count result")
	for i, kv := range countResult {
		if i >= *topNum {
			break
		}
		fmt.Printf("%s: %d[%d%%]\n", kv.Key, kv.Value, kv.Value*100/total)
	}

	// countResult.RLock()
	// for k, v := range countResult.m {
	// 	fmt.Printf("%v:%v\n", k, v)
	// }
	// countResult.RUnlock()
}

func usage() {
	fmt.Println("usage: gcl [flags] [dir]")
	fmt.Println("flags:")
	flag.PrintDefaults()
}

func readDir(dir string, taskQueue chan<- string, fileTypeList []string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		filePath := filepath.Join(dir, file.Name())
		if file.IsDir() {
			readDir(filePath, taskQueue, fileTypeList)
		} else {
			typeIsMatch := false
			if len(fileTypeList) > 0 {
				for _, fileType := range fileTypeList {
					if filepath.Ext(filePath) == fileType {
						typeIsMatch = true
						break
					}
				}
			} else {
				typeIsMatch = true
			}
			if !typeIsMatch {
				continue
			}
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
			lineStr := fmt.Sprintf("%s", part)
			lineStr = strings.Trim(lineStr, " ")
			if len(lineStr) == 0 {
				continue
			}
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
