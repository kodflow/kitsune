package main

import "fmt"

var version string
var commit string

func main() {
	fmt.Println("CACHE", version, commit)
}
