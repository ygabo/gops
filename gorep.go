package main

import (
	"flag"
	"fmt"
	"github.com/paddie/gokmp"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

var kmp *gokmp.KMP
var max_workers int
var lookingfor string
var wg sync.WaitGroup
var work_queue chan string
var inv_ext map[string]bool

func main() {
	flag.Parse()
	// TODO: handle incorrect input
	getReady(flag.Arg(0))
	crawlFolder()
	wg.Wait()

	fmt.Printf("Done.")
}

func getReady(lookingFor string) {
	// put cpu cores to work
	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs)
	setupInvalidFileExtSet()

	// TODO: handle worker pool and queue better
	max_workers = 6765 + 2584
	work_queue = make(chan string, 46368)
	for i := 0; i < max_workers-2584; i++ {
		go searchWorker(work_queue)
	}
	go func() { // wake the rest concurrently
		for i := 0; i < max_workers-6765; i++ {
			go searchWorker(work_queue)
		}
	}()

	lookingfor = flag.Arg(0)
	fmt.Println("Searching for ->", lookingfor)
	kmp, _ = gokmp.NewKMP(lookingfor)
}

func crawlFolder() {
	current := "."
	filepath.Walk(current, visit)
	close(work_queue)
}

func visit(path string, fileInfo os.FileInfo, err error) error {
	if fileInfo.IsDir() || notValidFileExtension(path) {
		return nil
	}
	work_queue <- path
	return nil
}

func searchWorker(work_queue <-chan string) {
	wg.Add(1)
	defer wg.Done()

	for {
		select {
		case path, open := <-work_queue:
			if !open {
				return
			}
			kmp := kmp
			x, err := ioutil.ReadFile(path)
			if err != nil {
				continue //fail gracefully
			}
			if kmp.ContainedIn(string(x)) {
				go printPath(path)
			}
		}
	}
}

func notValidFileExtension(path string) bool {
	ext := filepath.Ext(path)
	return inv_ext[ext]
}

func setupInvalidFileExtSet() {
	// TODO: exclude files better
	inv_ext = map[string]bool{
		".exe": true, ".dll": true, ".msi": true,
		".obj": true, ".bsc": true, ".pdb": true,
		".ilk": true, ".idb": true, ".psd": true,
		".sdf": true, "": true,
	}
}

func printPath(path string) {
	wg.Add(1)
	defer wg.Done()
	fmt.Println(path)
}
