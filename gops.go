package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"unicode/utf8"

	"github.com/paddie/gokmp"
)

var END_OF_WORK = ""
var MAX_OPEN_FILES uint64 = 17711
var WORK_QUEUE_SIZE int = 28657
var inv_ext map[string]bool
var lookingfor string
var kmp *gokmp.KMP
var max_workers int = 228
var wg sync.WaitGroup
var work_queue chan string
var done chan bool
var index int

func main() {
	flag.Parse()
	// TODO: handle incorrect input

	getReady(flag.Arg(0))
	crawlFolder()

	// tell workers it's over
	work_queue <- END_OF_WORK
	wg.Wait()

	fmt.Println("Done.")
}

func getReady(lookingFor string) {
	// put cpu cores to work
	numCPUs := runtime.NumCPU() - 1
	if numCPUs < 1 {
		numCPUs = 1
	}
	runtime.GOMAXPROCS(numCPUs)

	/*
		This doesn't even work in Windows... TODO: fix properly
		// Set max open files
		var rLimit syscall.Rlimit
		err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
		if err != nil {
			fmt.Println("Error Getting Rlimit ", err)
		}
		rLimit.Cur = MAX_OPEN_FILES
		err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
		if err != nil {
			fmt.Println("Error Setting Rlimit ", err)
		} else {
			max_workers = MAX_OPEN_FILES / 3
		}
	*/

	// Setup workers
	work_queue = make(chan string, WORK_QUEUE_SIZE)
	done = make(chan bool)

	for i := 0; i < max_workers/3; i++ {
		go searchWorker(work_queue)
	}
	go func() {
		for i := 0; i < (max_workers * (2 / 3)); i++ {
			go searchWorker(work_queue)
		}
	}()

	// Get KMP ready
	lookingfor = flag.Arg(0)
	fmt.Println("Searching for ->", lookingfor)
	kmp, _ = gokmp.NewKMP(lookingfor)
}

func crawlFolder() {
	current := "."
	filepath.Walk(current, visit)
}

func visit(path string, fileInfo os.FileInfo, err error) error {

	if fileInfo.IsDir() || err != nil {
		return nil
	}
	work_queue <- path
	return nil
}

func searchWorker(work_queue chan string) {
	wg.Add(1)
	defer wg.Done()
	localkmp := kmp
	var x []byte
	var err error
	for path := range work_queue {
		if path == END_OF_WORK {
			// tell everyone else we're done
			work_queue <- END_OF_WORK
			break
		}

		x, err = ioutil.ReadFile(path)
		if err != nil || x == nil {
			x = nil
			continue //fail gracefully
		}

		if index > len(x) {
			index = len(x)
		}
		if !utf8.ValidString(string(x[:index])) {
			continue
		}

		if localkmp == nil {
			localkmp = kmp
		}
		if localkmp.ContainedIn(string(x)) {
			fmt.Println(path)
		}
	}
}
