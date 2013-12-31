package main

import (
	"flag"
	"fmt"
	"github.com/paddie/gokmp"
	"io/ioutil"
	"sync"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"time"
)

func visit(path string, f os.FileInfo, err error) error {
	wg.Add(1)
	go findText(path, f)
	return nil
}

func findText(filename string, f os.FileInfo) {
	defer wg.Done()
	if f.IsDir() {
		return
	}
	kmp := kmp
	x, _ := ioutil.ReadFile(filename)
	if kmp.ContainedIn(string(x)) {
		fmt.Println(filename)
	}
}

var kmp *gokmp.KMP
var lookingfor string
var wg sync.WaitGroup
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs)
	flag.Parse()
	if *cpuprofile != "" {
        f, err := os.Create(*cpuprofile)
        if err != nil {
           fmt.Println("lol")
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
  }
	start := time.Now()
	lookingfor = flag.Arg(0)
	fmt.Println("Searching for -> ", lookingfor)
	kmp, _ = gokmp.NewKMP(lookingfor)
	root := "."
	filepath.Walk(root, visit)
	wg.Wait()
	elapsed := time.Since(start)
    fmt.Printf("took %s", elapsed)
}
