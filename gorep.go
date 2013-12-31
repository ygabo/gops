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
	"time"
)

type Work struct {
	path string
	f    os.FileInfo
}

var kmp *gokmp.KMP
var lookingfor string
var wg sync.WaitGroup
var num_workers int
var inv_ext map[string]bool
var max_workers int
var work_queue chan Work

func visit(path string, f os.FileInfo, err error) error {
	work_queue <- Work{path, f}
	return nil
}


func notValidFileExtension(path string) bool {
  ext := filepath.Ext(path)
  return inv_ext[ext]
}

func printPath(path string) {
  wg.Add(1)
  defer wg.Done()
  fmt.Println(path)
}

func searchWorker(work_queue <-chan Work) {
	wg.Add(1)
	defer wg.Done()
	
	for {
		select {
		case info, open := <-work_queue:
			if !open {
				return
			}
			if info.f.IsDir() || notValidFileExtension(info.path){
				continue
			}
			kmp := kmp
			x, _ := ioutil.ReadFile(info.path)     
      
			if kmp.ContainedIn(string(x)) {
				go printPath(info.path)
			}
		}
	}
}

func setupInvalidFileExtMap(){
  // TODO: exclude files better
  inv_ext = map[string]bool{
    "":true,     ".exe":true, ".dll":true,
    ".obj":true, ".bsc":true,
    ".ilk":true, ".pdb":true,
    ".msi":true, ".idb":true,
    ".sdf":true, ".psd":true,
    
  }
}

func getReady(lookingFor string){
  setupInvalidFileExtMap()
  max_workers = 6765+2584  
	work_queue = make(chan Work, 46368)
	for i := 0; i < max_workers-2584; i++ {
			go searchWorker(work_queue)
	}	
	go func() {
		for i := 0; i < max_workers-6765; i++ {
			go searchWorker(work_queue)
		}
	}()
  
	lookingfor = flag.Arg(0)
	fmt.Println("Searching for ->", lookingfor)
	kmp, _ = gokmp.NewKMP(lookingfor)
}

func crawlFolder(){
	current := "."
	filepath.Walk(current, visit)
	close(work_queue)
}

func main() {
	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs)
	flag.Parse()
	start := time.Now()
  
	getReady(flag.Arg(0))
  crawlFolder()
	wg.Wait()

	elapsed := time.Since(start)
	fmt.Printf("took %s", elapsed)
}
