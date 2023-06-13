package main

import "fmt"

var version string
var commit string

func main() {
	fmt.Println("GATEWAY", version, commit)
}
