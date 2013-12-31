#GOREP v0.01a
## 

###Fast linear string search in Go. Find files, fast.

* Uses the [KMP algorithm](http://en.wikipedia.org/wiki/Knuth%E2%80%93Morris%E2%80%93Pratt_algorithm). Linear, O(n), time complexity string search. ([Thanks paddie.](https://github.com/paddie/gokmp))
* Concurrent searching with goroutines.
* Recursively search current directory.
* UTF-8 compliant.

##### Todo
* Better memory handling of opened files
* Case sensitivity
* Better error handling
* Better flag commands
* Regex


###### HowTo

1. [install Go.](ttp://golang.org/doc/install)
2. Build the executable. (go build or go install)
3. gorep "SearchMe"

